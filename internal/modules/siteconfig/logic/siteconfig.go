// Package logic 站点配置发布业务实现(M3)。
//
// 关键设计:
//   - 发布记录存"未注入 database 段"的原文 —— 库里不落站点DB明文密码;
//   - 发布/回滚时保留原文; 仅覆盖 database.default.link 为站点登记库连接;
//   - paas 段只改/补 app_key、app_secret, 不整块替换;
//   - 查看当前配置时对 link 里的密码脱敏。
package logic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
	"github.com/JarvanDante/my_manage_service/internal/shared/nacospub"
)

var (
	passMaskRe          = regexp.MustCompile(`(pgsql:[^:@]+:)[^@]+(@)`)
	appSecretMaskRe     = regexp.MustCompile(`(?m)^([ \t]*app_secret\s*:\s*).+$`)
	linkLineRe          = regexp.MustCompile(`^([ \t]*)link\s*:.*$`)
	defaultRe           = regexp.MustCompile(`^([ \t]*)default\s*:\s*(#.*)?$`)
	paasAppKeyLineRe    = regexp.MustCompile(`^([ \t]*)app_key\s*:.*$`)
	paasAppSecretLineRe = regexp.MustCompile(`^([ \t]*)app_secret\s*:.*$`)
	paasKeyRe           = regexp.MustCompile(`^paas\s*:\s*(#.*)?$`)
)

type sSiteConfig struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.ISiteConfig { return &sSiteConfig{repo: repo} }

func (s *sSiteConfig) Current(ctx context.Context, siteId int64) (*service.CurrentConfigDTO, error) {
	site, ns, dataId, group, err := s.target(ctx, siteId)
	if err != nil {
		return nil, err
	}
	content, err := nacospub.Get(ctx, ns, dataId, group)
	if err != nil {
		return nil, err
	}
	masked := passMaskRe.ReplaceAllString(content, "${1}******${2}")
	masked = appSecretMaskRe.ReplaceAllString(masked, `${1}"******"`)
	return &service.CurrentConfigDTO{
		SiteCode: site.SiteCode, Env: site.Env, NamespaceId: ns,
		DataId: dataId, Group: group,
		Content: masked,
		Exists:  content != "",
	}, nil
}

func (s *sSiteConfig) Publish(ctx context.Context, in service.PublishInput) (*service.PublishDTO, error) {
	if in.Content == "" {
		return nil, gerror.New("配置内容不能为空")
	}
	return s.doPublish(ctx, in.SiteId, in.Content, in.Operator, in.Remark)
}

func (s *sSiteConfig) History(ctx context.Context, siteId int64, page, size int) (*service.HistoryDTO, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if _, err := s.mustSite(ctx, siteId); err != nil {
		return nil, err
	}
	list, total, err := s.repo.History(ctx, siteId, page, size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.HistoryItem, 0, len(list))
	for _, r := range list {
		out = append(out, &service.HistoryItem{
			Id: r.Id, Version: r.Version, Env: r.Env, Content: r.Content,
			Operator: r.Operator, Remark: r.Remark, CreatedAt: fmtTime(r.CreatedAt),
		})
	}
	return &service.HistoryDTO{List: out, Total: total, Page: page, Size: size}, nil
}

func (s *sSiteConfig) Rollback(ctx context.Context, siteId int64, version int, operator int64) (*service.PublishDTO, error) {
	if version <= 0 {
		return nil, gerror.New("版本号无效")
	}
	rec, err := s.repo.FindVersion(ctx, siteId, version)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, gerror.Newf("版本 v%d 不存在", version)
	}
	return s.doPublish(ctx, siteId, rec.Content, operator, fmt.Sprintf("回滚自 v%d", version))
}

// doPublish 发布主流程: 校验yaml -> 注入 database/paas 段 -> 发布Nacos -> 落发布记录(存原文)。
func (s *sSiteConfig) doPublish(ctx context.Context, siteId int64, content string, operator int64, remark string) (*service.PublishDTO, error) {
	site, ns, dataId, group, err := s.target(ctx, siteId)
	if err != nil {
		return nil, err
	}
	rendered, err := injectDatabase(ctx, content, site)
	if err != nil {
		return nil, err
	}
	rendered, err = injectPaas(ctx, rendered, site)
	if err != nil {
		return nil, err
	}
	if err = nacospub.Publish(ctx, ns, dataId, group, rendered); err != nil {
		return nil, err
	}
	version, err := s.repo.NextVersion(ctx, siteId)
	if err != nil {
		return nil, err
	}
	if _, err = s.repo.CreateRecord(ctx, &entity.SiteConfigPublish{
		SiteId: siteId, Env: site.Env, Version: version,
		Content: content, Operator: operator, Remark: remark,
	}); err != nil {
		return nil, err
	}
	return &service.PublishDTO{Version: version, DataId: dataId}, nil
}

// injectDatabase 校验 YAML 合法后, 全文文本级处理:
//   - 其它段落原样保留;
//   - 若已有 database 段: 只改/补 default.link, debug 等字段一字不动;
//   - 若无 database 段: 末尾追加仅含 link 的 database 段。
func injectDatabase(ctx context.Context, content string, site *entity.Site) (string, error) {
	if _, err := gjson.LoadContentType(gjson.ContentTypeYaml, []byte(content)); err != nil {
		return "", gerror.Wrap(err, "配置内容不是合法的 YAML")
	}
	pass, err := kit.Decrypt(site.DbPass)
	if err != nil {
		return "", gerror.Wrap(err, "解密站点数据库密码失败")
	}
	link := fmt.Sprintf("pgsql:%s:%s@tcp(%s:%d)/%s?sslmode=disable",
		site.DbUser, pass, site.DbHost, site.DbPort, site.DbName)

	if strings.TrimSpace(extractTopLevelYAMLBlock(content, "database")) != "" {
		return replaceOrInsertDatabaseLink(content, link), nil
	}
	dbBlock := "database:\n  default:\n    link: " + strconv.Quote(link) + "\n"
	body := strings.TrimRight(content, "\n")
	if body == "" {
		return dbBlock, nil
	}
	return body + "\n\n" + dbBlock, nil
}

// injectPaas 将站点 APPKEY/APPSECRET 写入顶层 paas 段。
// 已有 paas 段时只改/补 app_key、app_secret，ad_base / media_base 等其余字段一字不动。
// 历史站尚无凭证时跳过, 不阻断配置发布。
func injectPaas(ctx context.Context, content string, site *entity.Site) (string, error) {
	_ = ctx
	if strings.TrimSpace(site.AppKey) == "" || strings.TrimSpace(site.AppSecret) == "" {
		return content, nil
	}
	secret, err := kit.Decrypt(site.AppSecret)
	if err != nil {
		return "", gerror.Wrap(err, "解密站点 APPSECRET 失败")
	}
	if strings.TrimSpace(extractTopLevelYAMLBlock(content, "paas")) != "" {
		return replaceOrInsertPaasCreds(content, site.AppKey, secret), nil
	}
	block := "paas:\n  app_key: " + strconv.Quote(site.AppKey) + "\n  app_secret: " + strconv.Quote(secret) + "\n"
	body := strings.TrimRight(content, "\n")
	if body == "" {
		return block, nil
	}
	return body + "\n\n" + block, nil
}

// replaceOrInsertPaasCreds 在原文 paas 段内替换或插入凭证行, 其它行保持不变。
func replaceOrInsertPaasCreds(content, appKey, secret string) string {
	keyQuoted := strconv.Quote(appKey)
	secretQuoted := strconv.Quote(secret)
	block := extractTopLevelYAMLBlock(content, "paas")
	rest := strings.TrimRight(removeTopLevelYAMLKey(content, "paas"), "\n")

	lines := strings.Split(block, "\n")
	hasKey, hasSecret := false, false
	for i, line := range lines {
		if m := paasAppKeyLineRe.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + "app_key: " + keyQuoted
			hasKey = true
			continue
		}
		if m := paasAppSecretLineRe.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + "app_secret: " + secretQuoted
			hasSecret = true
		}
	}
	if !hasKey || !hasSecret {
		out := make([]string, 0, len(lines)+2)
		inserted := false
		for _, line := range lines {
			out = append(out, line)
			if !inserted && paasKeyRe.MatchString(line) {
				if !hasKey {
					out = append(out, "  app_key: "+keyQuoted)
				}
				if !hasSecret {
					out = append(out, "  app_secret: "+secretQuoted)
				}
				inserted = true
			}
		}
		if !inserted {
			if !hasKey {
				out = append(out, "  app_key: "+keyQuoted)
			}
			if !hasSecret {
				out = append(out, "  app_secret: "+secretQuoted)
			}
		}
		lines = out
	}
	newBlock := strings.Join(lines, "\n")
	if rest == "" {
		return newBlock
	}
	return rest + "\n\n" + newBlock
}

// replaceOrInsertDatabaseLink 在原文 database 段内替换或插入 link 行, 其它行保持不变。
func replaceOrInsertDatabaseLink(content, link string) string {
	quoted := strconv.Quote(link)
	block := extractTopLevelYAMLBlock(content, "database")
	rest := strings.TrimRight(removeTopLevelYAMLKey(content, "database"), "\n")

	lines := strings.Split(block, "\n")
	replaced := false
	for i, line := range lines {
		if m := linkLineRe.FindStringSubmatch(line); m != nil {
			lines[i] = m[1] + "link: " + quoted
			replaced = true
			break
		}
	}
	if !replaced {
		out := make([]string, 0, len(lines)+2)
		inserted := false
		for _, line := range lines {
			out = append(out, line)
			if m := defaultRe.FindStringSubmatch(line); m != nil {
				out = append(out, m[1]+"  link: "+quoted)
				inserted = true
			}
		}
		if !inserted {
			out = append(out, "  default:", "    link: "+quoted)
		}
		lines = out
	}
	newBlock := strings.Join(lines, "\n")
	if rest == "" {
		return newBlock
	}
	return rest + "\n\n" + newBlock
}

// extractTopLevelYAMLBlock 取出顶层 key 整块文本(含 key 行)。
func extractTopLevelYAMLBlock(content, key string) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0)
	capturing := false
	for _, line := range lines {
		if isYAMLTopLevelKeyLine(line) {
			name, ok := yamlTopLevelKeyName(line)
			if ok && name == key {
				capturing = true
				out = append(out, line)
				continue
			}
			if capturing {
				break
			}
		}
		if capturing {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// removeTopLevelYAMLKey 删除顶层 key 整块(含其缩进子行), 其它文本原样保留。
func removeTopLevelYAMLKey(content, key string) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	skipping := false
	for _, line := range lines {
		if isYAMLTopLevelKeyLine(line) {
			name, ok := yamlTopLevelKeyName(line)
			if ok && name == key {
				skipping = true
				continue
			}
			skipping = false
		}
		if skipping {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func isYAMLTopLevelKeyLine(line string) bool {
	if line == "" {
		return false
	}
	if line[0] == ' ' || line[0] == '\t' || line[0] == '#' {
		return false
	}
	return strings.Contains(line, ":")
}

func yamlTopLevelKeyName(line string) (string, bool) {
	// 形如 "database:" / "database: something" / "database: # cmt"
	i := strings.Index(line, ":")
	if i <= 0 {
		return "", false
	}
	name := strings.TrimSpace(line[:i])
	if name == "" || strings.ContainsAny(name, " \t") {
		return "", false
	}
	return name, true
}

// ---------- helpers ----------

func (s *sSiteConfig) mustSite(ctx context.Context, siteId int64) (*entity.Site, error) {
	if siteId <= 0 {
		return nil, gerror.New("站点ID无效")
	}
	site, err := s.repo.FindSite(ctx, siteId)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, gerror.New("站点不存在")
	}
	return site, nil
}

func (s *sSiteConfig) target(ctx context.Context, siteId int64) (site *entity.Site, ns, dataId, group string, err error) {
	site, err = s.mustSite(ctx, siteId)
	if err != nil {
		return
	}
	ns, err = nacospub.NamespaceForEnv(ctx, site.Env)
	if err != nil {
		return
	}
	dataId = site.SiteCode + ".yaml"
	group = nacospub.Group(ctx)
	return
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

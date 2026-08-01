// Package logic 站点配置发布业务实现(M3)。
//
// 关键设计:
//   - 发布记录存"未注入 database 段"的原文 —— 库里不落站点DB明文密码;
//   - 发布/回滚时动态注入: 用站点登记信息(AES解密)渲染 database.default.link;
//   - 查看当前配置时对 link 里的密码脱敏。
package logic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
	"github.com/JarvanDante/my_manage_service/internal/shared/nacospub"
)

var passMaskRe = regexp.MustCompile(`(pgsql:[^:@]+:)[^@]+(@)`)

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
	return &service.CurrentConfigDTO{
		SiteCode: site.SiteCode, Env: site.Env, NamespaceId: ns,
		DataId: dataId, Group: group,
		Content: passMaskRe.ReplaceAllString(content, "${1}******${2}"),
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

// doPublish 发布主流程: 校验yaml -> 注入database段 -> 发布Nacos -> 落发布记录(存原文)。
func (s *sSiteConfig) doPublish(ctx context.Context, siteId int64, content string, operator int64, remark string) (*service.PublishDTO, error) {
	site, ns, dataId, group, err := s.target(ctx, siteId)
	if err != nil {
		return nil, err
	}
	rendered, err := injectDatabase(ctx, content, site)
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

// injectDatabase 解析 yaml 并注入 database.default.link(覆盖同名段)。
func injectDatabase(ctx context.Context, content string, site *entity.Site) (string, error) {
	j, err := gjson.LoadContentType(gjson.ContentTypeYaml, []byte(content))
	if err != nil {
		return "", gerror.Wrap(err, "配置内容不是合法的 YAML")
	}
	pass, err := kit.Decrypt(site.DbPass)
	if err != nil {
		return "", gerror.Wrap(err, "解密站点数据库密码失败")
	}
	link := fmt.Sprintf("pgsql:%s:%s@tcp(%s:%d)/%s?sslmode=disable",
		site.DbUser, pass, site.DbHost, site.DbPort, site.DbName)
	if err = j.Set("database.default.link", link); err != nil {
		return "", err
	}
	out, err := j.ToYamlString()
	if err != nil {
		return "", gerror.Wrap(err, "生成 YAML 失败")
	}
	return out, nil
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

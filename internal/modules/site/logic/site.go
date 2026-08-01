// Package logic 站点模块业务实现(M2)。
package logic

import (
	"context"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/dao/sitedb"
	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
)

var (
	codeRe   = regexp.MustCompile(`^[a-z][a-z0-9]{1,31}$`)
	domainRe = regexp.MustCompile(`^([a-z0-9-]+\.)+[a-z]{2,}$`)
)

type sSite struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.ISite { return &sSite{repo: repo} }

func (s *sSite) List(ctx context.Context, in service.ListInput) (*service.ListDTO, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.List(ctx, domain.Filter{
		Keyword: in.Keyword, MerchantId: in.MerchantId, Env: in.Env, Status: in.Status,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.SiteDTO, 0, len(list))
	for _, it := range list {
		d := toDTO(&it.Site)
		d.MerchantName = it.MerchantName
		out = append(out, d)
	}
	return &service.ListDTO{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sSite) Detail(ctx context.Context, id int64) (*service.SiteDetailDTO, error) {
	site, err := s.mustSite(ctx, id)
	if err != nil {
		return nil, err
	}
	domains, err := s.repo.Domains(ctx, id)
	if err != nil {
		return nil, err
	}
	ds := make([]*service.DomainDTO, 0, len(domains))
	for _, d := range domains {
		ds = append(ds, &service.DomainDTO{
			Id: d.Id, Domain: d.Domain, IsMain: d.IsMain, Https: d.Https, Status: d.Status,
		})
	}
	return &service.SiteDetailDTO{SiteDTO: *toDTO(site), Domains: ds}, nil
}

func checkSave(in service.SaveInput) error {
	if in.Name == "" {
		return gerror.New("站点名必填")
	}
	if in.Env != "dev" && in.Env != "test" && in.Env != "prod" {
		return gerror.New("env 仅支持 dev/test/prod")
	}
	if in.DbHost == "" || in.DbName == "" || in.DbUser == "" {
		return gerror.New("数据库登记信息(host/name/user)必填")
	}
	if in.DbPort <= 0 || in.DbPort > 65535 {
		return gerror.New("数据库端口不合法")
	}
	return nil
}

func (s *sSite) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	in.SiteCode = strings.ToLower(strings.TrimSpace(in.SiteCode))
	if !codeRe.MatchString(in.SiteCode) {
		return 0, gerror.New("site_code 须为 2~32 位小写字母开头的字母数字")
	}
	if err := checkSave(in); err != nil {
		return 0, err
	}
	if in.DbPass == "" {
		return 0, gerror.New("数据库密码必填")
	}
	if in.MerchantId <= 0 {
		return 0, gerror.New("请指定商户")
	}
	ok, err := s.repo.MerchantExists(ctx, in.MerchantId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, gerror.New("商户不存在或已停用")
	}
	exist, err := s.repo.FindByCode(ctx, in.SiteCode)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("site_code 已存在")
	}
	encPass, err := kit.Encrypt(in.DbPass)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, &entity.Site{
		MerchantId: in.MerchantId, SiteCode: in.SiteCode, Name: in.Name, Env: in.Env,
		DbHost: in.DbHost, DbPort: in.DbPort, DbName: in.DbName, DbUser: in.DbUser,
		DbPass: encPass, Remark: in.Remark,
	})
}

func (s *sSite) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("站点ID无效")
	}
	if err := checkSave(in); err != nil {
		return err
	}
	if _, err := s.mustSite(ctx, in.Id); err != nil {
		return err
	}
	if in.MerchantId > 0 {
		ok, err := s.repo.MerchantExists(ctx, in.MerchantId)
		if err != nil {
			return err
		}
		if !ok {
			return gerror.New("商户不存在或已停用")
		}
	}
	encPass := ""
	if in.DbPass != "" {
		var err error
		if encPass, err = kit.Encrypt(in.DbPass); err != nil {
			return err
		}
	}
	return s.repo.Update(ctx, &entity.Site{
		Id: in.Id, MerchantId: in.MerchantId, Name: in.Name, Env: in.Env,
		DbHost: in.DbHost, DbPort: in.DbPort, DbName: in.DbName, DbUser: in.DbUser,
		DbPass: encPass, Remark: in.Remark,
	})
}

// SetStatus 0筹备 1上线 2停用。
// TODO(M4): status→1 前走 provision-check 清单(域名/DB/迁移/配置全过才放行)。
func (s *sSite) SetStatus(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.New("status 仅支持 0/1/2")
	}
	if _, err := s.mustSite(ctx, id); err != nil {
		return err
	}
	return s.repo.SetStatus(ctx, id, status)
}

// DbCheck 用登记信息真实连接站点库(半自动开站的校验环节)。
func (s *sSite) DbCheck(ctx context.Context, id int64) (*service.DbCheckDTO, error) {
	site, err := s.mustSite(ctx, id)
	if err != nil {
		return nil, err
	}
	pass, err := kit.Decrypt(site.DbPass)
	if err != nil {
		return nil, gerror.Wrap(err, "解密数据库密码失败(MANAGE_SECRET 是否与登记时一致?)")
	}
	ms, err := sitedb.Ping(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName)
	if err != nil {
		return &service.DbCheckDTO{Ok: false, Message: err.Error()}, nil
	}
	return &service.DbCheckDTO{Ok: true, LatencyMs: ms, Message: "连接正常"}, nil
}

// ---------- 域名 ----------

func (s *sSite) BindDomain(ctx context.Context, siteId int64, domainName string, https int) (int64, error) {
	if _, err := s.mustSite(ctx, siteId); err != nil {
		return 0, err
	}
	domainName = strings.ToLower(strings.TrimSpace(domainName))
	if !domainRe.MatchString(domainName) {
		return 0, gerror.New("域名格式不正确")
	}
	exist, err := s.repo.FindDomain(ctx, domainName)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("该域名已被绑定")
	}
	// 首个域名自动设为主域名
	list, err := s.repo.Domains(ctx, siteId)
	if err != nil {
		return 0, err
	}
	isMain := 0
	if len(list) == 0 {
		isMain = 1
	}
	if https != 0 && https != 1 {
		https = 1
	}
	return s.repo.BindDomain(ctx, siteId, domainName, https, isMain)
}

func (s *sSite) UnbindDomain(ctx context.Context, siteId, domainId int64) error {
	d, err := s.repo.FindDomainById(ctx, domainId)
	if err != nil {
		return err
	}
	if d == nil || d.SiteId != siteId {
		return gerror.New("域名不存在或不属于该站点")
	}
	if d.IsMain == 1 {
		return gerror.New("主域名不可解绑, 请先切换主域名")
	}
	return s.repo.UnbindDomain(ctx, domainId)
}

func (s *sSite) SetMainDomain(ctx context.Context, siteId, domainId int64) error {
	d, err := s.repo.FindDomainById(ctx, domainId)
	if err != nil {
		return err
	}
	if d == nil || d.SiteId != siteId {
		return gerror.New("域名不存在或不属于该站点")
	}
	return s.repo.SetMainDomain(ctx, siteId, domainId)
}

// ---------- helpers ----------

func (s *sSite) mustSite(ctx context.Context, id int64) (*entity.Site, error) {
	if id <= 0 {
		return nil, gerror.New("站点ID无效")
	}
	site, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, gerror.New("站点不存在")
	}
	return site, nil
}

func toDTO(e *entity.Site) *service.SiteDTO {
	return &service.SiteDTO{
		Id: e.Id, MerchantId: e.MerchantId, SiteCode: e.SiteCode, Name: e.Name,
		Env: e.Env, Status: e.Status,
		DbHost: e.DbHost, DbPort: e.DbPort, DbName: e.DbName, DbUser: e.DbUser,
		DbPassMask: "******", Remark: e.Remark, CreatedAt: fmtTime(e.CreatedAt),
	}
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

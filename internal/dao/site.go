package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	sitedomain "github.com/JarvanDante/my_manage_service/internal/modules/site/domain"
)

type siteRepo struct{}

// NewSiteRepo 返回 site 领域仓储实现。
func NewSiteRepo() sitedomain.Repository { return &siteRepo{} }

func (r *siteRepo) List(ctx context.Context, f sitedomain.Filter, page, size int) ([]*sitedomain.SiteItem, int, error) {
	m := g.Model("site s").Ctx(ctx)
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(s.site_code ILIKE ? OR s.name ILIKE ?)", kw, kw)
	}
	if f.MerchantId > 0 {
		m = m.Where("s.merchant_id", f.MerchantId)
	}
	if f.Env != "" {
		m = m.Where("s.env", f.Env)
	}
	if f.Status >= 0 {
		m = m.Where("s.status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*sitedomain.SiteItem
	err = m.Clone().
		LeftJoin("merchant mc", "mc.id=s.merchant_id").
		Fields("s.*, mc.name AS merchant_name").
		OrderDesc("s.id").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *siteRepo) FindById(ctx context.Context, id int64) (*entity.Site, error) {
	var s *entity.Site
	err := g.Model("site").Ctx(ctx).Where("id", id).Scan(&s)
	return s, err
}

// FindItemById 详情用: 联表商户名。
func (r *siteRepo) FindItemById(ctx context.Context, id int64) (*sitedomain.SiteItem, error) {
	var item *sitedomain.SiteItem
	err := g.Model("site s").Ctx(ctx).
		LeftJoin("merchant mc", "mc.id=s.merchant_id").
		Fields("s.*, mc.name AS merchant_name").
		Where("s.id", id).Scan(&item)
	return item, err
}

func (r *siteRepo) FindByCode(ctx context.Context, code string) (*entity.Site, error) {
	var s *entity.Site
	err := g.Model("site").Ctx(ctx).Where("site_code", code).Scan(&s)
	return s, err
}

func (r *siteRepo) Create(ctx context.Context, s *entity.Site) (int64, error) {
	res, err := g.Model("site").Ctx(ctx).Data(g.Map{
		"merchant_id": s.MerchantId, "site_code": s.SiteCode, "name": s.Name,
		"env": s.Env, "status": 0,
		"db_host": s.DbHost, "db_port": s.DbPort, "db_name": s.DbName,
		"db_user": s.DbUser, "db_pass": s.DbPass, "remark": s.Remark,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update 更新站点(site_code 不可改; db_pass 为空表示不改)。
func (r *siteRepo) Update(ctx context.Context, s *entity.Site) error {
	data := g.Map{
		"merchant_id": s.MerchantId, "name": s.Name, "env": s.Env,
		"db_host": s.DbHost, "db_port": s.DbPort, "db_name": s.DbName,
		"db_user": s.DbUser, "remark": s.Remark, "updated_at": gtime.Now(),
	}
	if s.DbPass != "" {
		data["db_pass"] = s.DbPass
	}
	_, err := g.Model("site").Ctx(ctx).Where("id", s.Id).Data(data).Update()
	return err
}

func (r *siteRepo) SetStatus(ctx context.Context, id int64, status int) error {
	_, err := g.Model("site").Ctx(ctx).Where("id", id).Data(g.Map{
		"status": status, "updated_at": gtime.Now(),
	}).Update()
	return err
}

// MerchantExists 校验商户存在且启用(控制面表同库, dao 层内联查询)。
func (r *siteRepo) MerchantExists(ctx context.Context, merchantId int64) (bool, error) {
	n, err := g.Model("merchant").Ctx(ctx).Where("id", merchantId).Where("status", 1).Count()
	return n > 0, err
}

// ---------- 域名 ----------

func (r *siteRepo) Domains(ctx context.Context, siteId int64) ([]*entity.SiteDomain, error) {
	var list []*entity.SiteDomain
	err := g.Model("site_domain").Ctx(ctx).Where("site_id", siteId).
		Order("is_main desc, id asc").Scan(&list)
	return list, err
}

func (r *siteRepo) FindDomain(ctx context.Context, domain string) (*entity.SiteDomain, error) {
	var d *entity.SiteDomain
	err := g.Model("site_domain").Ctx(ctx).Where("domain", domain).Scan(&d)
	return d, err
}

func (r *siteRepo) FindDomainById(ctx context.Context, id int64) (*entity.SiteDomain, error) {
	var d *entity.SiteDomain
	err := g.Model("site_domain").Ctx(ctx).Where("id", id).Scan(&d)
	return d, err
}

func (r *siteRepo) BindDomain(ctx context.Context, siteId int64, domain string, https, isMain int) (int64, error) {
	res, err := g.Model("site_domain").Ctx(ctx).Data(g.Map{
		"site_id": siteId, "domain": domain, "https": https, "is_main": isMain, "status": 1,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *siteRepo) UnbindDomain(ctx context.Context, id int64) error {
	_, err := g.Model("site_domain").Ctx(ctx).Where("id", id).Delete()
	return err
}

// SetMainDomain 事务: 清掉该站其它主域名标记, 设新主域名。
func (r *siteRepo) SetMainDomain(ctx context.Context, siteId, domainId int64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("site_domain").Ctx(ctx).
			Where("site_id", siteId).Data(g.Map{"is_main": 0}).Update(); err != nil {
			return err
		}
		_, err := tx.Model("site_domain").Ctx(ctx).
			Where("id", domainId).Where("site_id", siteId).
			Data(g.Map{"is_main": 1}).Update()
		return err
	})
}

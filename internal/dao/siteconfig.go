package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	scdomain "github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/domain"
)

type siteConfigRepo struct{}

// NewSiteConfigRepo 返回 siteconfig 领域仓储实现。
func NewSiteConfigRepo() scdomain.Repository { return &siteConfigRepo{} }

func (r *siteConfigRepo) FindSite(ctx context.Context, siteId int64) (*entity.Site, error) {
	var s *entity.Site
	err := g.Model("site").Ctx(ctx).Where("id", siteId).Scan(&s)
	return s, err
}

func (r *siteConfigRepo) NextVersion(ctx context.Context, siteId int64) (int, error) {
	v, err := g.Model("site_config_publish").Ctx(ctx).
		Where("site_id", siteId).Fields("coalesce(max(version), 0)").Value()
	if err != nil {
		return 0, err
	}
	return v.Int() + 1, nil
}

func (r *siteConfigRepo) CreateRecord(ctx context.Context, rec *entity.SiteConfigPublish) (int64, error) {
	res, err := g.Model("site_config_publish").Ctx(ctx).Data(g.Map{
		"site_id": rec.SiteId, "env": rec.Env, "version": rec.Version,
		"content": rec.Content, "operator": rec.Operator, "remark": rec.Remark,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *siteConfigRepo) History(ctx context.Context, siteId int64, page, size int) ([]*entity.SiteConfigPublish, int, error) {
	m := g.Model("site_config_publish").Ctx(ctx).Where("site_id", siteId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.SiteConfigPublish
	err = m.Clone().OrderDesc("version").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *siteConfigRepo) FindVersion(ctx context.Context, siteId int64, version int) (*entity.SiteConfigPublish, error) {
	var rec *entity.SiteConfigPublish
	err := g.Model("site_config_publish").Ctx(ctx).
		Where("site_id", siteId).Where("version", version).Scan(&rec)
	return rec, err
}

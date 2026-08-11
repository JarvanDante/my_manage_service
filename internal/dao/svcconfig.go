package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	scdomain "github.com/JarvanDante/my_manage_service/internal/modules/svcconfig/domain"
)

type svcConfigRepo struct{}

// NewSvcConfigRepo 返回 svcconfig 领域仓储实现。
func NewSvcConfigRepo() scdomain.Repository { return &svcConfigRepo{} }

func (r *svcConfigRepo) NextVersion(ctx context.Context, service, env string) (int, error) {
	v, err := g.Model("service_config_publish").Ctx(ctx).
		Where("service", service).Where("env", env).
		Fields("coalesce(max(version), 0)").Value()
	if err != nil {
		return 0, err
	}
	return v.Int() + 1, nil
}

func (r *svcConfigRepo) CreateRecord(ctx context.Context, rec *entity.ServiceConfigPublish) (int64, error) {
	res, err := g.Model("service_config_publish").Ctx(ctx).Data(g.Map{
		"service": rec.Service, "env": rec.Env, "version": rec.Version,
		"content": rec.Content, "operator": rec.Operator, "remark": rec.Remark,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *svcConfigRepo) History(ctx context.Context, service, env string, page, size int) ([]*entity.ServiceConfigPublish, int, error) {
	m := g.Model("service_config_publish").Ctx(ctx).Where("service", service).Where("env", env)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ServiceConfigPublish
	err = m.Clone().OrderDesc("version").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *svcConfigRepo) FindVersion(ctx context.Context, service, env string, version int) (*entity.ServiceConfigPublish, error) {
	var rec *entity.ServiceConfigPublish
	err := g.Model("service_config_publish").Ctx(ctx).
		Where("service", service).Where("env", env).Where("version", version).Scan(&rec)
	return rec, err
}

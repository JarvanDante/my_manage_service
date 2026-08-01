package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	mondomain "github.com/JarvanDante/my_manage_service/internal/modules/monitor/domain"
)

type monitorRepo struct{}

// NewMonitorRepo 返回 monitor 领域仓储实现。
func NewMonitorRepo() mondomain.Repository { return &monitorRepo{} }

func (r *monitorRepo) FindSite(ctx context.Context, siteId int64) (*entity.Site, error) {
	var s *entity.Site
	err := g.Model("site").Ctx(ctx).Where("id", siteId).Scan(&s)
	return s, err
}

func (r *monitorRepo) ListOnlineSites(ctx context.Context) ([]*entity.Site, error) {
	var list []*entity.Site
	err := g.Model("site").Ctx(ctx).Where("status", 1).Order("id asc").Scan(&list)
	return list, err
}

func (r *monitorRepo) Logs(ctx context.Context, f mondomain.LogFilter, page, size int) ([]*mondomain.LogItem, int, error) {
	m := g.Model("manage_log ml").Ctx(ctx)
	if f.AdminId > 0 {
		m = m.Where("ml.admin_id", f.AdminId)
	}
	if f.Method != "" {
		m = m.Where("ml.method", f.Method)
	}
	if f.Keyword != "" {
		m = m.Where("ml.path ILIKE ?", "%"+f.Keyword+"%")
	}
	if f.StartDate != "" {
		m = m.Where("ml.created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		m = m.Where("ml.created_at < ?::date + interval '1 day'", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*mondomain.LogItem
	err = m.Clone().
		LeftJoin("manage_admin ma", "ma.id=ml.admin_id").
		Fields("ml.*, ma.username").
		OrderDesc("ml.id").Page(page, size).Scan(&list)
	return list, total, err
}

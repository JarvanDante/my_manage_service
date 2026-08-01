// Package controller 监控/汇总控制器(M5)。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/monitor/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/service"
)

type Controller struct{ mon service.IMonitor }

func New(svc service.IMonitor) *Controller { return &Controller{mon: svc} }

// Health 站点健康。
func (c *Controller) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	d, err := c.mon.Health(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.HealthRes{
		SiteCode: d.SiteCode, Status: d.Status, DbOk: d.DbOk,
		LatencyMs: d.LatencyMs, GooseVersion: d.GooseVersion, Message: d.Message,
	}, nil
}

// Overview 跨站汇总。
func (c *Controller) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	d, err := c.mon.Overview(ctx)
	if err != nil {
		return nil, err
	}
	sites := make([]v1.SiteOverviewItem, 0, len(d.Sites))
	for _, it := range d.Sites {
		sites = append(sites, v1.SiteOverviewItem{
			SiteId: it.SiteId, SiteCode: it.SiteCode, Name: it.Name, Env: it.Env,
			Ok: it.Ok, Error: it.Error,
			UserCount: it.UserCount, TodayNew: it.TodayNew, TodayActive: it.TodayActive,
			PaidAmount: it.PaidAmount, PaidOrders: it.PaidOrders,
		})
	}
	return &v1.OverviewRes{
		SiteCount: d.SiteCount, OkCount: d.OkCount,
		UserCount: d.UserCount, TodayNew: d.TodayNew, TodayActive: d.TodayActive,
		PaidAmount: d.PaidAmount, PaidOrders: d.PaidOrders, Sites: sites,
	}, nil
}

// Logs 审计日志。
func (c *Controller) Logs(ctx context.Context, req *v1.LogListReq) (res *v1.LogListRes, err error) {
	d, err := c.mon.Logs(ctx, service.LogListInput{
		AdminId: req.AdminId, Method: req.Method, Keyword: req.Keyword,
		StartDate: req.StartDate, EndDate: req.EndDate, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.LogItem, 0, len(d.List))
	for _, l := range d.List {
		items = append(items, v1.LogItem{
			Id: l.Id, AdminId: l.AdminId, Username: l.Username,
			Method: l.Method, Path: l.Path, Ip: l.Ip, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.LogListRes{List: items, Total: d.Total, Page: d.Page, Size: d.Size}, nil
}

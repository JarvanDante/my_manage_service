// Package logic 监控/汇总业务实现(M5)。
// 说明: 跨站汇总当前为串行遍历(站点少); 站点多了再改并发 + 短缓存。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_manage_service/internal/dao/sitedb"
	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
)

type sMonitor struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IMonitor { return &sMonitor{repo: repo} }

// Health 站点健康: DB 连通 + 迁移版本。
func (s *sMonitor) Health(ctx context.Context, siteId int64) (*service.HealthDTO, error) {
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
	out := &service.HealthDTO{SiteCode: site.SiteCode, Status: site.Status}
	pass, err := kit.Decrypt(site.DbPass)
	if err != nil {
		out.Message = "解密登记密码失败: " + err.Error()
		return out, nil
	}
	ms, err := sitedb.Ping(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName)
	if err != nil {
		out.Message = err.Error()
		return out, nil
	}
	out.DbOk = true
	out.LatencyMs = ms
	if ver, err := sitedb.GooseVersion(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName); err == nil {
		out.GooseVersion = ver
	}
	out.Message = "正常"
	return out, nil
}

// Overview 跨站汇总: 遍历上线站点的登记库聚合核心指标, 单站失败不影响整体。
func (s *sMonitor) Overview(ctx context.Context) (*service.OverviewDTO, error) {
	sites, err := s.repo.ListOnlineSites(ctx)
	if err != nil {
		return nil, err
	}
	today := gconv.Int(gtime.Now().Format("Ymd"))
	out := &service.OverviewDTO{SiteCount: len(sites)}
	for _, site := range sites {
		item := s.collect(ctx, site, today)
		out.Sites = append(out.Sites, item)
		if item.Ok {
			out.OkCount++
			out.UserCount += item.UserCount
			out.TodayNew += item.TodayNew
			out.TodayActive += item.TodayActive
			out.PaidAmount += item.PaidAmount
			out.PaidOrders += item.PaidOrders
		}
	}
	return out, nil
}

func (s *sMonitor) collect(ctx context.Context, site *entity.Site, today int) *service.SiteOverviewItem {
	item := &service.SiteOverviewItem{
		SiteId: site.Id, SiteCode: site.SiteCode, Name: site.Name, Env: site.Env,
	}
	pass, err := kit.Decrypt(site.DbPass)
	if err != nil {
		item.Error = "解密登记密码失败"
		return item
	}
	st, err := sitedb.QueryStats(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName, today)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	item.Ok = true
	item.UserCount = st.UserCount
	item.TodayNew = st.TodayNew
	item.TodayActive = st.TodayActive
	item.PaidAmount = st.PaidAmount
	item.PaidOrders = st.PaidOrders
	return item
}

// Logs 审计日志查询。
func (s *sMonitor) Logs(ctx context.Context, in service.LogListInput) (*service.LogListDTO, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.Logs(ctx, domain.LogFilter{
		AdminId: in.AdminId, Method: in.Method, Keyword: in.Keyword,
		StartDate: in.StartDate, EndDate: in.EndDate,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.LogDTO, 0, len(list))
	for _, l := range list {
		created := ""
		if l.CreatedAt != nil {
			created = l.CreatedAt.String()
		}
		out = append(out, &service.LogDTO{
			Id: l.Id, AdminId: l.AdminId, Username: l.Username,
			Method: l.Method, Path: l.Path, Ip: l.Ip, CreatedAt: created,
		})
	}
	return &service.LogListDTO{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

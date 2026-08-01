// Package service 监控/汇总对外接口(M5)。
package service

import "context"

type HealthDTO struct {
	SiteCode     string
	Status       int
	DbOk         bool
	LatencyMs    int64
	GooseVersion int64
	Message      string
}

// SiteOverviewItem 单站指标。
type SiteOverviewItem struct {
	SiteId      int64
	SiteCode    string
	Name        string
	Env         string
	Ok          bool
	Error       string
	UserCount   int
	TodayNew    int
	TodayActive int
	PaidAmount  float64
	PaidOrders  int
}

type OverviewDTO struct {
	SiteCount   int // 上线站点数
	OkCount     int // 采集成功数
	UserCount   int // 汇总(仅成功站点)
	TodayNew    int
	TodayActive int
	PaidAmount  float64
	PaidOrders  int
	Sites       []*SiteOverviewItem
}

type LogListInput struct {
	AdminId   int64
	Method    string
	Keyword   string
	StartDate string
	EndDate   string
	Page      int
	Size      int
}

type LogDTO struct {
	Id        int64
	AdminId   int64
	Username  string
	Method    string
	Path      string
	Ip        string
	CreatedAt string
}

type LogListDTO struct {
	List  []*LogDTO
	Total int
	Page  int
	Size  int
}

type IMonitor interface {
	Health(ctx context.Context, siteId int64) (*HealthDTO, error)
	Overview(ctx context.Context) (*OverviewDTO, error)
	Logs(ctx context.Context, in LogListInput) (*LogListDTO, error)
}

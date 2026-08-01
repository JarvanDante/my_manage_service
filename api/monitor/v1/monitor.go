// Package v1 监控/汇总接口契约(M5)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// 站点健康
type HealthReq struct {
	g.Meta `path:"/sites/{id}/health" method:"get" tags:"Manage/Monitor" summary:"站点健康检查"`
	Id     int64 `json:"id" v:"required|min:1#站点ID必填|站点ID必须大于0"`
}
type HealthRes struct {
	SiteCode     string `json:"site_code"`
	Status       int    `json:"status"`
	DbOk         bool   `json:"db_ok"`
	LatencyMs    int64  `json:"latency_ms"`
	GooseVersion int64  `json:"goose_version"`
	Message      string `json:"message"`
}

// 跨站汇总
type SiteOverviewItem struct {
	SiteId      int64   `json:"site_id"`
	SiteCode    string  `json:"site_code"`
	Name        string  `json:"name"`
	Env         string  `json:"env"`
	Ok          bool    `json:"ok"`
	Error       string  `json:"error"`
	UserCount   int     `json:"user_count"`
	TodayNew    int     `json:"today_new"`
	TodayActive int     `json:"today_active"`
	PaidAmount  float64 `json:"paid_amount"`
	PaidOrders  int     `json:"paid_orders"`
}
type OverviewReq struct {
	g.Meta `path:"/overview" method:"get" tags:"Manage/Monitor" summary:"跨站数据汇总(上线站点)"`
}
type OverviewRes struct {
	SiteCount   int                `json:"site_count"`
	OkCount     int                `json:"ok_count"`
	UserCount   int                `json:"user_count"`
	TodayNew    int                `json:"today_new"`
	TodayActive int                `json:"today_active"`
	PaidAmount  float64            `json:"paid_amount"`
	PaidOrders  int                `json:"paid_orders"`
	Sites       []SiteOverviewItem `json:"sites"`
}

// 审计日志
type LogItem struct {
	Id        int64  `json:"id"`
	AdminId   int64  `json:"admin_id"`
	Username  string `json:"username"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Ip        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}
type LogListReq struct {
	g.Meta    `path:"/logs" method:"get" tags:"Manage/Monitor" summary:"总后台操作审计"`
	AdminId   int64  `json:"admin_id" v:"min:0#admin_id 不合法"`
	Method    string `json:"method"   v:"in:,POST,PUT,DELETE#method 仅支持 POST/PUT/DELETE"`
	Keyword   string `json:"keyword"` // path 模糊
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type LogListRes struct {
	List  []LogItem `json:"list"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

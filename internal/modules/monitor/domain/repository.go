// Package domain 监控/汇总领域层(M5)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

// LogFilter 审计日志筛选。
type LogFilter struct {
	AdminId   int64
	Method    string
	Keyword   string // path 模糊
	StartDate string // created_at >= (YYYY-MM-DD)
	EndDate   string
}

// LogItem 审计日志(联表管理员名)。
type LogItem struct {
	entity.ManageLogBase
	Username string `orm:"username"`
}

type Repository interface {
	FindSite(ctx context.Context, siteId int64) (*entity.Site, error)
	ListOnlineSites(ctx context.Context) ([]*entity.Site, error)
	Logs(ctx context.Context, f LogFilter, page, size int) ([]*LogItem, int, error)
}

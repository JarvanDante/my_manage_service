// Package domain 站点配置发布领域层(M3)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

type Repository interface {
	// 站点信息(控制面同库)
	FindSite(ctx context.Context, siteId int64) (*entity.Site, error)

	// 发布记录
	NextVersion(ctx context.Context, siteId int64) (int, error)
	CreateRecord(ctx context.Context, rec *entity.SiteConfigPublish) (int64, error)
	History(ctx context.Context, siteId int64, page, size int) ([]*entity.SiteConfigPublish, int, error)
	FindVersion(ctx context.Context, siteId int64, version int) (*entity.SiteConfigPublish, error)
}

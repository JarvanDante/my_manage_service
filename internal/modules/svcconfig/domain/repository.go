// Package domain 平台服务配置发布领域层。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

type Repository interface {
	NextVersion(ctx context.Context, service, env string) (int, error)
	CreateRecord(ctx context.Context, rec *entity.ServiceConfigPublish) (int64, error)
	History(ctx context.Context, service, env string, page, size int) ([]*entity.ServiceConfigPublish, int, error)
	FindVersion(ctx context.Context, service, env string, version int) (*entity.ServiceConfigPublish, error)
}

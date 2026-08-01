// Package domain 商户领域层(M1)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

// Filter 商户列表筛选。
type Filter struct {
	Keyword string // 名称/联系人 模糊
	Status  int    // 0全部 1启用 2停用
}

type Repository interface {
	List(ctx context.Context, f Filter, page, size int) ([]*entity.Merchant, int, error)
	FindById(ctx context.Context, id int64) (*entity.Merchant, error)
	FindByName(ctx context.Context, name string) (*entity.Merchant, error)
	Create(ctx context.Context, m *entity.Merchant) (int64, error)
	Update(ctx context.Context, m *entity.Merchant) error
	SetStatus(ctx context.Context, id int64, status int) error
}

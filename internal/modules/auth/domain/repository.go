// Package domain 总后台管理员领域层。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*entity.ManageAdmin, error)
	FindById(ctx context.Context, id int64) (*entity.ManageAdmin, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
}

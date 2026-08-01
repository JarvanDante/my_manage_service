package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	authdomain "github.com/JarvanDante/my_manage_service/internal/modules/auth/domain"
)

type authRepo struct{}

// NewAuthRepo 返回 auth 领域仓储实现。
func NewAuthRepo() authdomain.Repository { return &authRepo{} }

func (r *authRepo) FindByUsername(ctx context.Context, username string) (*entity.ManageAdmin, error) {
	var a *entity.ManageAdmin
	err := g.Model("manage_admin").Ctx(ctx).Where("username", username).Scan(&a)
	return a, err
}

func (r *authRepo) FindById(ctx context.Context, id int64) (*entity.ManageAdmin, error) {
	var a *entity.ManageAdmin
	err := g.Model("manage_admin").Ctx(ctx).Where("id", id).Scan(&a)
	return a, err
}

func (r *authRepo) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	_, err := g.Model("manage_admin").Ctx(ctx).Where("id", id).Data(g.Map{
		"last_login_at": gtime.Now(),
		"last_ip":       ip,
	}).Update()
	return err
}

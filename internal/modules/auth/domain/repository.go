// Package domain 总后台管理员领域层。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

type Repository interface {
	// 认证
	FindByUsername(ctx context.Context, username string) (*entity.ManageAdmin, error)
	FindById(ctx context.Context, id int64) (*entity.ManageAdmin, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error

	// 角色(M6)
	ListRoles(ctx context.Context) ([]*entity.ManageRole, error)
	FindRoleById(ctx context.Context, id int64) (*entity.ManageRole, error)
	FindRoleByCode(ctx context.Context, code string) (*entity.ManageRole, error)
	CreateRole(ctx context.Context, r *entity.ManageRole) (int64, error)
	UpdateRole(ctx context.Context, id int64, name, remark string, status int, permissions string) error
	DeleteRole(ctx context.Context, id int64) error
	CountAdminsByRoleId(ctx context.Context, roleId int64) (int, error)

	// 管理员账号(M6)
	ListAdmins(ctx context.Context, page, size int) ([]*entity.ManageAdmin, int, error)
	CreateAdmin(ctx context.Context, a *entity.ManageAdmin) (int64, error)
	UpdateAdmin(ctx context.Context, id int64, nickname string, roleId int64, status int, password, salt string) error
	DeleteAdmin(ctx context.Context, id int64) error

	// 权限树(RBAC 菜单+接口权限)
	ListPermissions(ctx context.Context) ([]*entity.ManagePermission, error)
	FindPermissionsByIds(ctx context.Context, ids []int64) ([]*entity.ManagePermission, error)
	FindPermissionById(ctx context.Context, id int64) (*entity.ManagePermission, error)
	CreatePermission(ctx context.Context, p *entity.ManagePermission) (int64, error)
	UpdatePermission(ctx context.Context, p *entity.ManagePermission) error
	DeletePermission(ctx context.Context, id int64) error
	CountPermissionChildren(ctx context.Context, parentId int64) (int, error)
}

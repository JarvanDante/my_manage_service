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

// ==================== 角色(M6) ====================

func (r *authRepo) ListRoles(ctx context.Context) ([]*entity.ManageRole, error) {
	var list []*entity.ManageRole
	err := g.Model("manage_role").Ctx(ctx).Order("id asc").Scan(&list)
	return list, err
}

func (r *authRepo) FindRoleById(ctx context.Context, id int64) (*entity.ManageRole, error) {
	var role *entity.ManageRole
	err := g.Model("manage_role").Ctx(ctx).Where("id", id).Scan(&role)
	return role, err
}

func (r *authRepo) FindRoleByCode(ctx context.Context, code string) (*entity.ManageRole, error) {
	var role *entity.ManageRole
	err := g.Model("manage_role").Ctx(ctx).Where("code", code).Scan(&role)
	return role, err
}

func (r *authRepo) CreateRole(ctx context.Context, role *entity.ManageRole) (int64, error) {
	res, err := g.Model("manage_role").Ctx(ctx).Data(g.Map{
		"name": role.Name, "code": role.Code, "remark": role.Remark, "status": 1, "permissions": role.Permissions,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *authRepo) UpdateRole(ctx context.Context, id int64, name, remark string, status int, permissions string) error {
	_, err := g.Model("manage_role").Ctx(ctx).Where("id", id).Data(g.Map{
		"name": name, "remark": remark, "status": status, "permissions": permissions, "updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *authRepo) DeleteRole(ctx context.Context, id int64) error {
	_, err := g.Model("manage_role").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *authRepo) CountAdminsByRoleId(ctx context.Context, roleId int64) (int, error) {
	return g.Model("manage_admin").Ctx(ctx).Where("role_id", roleId).Count()
}

// ==================== 管理员账号(M6) ====================

func (r *authRepo) ListAdmins(ctx context.Context, page, size int) ([]*entity.ManageAdmin, int, error) {
	m := g.Model("manage_admin").Ctx(ctx)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ManageAdmin
	err = m.Order("id asc").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *authRepo) CreateAdmin(ctx context.Context, a *entity.ManageAdmin) (int64, error) {
	res, err := g.Model("manage_admin").Ctx(ctx).Data(g.Map{
		"username": a.Username, "password": a.Password, "salt": "",
		"nickname": a.Nickname, "role_id": a.RoleId, "status": 1,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *authRepo) UpdateAdmin(ctx context.Context, id int64, nickname string, roleId int64, status int, password string) error {
	data := g.Map{
		"nickname": nickname, "role_id": roleId, "status": status, "updated_at": gtime.Now(),
	}
	if password != "" {
		data["password"] = password
		data["salt"] = "" // bcrypt 盐内嵌于哈希, 旧 salt 字段清空
	}
	_, err := g.Model("manage_admin").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

func (r *authRepo) DeleteAdmin(ctx context.Context, id int64) error {
	_, err := g.Model("manage_admin").Ctx(ctx).Where("id", id).Delete()
	return err
}

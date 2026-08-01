// Package logic — M6 角色/权限/账号管理实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/consts"
	"github.com/JarvanDante/my_manage_service/internal/shared/rbac"
)

// ---------- 角色 ----------

func (s *sAuth) ListRoles(ctx context.Context) ([]*service.RoleDTO, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.RoleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, &service.RoleDTO{
			Id: r.Id, Name: r.Name, Code: r.Code, Remark: r.Remark, Status: r.Status,
		})
	}
	return out, nil
}

func (s *sAuth) CreateRole(ctx context.Context, in service.RoleCreateInput) (int64, error) {
	if in.Name == "" || in.Code == "" {
		return 0, gerror.New("角色名与角色码必填")
	}
	if in.Code == consts.SuperAdminCode {
		return 0, gerror.New("superadmin 为系统保留角色码")
	}
	exist, err := s.repo.FindRoleByCode(ctx, in.Code)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("角色码已存在")
	}
	return s.repo.CreateRole(ctx, &entity.ManageRole{Name: in.Name, Code: in.Code, Remark: in.Remark})
}

func (s *sAuth) UpdateRole(ctx context.Context, in service.RoleUpdateInput) error {
	role, err := s.mustRole(ctx, in.Id)
	if err != nil {
		return err
	}
	if role.Code == consts.SuperAdminCode {
		return gerror.New("超级管理员角色不可修改")
	}
	if in.Name == "" {
		return gerror.New("角色名必填")
	}
	return s.repo.UpdateRole(ctx, in.Id, in.Name, in.Remark, in.Status)
}

func (s *sAuth) DeleteRole(ctx context.Context, id int64) error {
	role, err := s.mustRole(ctx, id)
	if err != nil {
		return err
	}
	if role.Code == consts.SuperAdminCode {
		return gerror.New("超级管理员角色不可删除")
	}
	cnt, err := s.repo.CountAdminsByRoleId(ctx, id)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return gerror.New("该角色下仍有管理员, 不能删除")
	}
	if err = s.repo.DeleteRole(ctx, id); err != nil {
		return err
	}
	// 清理该角色的 Casbin 策略
	return rbac.RemoveRolePolicies(role.Code)
}

// ---------- 权限(Casbin) ----------

func (s *sAuth) ListPerms(ctx context.Context, roleCode string) ([]*service.PermDTO, error) {
	if roleCode == "" {
		return nil, gerror.New("角色码不能为空")
	}
	rows := rbac.PermsForRole(roleCode)
	out := make([]*service.PermDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, &service.PermDTO{Path: r[0], Method: r[1]})
	}
	return out, nil
}

func (s *sAuth) AddPerm(ctx context.Context, roleCode, path, method string) error {
	if roleCode == "" || path == "" || method == "" {
		return gerror.New("角色码/路径/方法均不能为空")
	}
	if roleCode == consts.SuperAdminCode {
		return gerror.New("超级管理员无需配置权限")
	}
	_, err := rbac.AddPolicy(roleCode, path, method)
	return err
}

func (s *sAuth) RemovePerm(ctx context.Context, roleCode, path, method string) error {
	if roleCode == "" || path == "" || method == "" {
		return gerror.New("角色码/路径/方法均不能为空")
	}
	_, err := rbac.RemovePolicy(roleCode, path, method)
	return err
}

// ---------- 管理员账号 ----------

func (s *sAuth) ListAdmins(ctx context.Context, page, size int) (*service.AdminListDTO, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := s.repo.ListAdmins(ctx, page, size)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleName := make(map[int64]string, len(roles))
	for _, r := range roles {
		roleName[r.Id] = r.Name
	}
	items := make([]*service.AdminItemDTO, 0, len(list))
	for _, a := range list {
		last := ""
		if a.LastLoginAt != nil {
			last = a.LastLoginAt.String()
		}
		items = append(items, &service.AdminItemDTO{
			Id: a.Id, Username: a.Username, Nickname: a.Nickname,
			RoleId: a.RoleId, RoleName: roleName[a.RoleId], Status: a.Status, LastLoginAt: last,
		})
	}
	return &service.AdminListDTO{List: items, Total: total, Page: page, Size: size}, nil
}

func (s *sAuth) CreateAdmin(ctx context.Context, in service.AdminCreateInput) (int64, error) {
	if in.Username == "" || in.Password == "" {
		return 0, gerror.New("账号与密码必填")
	}
	if in.RoleId <= 0 {
		return 0, gerror.New("请指定角色")
	}
	exist, err := s.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("账号已存在")
	}
	if _, err = s.mustRole(ctx, in.RoleId); err != nil {
		return 0, err
	}
	salt := grand.S(8)
	pwd := gmd5.MustEncryptString(in.Password + salt)
	return s.repo.CreateAdmin(ctx, &entity.ManageAdmin{
		Username: in.Username, Password: pwd, Salt: salt, Nickname: in.Nickname, RoleId: in.RoleId,
	})
}

func (s *sAuth) UpdateAdmin(ctx context.Context, in service.AdminUpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("管理员ID无效")
	}
	a, err := s.repo.FindById(ctx, in.Id)
	if err != nil {
		return err
	}
	if a == nil {
		return gerror.New("管理员不存在")
	}
	if _, err = s.mustRole(ctx, in.RoleId); err != nil {
		return err
	}
	var pwd, salt string
	if in.Password != "" {
		salt = grand.S(8)
		pwd = gmd5.MustEncryptString(in.Password + salt)
	}
	return s.repo.UpdateAdmin(ctx, in.Id, in.Nickname, in.RoleId, in.Status, pwd, salt)
}

func (s *sAuth) DeleteAdmin(ctx context.Context, id, operatorId int64) error {
	if id <= 0 {
		return gerror.New("管理员ID无效")
	}
	if id == operatorId {
		return gerror.New("不能删除当前登录账号")
	}
	a, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return gerror.New("管理员不存在")
	}
	return s.repo.DeleteAdmin(ctx, id)
}

func (s *sAuth) mustRole(ctx context.Context, id int64) (*entity.ManageRole, error) {
	if id <= 0 {
		return nil, gerror.New("角色ID无效")
	}
	role, err := s.repo.FindRoleById(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, gerror.New("角色不存在")
	}
	return role, nil
}

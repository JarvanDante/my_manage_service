// Package controller — M6 角色/权限/账号管理控制器。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/auth/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
)

// ---------- 角色 ----------

func (c *Controller) RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	list, err := c.auth.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.RoleItem, 0, len(list))
	for _, r := range list {
		items = append(items, v1.RoleItem{Id: r.Id, Name: r.Name, Code: r.Code, Remark: r.Remark, Status: r.Status})
	}
	return &v1.RoleListRes{List: items}, nil
}

func (c *Controller) RoleCreate(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	id, err := c.auth.CreateRole(ctx, service.RoleCreateInput{Name: req.Name, Code: req.Code, Remark: req.Remark})
	if err != nil {
		return nil, err
	}
	return &v1.RoleCreateRes{Id: id}, nil
}

func (c *Controller) RoleUpdate(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	if err = c.auth.UpdateRole(ctx, service.RoleUpdateInput{Id: req.Id, Name: req.Name, Remark: req.Remark, Status: req.Status}); err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}

func (c *Controller) RoleDelete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	if err = c.auth.DeleteRole(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}

// ---------- 权限 ----------

func (c *Controller) PermList(ctx context.Context, req *v1.PermListReq) (res *v1.PermListRes, err error) {
	list, err := c.auth.ListPerms(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	items := make([]v1.PermItem, 0, len(list))
	for _, p := range list {
		items = append(items, v1.PermItem{Path: p.Path, Method: p.Method})
	}
	return &v1.PermListRes{List: items}, nil
}

func (c *Controller) PermAdd(ctx context.Context, req *v1.PermAddReq) (res *v1.PermAddRes, err error) {
	if err = c.auth.AddPerm(ctx, req.Code, req.Path, req.Method); err != nil {
		return nil, err
	}
	return &v1.PermAddRes{}, nil
}

func (c *Controller) PermDel(ctx context.Context, req *v1.PermDelReq) (res *v1.PermDelRes, err error) {
	if err = c.auth.RemovePerm(ctx, req.Code, req.Path, req.Method); err != nil {
		return nil, err
	}
	return &v1.PermDelRes{}, nil
}

// ---------- 管理员账号 ----------

func (c *Controller) AdminList(ctx context.Context, req *v1.AdminListReq) (res *v1.AdminListRes, err error) {
	dto, err := c.auth.ListAdmins(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.AdminItem, 0, len(dto.List))
	for _, a := range dto.List {
		items = append(items, v1.AdminItem{
			Id: a.Id, Username: a.Username, Nickname: a.Nickname,
			RoleId: a.RoleId, RoleName: a.RoleName, Status: a.Status, LastLoginAt: a.LastLoginAt,
		})
	}
	return &v1.AdminListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) AdminCreate(ctx context.Context, req *v1.AdminCreateReq) (res *v1.AdminCreateRes, err error) {
	id, err := c.auth.CreateAdmin(ctx, service.AdminCreateInput{
		Username: req.Username, Password: req.Password, Nickname: req.Nickname, RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminCreateRes{Id: id}, nil
}

func (c *Controller) AdminUpdate(ctx context.Context, req *v1.AdminUpdateReq) (res *v1.AdminUpdateRes, err error) {
	if err = c.auth.UpdateAdmin(ctx, service.AdminUpdateInput{
		Id: req.Id, Nickname: req.Nickname, RoleId: req.RoleId, Status: req.Status, Password: req.Password,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminUpdateRes{}, nil
}

func (c *Controller) AdminDelete(ctx context.Context, req *v1.AdminDeleteReq) (res *v1.AdminDeleteRes, err error) {
	operator, err := adminId(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.auth.DeleteAdmin(ctx, req.Id, operator); err != nil {
		return nil, err
	}
	return &v1.AdminDeleteRes{}, nil
}

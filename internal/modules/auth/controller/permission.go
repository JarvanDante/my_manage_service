// Package controller — RBAC 权限树控制器。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/auth/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
)

func (c *Controller) PermTree(ctx context.Context, req *v1.PermTreeReq) (res *v1.PermTreeRes, err error) {
	if _, err = adminId(ctx); err != nil {
		return nil, err
	}
	list, err := c.auth.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PermTreeRes{List: toPermNodes(list)}, nil
}

func toPermNodes(list []*service.PermissionDTO) []v1.PermNode {
	out := make([]v1.PermNode, 0, len(list))
	for _, p := range list {
		out = append(out, v1.PermNode{
			Id: p.Id, ParentId: p.ParentId, Name: p.Name, RouteUrl: p.RouteUrl,
			Component: p.Component, Method: p.Method, Icon: p.Icon, IsMenu: p.IsMenu,
			HideInMenu: p.HideInMenu, AffixTab: p.AffixTab, ActivePath: p.ActivePath,
			Sort: p.Sort, Status: p.Status, Children: toPermNodes(p.Children),
		})
	}
	return out
}

func (c *Controller) PermCreate(ctx context.Context, req *v1.PermCreateReq) (res *v1.PermCreateRes, err error) {
	id, err := c.auth.CreatePermission(ctx, service.PermissionInput{
		ParentId: req.ParentId, Name: req.Name, RouteUrl: req.RouteUrl, Component: req.Component,
		Method: req.Method, Icon: req.Icon, IsMenu: req.IsMenu, HideInMenu: req.HideInMenu,
		AffixTab: req.AffixTab, ActivePath: req.ActivePath, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PermCreateRes{Id: id}, nil
}

func (c *Controller) PermUpdate(ctx context.Context, req *v1.PermUpdateReq) (res *v1.PermUpdateRes, err error) {
	if err = c.auth.UpdatePermission(ctx, service.PermissionInput{
		Id: req.Id, ParentId: req.ParentId, Name: req.Name, RouteUrl: req.RouteUrl, Component: req.Component,
		Method: req.Method, Icon: req.Icon, IsMenu: req.IsMenu, HideInMenu: req.HideInMenu,
		AffixTab: req.AffixTab, ActivePath: req.ActivePath, Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.PermUpdateRes{}, nil
}

func (c *Controller) PermDelete(ctx context.Context, req *v1.PermDeleteReq) (res *v1.PermDeleteRes, err error) {
	if err = c.auth.DeletePermission(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.PermDeleteRes{}, nil
}

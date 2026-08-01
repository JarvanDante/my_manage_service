package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/auth/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
)

// Menus 返回当前管理员可见的菜单树(DB 驱动, 按角色权限过滤; 超管全量)。
func (c *Controller) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	id, err := adminId(ctx)
	if err != nil {
		return nil, err
	}
	nodes, err := c.auth.MenusForAdmin(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.MenusRes{List: toMenuNodes(nodes)}, nil
}

func toMenuNodes(list []*service.MenuNodeDTO) []v1.MenuNode {
	out := make([]v1.MenuNode, 0, len(list))
	for _, n := range list {
		out = append(out, v1.MenuNode{
			Name: n.Name, Path: n.Path, Component: n.Component, Redirect: n.Redirect,
			Meta: v1.MenuMeta{
				Title: n.Title, Icon: n.Icon, Order: n.Order,
				AffixTab: n.AffixTab, HideInMenu: n.HideInMenu, ActivePath: n.ActivePath,
			},
			Children: toMenuNodes(n.Children),
		})
	}
	return out
}

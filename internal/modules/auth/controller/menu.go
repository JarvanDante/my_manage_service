package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/auth/v1"
)

// Menus 返回当前管理员可见的菜单树(动态菜单 method B)。
// 目前返回总后台完整菜单;后续可按 adminId 的角色/Casbin 权限做过滤。
func (c *Controller) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	if _, err = adminId(ctx); err != nil {
		return nil, err
	}
	return &v1.MenusRes{List: buildManageMenus()}, nil
}

// buildManageMenus 总后台菜单树。component 为相对 apps/web-ele/src/views 的路径(不含 .vue)。
func buildManageMenus() []v1.MenuNode {
	return []v1.MenuNode{
		{
			Name:     "Dashboard",
			Path:     "/dashboard",
			Redirect: "/analytics",
			Meta:     v1.MenuMeta{Title: "仪表盘", Icon: "lucide:layout-dashboard", Order: -1},
			Children: []v1.MenuNode{
				{
					Name:      "Analytics",
					Path:      "/analytics",
					Component: "dashboard/analytics/index",
					Meta:      v1.MenuMeta{Title: "分析页", Icon: "lucide:area-chart", AffixTab: true},
				},
			},
		},
		{
			Name:      "Merchant",
			Path:      "/merchant",
			Component: "merchant/index",
			Meta:      v1.MenuMeta{Title: "商户管理", Icon: "lucide:store", Order: 10},
		},
		{
			Name:      "Site",
			Path:      "/site",
			Component: "site/index",
			Meta:      v1.MenuMeta{Title: "站点管理", Icon: "lucide:globe", Order: 20},
		},
		{
			Name:      "SiteDetail",
			Path:      "/site/detail",
			Component: "site/detail",
			Meta:      v1.MenuMeta{Title: "站点详情", HideInMenu: true, ActivePath: "/site"},
		},
		{
			Name:      "ManageLogs",
			Path:      "/monitor/logs",
			Component: "monitor/logs",
			Meta:      v1.MenuMeta{Title: "操作审计", Icon: "lucide:scroll-text", Order: 80},
		},
		{
			Name:     "System",
			Path:     "/system",
			Redirect: "/system/role",
			Meta:     v1.MenuMeta{Title: "系统管理", Icon: "lucide:settings", Order: 90},
			Children: []v1.MenuNode{
				{
					Name:      "SystemRole",
					Path:      "role",
					Component: "system/role",
					Meta:      v1.MenuMeta{Title: "角色权限", Icon: "lucide:shield-check"},
				},
				{
					Name:      "SystemAdmin",
					Path:      "admin",
					Component: "system/admin",
					Meta:      v1.MenuMeta{Title: "管理员", Icon: "lucide:users"},
				},
			},
		},
	}
}

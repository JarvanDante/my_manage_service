package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

func (r *authRepo) ListPermissions(ctx context.Context) ([]*entity.ManagePermission, error) {
	var list []*entity.ManagePermission
	err := g.Model("manage_permission").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *authRepo) FindPermissionsByIds(ctx context.Context, ids []int64) ([]*entity.ManagePermission, error) {
	var list []*entity.ManagePermission
	if len(ids) == 0 {
		return list, nil
	}
	err := g.Model("manage_permission").Ctx(ctx).WhereIn("id", ids).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *authRepo) FindPermissionById(ctx context.Context, id int64) (*entity.ManagePermission, error) {
	var p *entity.ManagePermission
	err := g.Model("manage_permission").Ctx(ctx).Where("id", id).Scan(&p)
	return p, err
}

func (r *authRepo) CreatePermission(ctx context.Context, p *entity.ManagePermission) (int64, error) {
	res, err := g.Model("manage_permission").Ctx(ctx).Data(g.Map{
		"parent_id": p.ParentId, "name": p.Name, "route_url": p.RouteUrl,
		"component": p.Component, "method": p.Method, "icon": p.Icon,
		"is_menu": p.IsMenu, "hide_in_menu": p.HideInMenu, "affix_tab": p.AffixTab,
		"active_path": p.ActivePath, "sort": p.Sort, "status": p.Status,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *authRepo) UpdatePermission(ctx context.Context, p *entity.ManagePermission) error {
	_, err := g.Model("manage_permission").Ctx(ctx).Where("id", p.Id).Data(g.Map{
		"parent_id": p.ParentId, "name": p.Name, "route_url": p.RouteUrl,
		"component": p.Component, "method": p.Method, "icon": p.Icon,
		"is_menu": p.IsMenu, "hide_in_menu": p.HideInMenu, "affix_tab": p.AffixTab,
		"active_path": p.ActivePath, "sort": p.Sort, "status": p.Status,
		"updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *authRepo) DeletePermission(ctx context.Context, id int64) error {
	_, err := g.Model("manage_permission").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *authRepo) CountPermissionChildren(ctx context.Context, parentId int64) (int, error) {
	return g.Model("manage_permission").Ctx(ctx).Where("parent_id", parentId).Count()
}

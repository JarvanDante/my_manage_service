// Package auth 总后台认证模块装配。
package auth

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/modules/auth/controller"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/logic"
)

// RegisterPublic 公开接口(登录)。
func RegisterPublic(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(controller.New(logic.New(repo)).Login)
}

// RegisterAuthed 仅需登录的接口(退出/自身信息)。
func RegisterAuthed(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(ctrl.Logout, ctrl.Info)
}

// RegisterPermManage 需权限校验的接口(角色/权限/账号管理)。
func RegisterPermManage(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(
		ctrl.RoleList, ctrl.RoleCreate, ctrl.RoleUpdate, ctrl.RoleDelete,
		ctrl.PermList, ctrl.PermAdd, ctrl.PermDel,
		ctrl.AdminList, ctrl.AdminCreate, ctrl.AdminUpdate, ctrl.AdminDelete,
	)
}

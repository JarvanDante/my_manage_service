// Package siteconfig 站点配置发布模块装配(M3)。
package siteconfig

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/controller"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/siteconfig/logic"
)

// Register 配置发布接口(挂权限组)。
func Register(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(ctrl.Current, ctrl.Publish, ctrl.History, ctrl.Rollback)
}

// Package site 站点模块装配(M2)。
package site

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/modules/site/controller"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/logic"
)

// Register 站点管理接口(挂权限组)。
func Register(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(
		ctrl.List, ctrl.Detail, ctrl.Create, ctrl.Update, ctrl.Status, ctrl.DbCheck,
		ctrl.DomainBind, ctrl.DomainUnbind, ctrl.DomainMain,
	)
}

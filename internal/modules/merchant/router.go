// Package merchant 商户模块装配(M1)。
package merchant

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/controller"
	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/logic"
)

// Register 商户管理接口(挂权限组)。
func Register(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Status)
}

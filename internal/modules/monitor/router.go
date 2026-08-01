// Package monitor 监控/汇总模块装配(M5)。
package monitor

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/controller"
	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/monitor/logic"
)

// Register 监控接口(挂权限组)。
func Register(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := controller.New(logic.New(repo))
	group.Bind(ctrl.Health, ctrl.Overview, ctrl.Logs)
}

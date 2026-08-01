// Package cmd 启动装配。依赖(dao)在此注入各模块。
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_manage_service/internal/dao"
	authmod "github.com/JarvanDante/my_manage_service/internal/modules/auth"
	merchmod "github.com/JarvanDante/my_manage_service/internal/modules/merchant"
	monmod "github.com/JarvanDante/my_manage_service/internal/modules/monitor"
	sitemod "github.com/JarvanDante/my_manage_service/internal/modules/site"
	scmod "github.com/JarvanDante/my_manage_service/internal/modules/siteconfig"
	"github.com/JarvanDante/my_manage_service/internal/shared/middleware"
)

// Main 总后台单进程入口(:8003)。
var Main = gcmd.Command{
	Name:  "manageapi",
	Brief: "总后台(控制面)API",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server()
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)

		s.Group("/manage", func(group *ghttp.RouterGroup) {
			// 公开: 登录
			authmod.RegisterPublic(group, dao.NewAuthRepo())
			// 仅需登录: 退出 / 查看自身信息(任何管理员)
			group.Group("/", func(authed *ghttp.RouterGroup) {
				authed.Middleware(middleware.ManageAuth, middleware.ManageOpLog)
				authmod.RegisterAuthed(authed, dao.NewAuthRepo())
			})
			// 需权限校验: 业务接口(M1 商户起挂这里)
			group.Group("/", func(perm *ghttp.RouterGroup) {
				perm.Middleware(middleware.ManageAuth, middleware.ManagePerm, middleware.ManageOpLog)
				authmod.RegisterPermManage(perm, dao.NewAuthRepo())
				merchmod.Register(perm, dao.NewMerchantRepo())
				sitemod.Register(perm, dao.NewSiteRepo())
				scmod.Register(perm, dao.NewSiteConfigRepo())
				monmod.Register(perm, dao.NewMonitorRepo())
			})
		})

		s.Run()
		return nil
	},
}

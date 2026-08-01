// Package middleware 全局/分组中间件。
package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/shared/consts"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
	"github.com/JarvanDante/my_manage_service/internal/shared/rbac"
)

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// ManageAuth 总后台鉴权: 解析 token, 成功则把 adminId 写入 ctx。
func ManageAuth(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.Get("token").String()
	}
	adminId, err := kit.ParseToken(r.GetCtx(), token)
	if err != nil || adminId <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未登录或登录已失效"))
		return
	}
	r.SetCtxVar(consts.CtxAdminId, adminId)
	r.Middleware.Next()
}

// ManagePerm 权限校验(需在 ManageAuth 之后): 超管放行, 其余角色走 Casbin。
func ManagePerm(r *ghttp.Request) {
	adminId := r.GetCtxVar(consts.CtxAdminId).Int64()
	if adminId <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未登录"))
		return
	}
	val, err := g.Model("manage_admin ma").Ctx(r.GetCtx()).
		LeftJoin("manage_role mr", "mr.id=ma.role_id").
		Where("ma.id", adminId).
		Fields("mr.code").
		Value()
	if err != nil {
		r.SetError(err)
		return
	}
	roleCode := val.String()
	if roleCode == consts.SuperAdminCode {
		r.Middleware.Next()
		return
	}
	ok, err := rbac.Enforce(roleCode, r.URL.Path, r.Method)
	if err != nil {
		r.SetError(err)
		return
	}
	if !ok {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "无权限访问该资源"))
		return
	}
	r.Middleware.Next()
}

// ManageOpLog 写操作审计(POST/PUT/DELETE)。
func ManageOpLog(r *ghttp.Request) {
	r.Middleware.Next()
	m := r.Method
	if m == "POST" || m == "PUT" || m == "DELETE" {
		adminId := r.GetCtxVar(consts.CtxAdminId).Int64()
		if adminId > 0 {
			_, _ = g.Model("manage_log").Ctx(r.GetCtx()).Data(g.Map{
				"admin_id": adminId,
				"method":   m,
				"path":     r.URL.Path,
				"ip":       r.GetClientIp(),
			}).Insert()
		}
	}
}

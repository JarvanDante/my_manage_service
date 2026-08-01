// Package middleware 全局/分组中间件。
package middleware

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/shared/consts"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
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

// ManagePerm 权限校验(需在 ManageAuth 之后): 超管放行;
// 其余角色按 role.permissions 里勾选的「接口权限(is_menu=0)」匹配请求路径+方法。
// 无权限返回 403(区别于未登录的 61), 避免前端把无权限当成掉登录。
func ManagePerm(r *ghttp.Request) {
	ctx := r.GetCtx()
	adminId := r.GetCtxVar(consts.CtxAdminId).Int64()
	if adminId <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未登录"))
		return
	}
	rec, err := g.Model("manage_admin ma").Ctx(ctx).
		LeftJoin("manage_role mr", "mr.id=ma.role_id").
		Where("ma.id", adminId).
		Fields("mr.code AS code, mr.permissions AS permissions").
		One()
	if err != nil {
		r.SetError(err)
		return
	}
	if rec == nil {
		r.SetError(forbidden())
		return
	}
	if rec["code"].String() == consts.SuperAdminCode {
		r.Middleware.Next()
		return
	}
	ids := parsePermIds(rec["permissions"].String())
	if len(ids) > 0 {
		var perms []*entity.ManagePermission
		_ = g.Model("manage_permission").Ctx(ctx).
			WhereIn("id", ids).Where("is_menu", 0).Where("status", 1).
			Scan(&perms)
		for _, p := range perms {
			if methodMatch(p.Method, r.Method) && pathMatch(p.RouteUrl, r.URL.Path) {
				r.Middleware.Next()
				return
			}
		}
	}
	r.SetError(forbidden())
}

func forbidden() error {
	return gerror.NewCode(gcode.New(403, "无权限访问该资源", nil), "无权限访问该资源")
}

func parsePermIds(s string) []int64 {
	out := make([]int64, 0)
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			out = append(out, v)
		}
	}
	return out
}

func methodMatch(permMethod, reqMethod string) bool {
	permMethod = strings.TrimSpace(permMethod)
	return permMethod == "" || permMethod == "*" || strings.EqualFold(permMethod, reqMethod)
}

// pathMatch 支持 {param} 通配: /manage/sites/{id}/domains/{domainId}
func pathMatch(pattern, path string) bool {
	if pattern == path {
		return true
	}
	if !strings.Contains(pattern, "{") {
		return false
	}
	var re strings.Builder
	re.WriteString("^")
	rest := pattern
	for {
		i := strings.Index(rest, "{")
		if i < 0 {
			re.WriteString(regexp.QuoteMeta(rest))
			break
		}
		j := strings.Index(rest, "}")
		if j < 0 || j < i {
			re.WriteString(regexp.QuoteMeta(rest))
			break
		}
		re.WriteString(regexp.QuoteMeta(rest[:i]))
		re.WriteString("[^/]+")
		rest = rest[j+1:]
	}
	re.WriteString("$")
	ok, _ := regexp.MatchString(re.String(), path)
	return ok
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

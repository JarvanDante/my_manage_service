// Package v1 总后台角色/权限/账号管理接口契约(M6)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 角色 ----------

type RoleItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
	Permissions string `json:"permissions"`
}

type RoleListReq struct {
	g.Meta `path:"/roles" method:"get" tags:"Manage/Role" summary:"角色列表"`
}
type RoleListRes struct {
	List []RoleItem `json:"list"`
}

type RoleCreateReq struct {
	g.Meta      `path:"/roles" method:"post" tags:"Manage/Role" summary:"创建角色"`
	Name        string `json:"name" v:"required#角色名必填"`
	Code        string `json:"code" v:"required#角色码必填"`
	Remark      string `json:"remark"`
	Permissions string `json:"permissions"`
}
type RoleCreateRes struct {
	Id int64 `json:"id"`
}

type RoleUpdateReq struct {
	g.Meta      `path:"/roles/{id}" method:"put" tags:"Manage/Role" summary:"更新角色"`
	Id          int64  `json:"id"   v:"required|min:1#角色ID必填|角色ID必须大于0"`
	Name        string `json:"name" v:"required#角色名必填"`
	Remark      string `json:"remark"`
	Status      int    `json:"status" v:"in:0,1#status 仅支持 0/1"`
	Permissions string `json:"permissions"`
}
type RoleUpdateRes struct{}

type RoleDeleteReq struct {
	g.Meta `path:"/roles/{id}" method:"delete" tags:"Manage/Role" summary:"删除角色"`
	Id     int64 `json:"id" v:"required|min:1#角色ID必填|角色ID必须大于0"`
}
type RoleDeleteRes struct{}

// ---------- 管理员账号 ----------

type AdminItem struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	RoleId      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	Status      int    `json:"status"`
	LastLoginAt string `json:"last_login_at"`
}

type AdminListReq struct {
	g.Meta `path:"/admins" method:"get" tags:"Manage/Admin" summary:"管理员列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type AdminListRes struct {
	List  []AdminItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type AdminCreateReq struct {
	g.Meta   `path:"/admins" method:"post" tags:"Manage/Admin" summary:"创建管理员"`
	Username string `json:"username" v:"required#账号必填"`
	Password string `json:"password" v:"required#密码必填"`
	Nickname string `json:"nickname"`
	RoleId   int64  `json:"role_id"  v:"required|min:1#角色必填|角色不合法"`
}
type AdminCreateRes struct {
	Id int64 `json:"id"`
}

type AdminUpdateReq struct {
	g.Meta   `path:"/admins/{id}" method:"put" tags:"Manage/Admin" summary:"更新管理员(password空=不改)"`
	Id       int64  `json:"id"       v:"required|min:1#管理员ID必填|管理员ID必须大于0"`
	Nickname string `json:"nickname"`
	RoleId   int64  `json:"role_id"  v:"required|min:1#角色必填|角色不合法"`
	Status   int    `json:"status"   v:"in:0,1#status 仅支持 0/1"`
	Password string `json:"password"`
}
type AdminUpdateRes struct{}

type AdminDeleteReq struct {
	g.Meta `path:"/admins/{id}" method:"delete" tags:"Manage/Admin" summary:"删除管理员"`
	Id     int64 `json:"id" v:"required|min:1#管理员ID必填|管理员ID必须大于0"`
}
type AdminDeleteRes struct{}

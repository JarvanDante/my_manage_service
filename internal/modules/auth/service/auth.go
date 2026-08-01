// Package service 总后台管理员对外接口。
package service

import "context"

type LoginInput struct {
	Username string
	Password string
	Ip       string
}

type AdminInfoDTO struct {
	Id       int64
	Username string
	Nickname string
	RoleId   int64
}

type LoginDTO struct {
	Token string
	Admin *AdminInfoDTO
}

// ---- M6 角色 / 权限 / 账号 ----

type RoleDTO struct {
	Id     int64
	Name   string
	Code   string
	Remark string
	Status int
}

type PermDTO struct {
	Path   string
	Method string
}

type RoleCreateInput struct {
	Name   string
	Code   string
	Remark string
}

type RoleUpdateInput struct {
	Id     int64
	Name   string
	Remark string
	Status int
}

type AdminItemDTO struct {
	Id          int64
	Username    string
	Nickname    string
	RoleId      int64
	RoleName    string
	Status      int
	LastLoginAt string
}

type AdminListDTO struct {
	List  []*AdminItemDTO
	Total int
	Page  int
	Size  int
}

type AdminCreateInput struct {
	Username string
	Password string
	Nickname string
	RoleId   int64
}

type AdminUpdateInput struct {
	Id       int64
	Nickname string
	RoleId   int64
	Status   int
	Password string // 空表示不改
}

type IAuth interface {
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	Logout(ctx context.Context, adminId int64) error
	Info(ctx context.Context, adminId int64) (*AdminInfoDTO, error)

	// M6 角色权限管理
	ListRoles(ctx context.Context) ([]*RoleDTO, error)
	CreateRole(ctx context.Context, in RoleCreateInput) (int64, error)
	UpdateRole(ctx context.Context, in RoleUpdateInput) error
	DeleteRole(ctx context.Context, id int64) error
	ListPerms(ctx context.Context, roleCode string) ([]*PermDTO, error)
	AddPerm(ctx context.Context, roleCode, path, method string) error
	RemovePerm(ctx context.Context, roleCode, path, method string) error

	// M6 管理员账号管理
	ListAdmins(ctx context.Context, page, size int) (*AdminListDTO, error)
	CreateAdmin(ctx context.Context, in AdminCreateInput) (int64, error)
	UpdateAdmin(ctx context.Context, in AdminUpdateInput) error
	DeleteAdmin(ctx context.Context, id, operatorId int64) error
}

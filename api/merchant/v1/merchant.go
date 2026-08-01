// Package v1 总后台商户管理接口契约(M1)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type MerchantItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Contact   string `json:"contact"`
	Phone     string `json:"phone"`
	Status    int    `json:"status"` // 1启用 0停用
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}

// 商户列表
type ListReq struct {
	g.Meta  `path:"/merchants" method:"get" tags:"Manage/Merchant" summary:"商户列表"`
	Keyword string `json:"keyword"`                        // 名称/联系人模糊
	Status  int    `json:"status" v:"in:0,1,2#状态仅支持0/1/2"` // 0全部 1启用 2停用
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []MerchantItem `json:"list"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// 创建商户
type CreateReq struct {
	g.Meta  `path:"/merchants" method:"post" tags:"Manage/Merchant" summary:"创建商户"`
	Name    string `json:"name" v:"required#商户名必填"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Remark  string `json:"remark"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// 更新商户
type UpdateReq struct {
	g.Meta  `path:"/merchants/{id}" method:"put" tags:"Manage/Merchant" summary:"更新商户"`
	Id      int64  `json:"id"   v:"required|min:1#商户ID必填|商户ID必须大于0"`
	Name    string `json:"name" v:"required#商户名必填"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Remark  string `json:"remark"`
}
type UpdateRes struct{}

// 启用/停用
type StatusReq struct {
	g.Meta `path:"/merchants/{id}/status" method:"post" tags:"Manage/Merchant" summary:"启用/停用商户"`
	Id     int64 `json:"id"     v:"required|min:1#商户ID必填|商户ID必须大于0"`
	Status int   `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type StatusRes struct{}

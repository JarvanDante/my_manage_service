// Package v1 总后台站点管理接口契约(M2)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type DomainItem struct {
	Id     int64  `json:"id"`
	Domain string `json:"domain"`
	IsMain int    `json:"is_main"`
	Https  int    `json:"https"`
	Status int    `json:"status"`
}

type SiteItem struct {
	Id           int64  `json:"id"`
	MerchantId   int64  `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	SiteCode     string `json:"site_code"`
	Name         string `json:"name"`
	Env          string `json:"env"`
	Status       int    `json:"status"` // 0筹备 1上线 2停用
	DbHost       string `json:"db_host"`
	DbPort       int    `json:"db_port"`
	DbName       string `json:"db_name"`
	DbUser       string `json:"db_user"`
	DbPass       string `json:"db_pass"` // 永远脱敏 ******
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
}

// 站点列表
type ListReq struct {
	g.Meta     `path:"/sites" method:"get" tags:"Manage/Site" summary:"站点列表"`
	Keyword    string `json:"keyword"`
	MerchantId int64  `json:"merchant_id" v:"min:0#merchant_id 不合法"`
	Env        string `json:"env"    v:"in:,dev,test,prod#env 仅支持 dev/test/prod"`
	Status     int    `json:"status" d:"-1" v:"in:-1,0,1,2#status 仅支持 -1/0/1/2"` // -1全部
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}
type ListRes struct {
	List  []SiteItem `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// 站点详情(含域名)
type DetailReq struct {
	g.Meta `path:"/sites/{id}" method:"get" tags:"Manage/Site" summary:"站点详情"`
	Id     int64 `json:"id" v:"required|min:1#站点ID必填|站点ID必须大于0"`
}
type DetailRes struct {
	SiteItem
	Domains []DomainItem `json:"domains"`
}

// 创建站点
type CreateReq struct {
	g.Meta     `path:"/sites" method:"post" tags:"Manage/Site" summary:"创建站点"`
	MerchantId int64  `json:"merchant_id" v:"required|min:1#商户必填|商户不合法"`
	SiteCode   string `json:"site_code"   v:"required#site_code必填"`
	Name       string `json:"name"        v:"required#站点名必填"`
	Env        string `json:"env"         v:"required|in:dev,test,prod#env必填|env 仅支持 dev/test/prod"`
	DbHost     string `json:"db_host"     v:"required#db_host必填"`
	DbPort     int    `json:"db_port"     v:"required|between:1,65535#db_port必填|db_port不合法"`
	DbName     string `json:"db_name"     v:"required#db_name必填"`
	DbUser     string `json:"db_user"     v:"required#db_user必填"`
	DbPass     string `json:"db_pass"     v:"required#db_pass必填"`
	Remark     string `json:"remark"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// 更新站点(site_code 不可改; db_pass 空=不改)
type UpdateReq struct {
	g.Meta     `path:"/sites/{id}" method:"put" tags:"Manage/Site" summary:"更新站点"`
	Id         int64  `json:"id"          v:"required|min:1#站点ID必填|站点ID必须大于0"`
	MerchantId int64  `json:"merchant_id" v:"required|min:1#商户必填|商户不合法"`
	Name       string `json:"name"        v:"required#站点名必填"`
	Env        string `json:"env"         v:"required|in:dev,test,prod#env必填|env 仅支持 dev/test/prod"`
	DbHost     string `json:"db_host"     v:"required#db_host必填"`
	DbPort     int    `json:"db_port"     v:"required|between:1,65535#db_port必填|db_port不合法"`
	DbName     string `json:"db_name"     v:"required#db_name必填"`
	DbUser     string `json:"db_user"     v:"required#db_user必填"`
	DbPass     string `json:"db_pass"` // 空=不改
	Remark     string `json:"remark"`
}
type UpdateRes struct{}

// 站点状态
type StatusReq struct {
	g.Meta `path:"/sites/{id}/status" method:"post" tags:"Manage/Site" summary:"站点状态(0筹备/1上线/2停用)"`
	Id     int64 `json:"id"     v:"required|min:1#站点ID必填|站点ID必须大于0"`
	Status int   `json:"status" v:"in:0,1,2#status 仅支持 0/1/2"`
}
type StatusRes struct{}

// DB 连通性校验
type DbCheckReq struct {
	g.Meta `path:"/sites/{id}/db-check" method:"post" tags:"Manage/Site" summary:"站点库连通性校验"`
	Id     int64 `json:"id" v:"required|min:1#站点ID必填|站点ID必须大于0"`
}
type DbCheckRes struct {
	Ok        bool   `json:"ok"`
	LatencyMs int64  `json:"latency_ms"`
	Message   string `json:"message"`
}

// 绑定域名
type DomainBindReq struct {
	g.Meta `path:"/sites/{id}/domains" method:"post" tags:"Manage/Site" summary:"绑定域名"`
	Id     int64  `json:"id"     v:"required|min:1#站点ID必填|站点ID必须大于0"`
	Domain string `json:"domain" v:"required#域名必填"`
	Https  int    `json:"https"  v:"in:0,1#https 仅支持 0/1"`
}
type DomainBindRes struct {
	DomainId int64 `json:"domain_id"`
}

// 解绑域名
type DomainUnbindReq struct {
	g.Meta   `path:"/sites/{id}/domains/{domainId}" method:"delete" tags:"Manage/Site" summary:"解绑域名"`
	Id       int64 `json:"id"       v:"required|min:1#站点ID必填|站点ID必须大于0"`
	DomainId int64 `json:"domainId" v:"required|min:1#域名ID必填|域名ID必须大于0"`
}
type DomainUnbindRes struct{}

// 设主域名
type DomainMainReq struct {
	g.Meta   `path:"/sites/{id}/domains/{domainId}/main" method:"post" tags:"Manage/Site" summary:"设为主域名"`
	Id       int64 `json:"id"       v:"required|min:1#站点ID必填|站点ID必须大于0"`
	DomainId int64 `json:"domainId" v:"required|min:1#域名ID必填|域名ID必须大于0"`
}
type DomainMainRes struct{}

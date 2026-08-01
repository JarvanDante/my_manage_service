// Package v1 站点配置发布接口契约(M3)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// 当前配置(Nacos 上的实际内容, DB密码脱敏)
type CurrentReq struct {
	g.Meta `path:"/sites/{id}/config" method:"get" tags:"Manage/SiteConfig" summary:"查看站点当前Nacos配置"`
	Id     int64 `json:"id" v:"required|min:1#站点ID必填|站点ID必须大于0"`
}
type CurrentRes struct {
	SiteCode    string `json:"site_code"`
	Env         string `json:"env"`
	NamespaceId string `json:"namespace_id"`
	DataId      string `json:"data_id"`
	Group       string `json:"group"`
	Exists      bool   `json:"exists"`
	Content     string `json:"content"`
}

// 发布配置
type PublishReq struct {
	g.Meta  `path:"/sites/{id}/config/publish" method:"post" tags:"Manage/SiteConfig" summary:"发布站点配置到Nacos"`
	Id      int64  `json:"id"      v:"required|min:1#站点ID必填|站点ID必须大于0"`
	Content string `json:"content" v:"required#配置内容必填"` // YAML 原文; database 段发布时自动注入
	Remark  string `json:"remark"`
}
type PublishRes struct {
	Version int    `json:"version"`
	DataId  string `json:"data_id"`
}

// 发布历史
type HistoryReq struct {
	g.Meta `path:"/sites/{id}/config/history" method:"get" tags:"Manage/SiteConfig" summary:"发布历史"`
	Id     int64 `json:"id" v:"required|min:1#站点ID必填|站点ID必须大于0"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type HistoryItem struct {
	Id        int64  `json:"id"`
	Version   int    `json:"version"`
	Env       string `json:"env"`
	Content   string `json:"content"`
	Operator  int64  `json:"operator"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}
type HistoryRes struct {
	List  []HistoryItem `json:"list"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// 回滚
type RollbackReq struct {
	g.Meta  `path:"/sites/{id}/config/rollback" method:"post" tags:"Manage/SiteConfig" summary:"回滚到历史版本(生成新版本)"`
	Id      int64 `json:"id"      v:"required|min:1#站点ID必填|站点ID必须大于0"`
	Version int   `json:"version" v:"required|min:1#版本号必填|版本号必须大于0"`
}
type RollbackRes struct {
	Version int    `json:"version"`
	DataId  string `json:"data_id"`
}

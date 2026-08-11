// Package v1 平台服务配置发布接口契约(每个服务一份 dataId=<binary>.yaml, 命名空间按环境)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// Current 查看服务在某环境的当前 Nacos 配置。
type CurrentReq struct {
	g.Meta  `path:"/paas/services/{service}/config" method:"get" tags:"Manage/SvcConfig" summary:"查看服务当前Nacos配置"`
	Service string `json:"service" in:"path" v:"required#服务标识必填"`
	Env     string `json:"env"` // dev/test/prod, 空默认 dev
}
type CurrentRes struct {
	Service     string `json:"service"`
	Binary      string `json:"binary"`
	Env         string `json:"env"`
	NamespaceId string `json:"namespace_id"`
	DataId      string `json:"data_id"`
	Group       string `json:"group"`
	Exists      bool   `json:"exists"`
	Content     string `json:"content"`
}

// Publish 发布配置到 Nacos(YAML 原文, 校验合法后原样发布)。
type PublishReq struct {
	g.Meta  `path:"/paas/services/{service}/config/publish" method:"post" tags:"Manage/SvcConfig" summary:"发布服务配置到Nacos"`
	Service string `json:"service" in:"path" v:"required#服务标识必填"`
	Env     string `json:"env"`
	Content string `json:"content" v:"required#配置内容必填"`
	Remark  string `json:"remark"`
}
type PublishRes struct {
	Version int    `json:"version"`
	DataId  string `json:"data_id"`
}

// History 发布历史(按 服务+环境)。
type HistoryReq struct {
	g.Meta  `path:"/paas/services/{service}/config/history" method:"get" tags:"Manage/SvcConfig" summary:"服务配置发布历史"`
	Service string `json:"service" in:"path" v:"required#服务标识必填"`
	Env     string `json:"env"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
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

// Rollback 回滚到历史版本(重新发布该版本原文, 生成新版本号)。
type RollbackReq struct {
	g.Meta  `path:"/paas/services/{service}/config/rollback" method:"post" tags:"Manage/SvcConfig" summary:"回滚服务配置到历史版本"`
	Service string `json:"service" in:"path" v:"required#服务标识必填"`
	Env     string `json:"env"`
	Version int    `json:"version" v:"required|min:1#版本号必填|版本号必须大于0"`
}
type RollbackRes struct {
	Version int    `json:"version"`
	DataId  string `json:"data_id"`
}

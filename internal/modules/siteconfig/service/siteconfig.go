// Package service 站点配置发布对外接口(M3)。
package service

import "context"

type CurrentConfigDTO struct {
	SiteCode    string
	Env         string
	NamespaceId string
	DataId      string
	Group       string
	Content     string // 已脱敏(DB密码)
	Exists      bool
}

type PublishInput struct {
	SiteId   int64
	Content  string // 配置原文(database 段发布时自动注入, 无需填写)
	Operator int64
	Remark   string
}

type PublishDTO struct {
	Version int
	DataId  string
}

type HistoryItem struct {
	Id        int64
	Version   int
	Env       string
	Content   string
	Operator  int64
	Remark    string
	CreatedAt string
}

type HistoryDTO struct {
	List  []*HistoryItem
	Total int
	Page  int
	Size  int
}

type ISiteConfig interface {
	Current(ctx context.Context, siteId int64) (*CurrentConfigDTO, error)
	Publish(ctx context.Context, in PublishInput) (*PublishDTO, error)
	History(ctx context.Context, siteId int64, page, size int) (*HistoryDTO, error)
	Rollback(ctx context.Context, siteId int64, version int, operator int64) (*PublishDTO, error)
}

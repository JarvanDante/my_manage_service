// Package service 平台服务配置发布对外接口。
package service

import "context"

type CurrentConfigDTO struct {
	Service     string
	Binary      string
	Env         string
	NamespaceId string
	DataId      string
	Group       string
	Content     string
	Exists      bool
}

type PublishInput struct {
	Service  string
	Env      string
	Content  string
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

type ISvcConfig interface {
	Current(ctx context.Context, service, env string) (*CurrentConfigDTO, error)
	Publish(ctx context.Context, in PublishInput) (*PublishDTO, error)
	History(ctx context.Context, service, env string, page, size int) (*HistoryDTO, error)
	Rollback(ctx context.Context, service, env string, version int, operator int64) (*PublishDTO, error)
}

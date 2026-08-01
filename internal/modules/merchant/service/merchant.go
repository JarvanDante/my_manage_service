// Package service 商户模块对外接口(M1)。
package service

import "context"

type MerchantDTO struct {
	Id        int64
	Name      string
	Contact   string
	Phone     string
	Status    int
	Remark    string
	CreatedAt string
}

type ListInput struct {
	Keyword string
	Status  int // 0全部 1启用 2停用
	Page    int
	Size    int
}

type ListDTO struct {
	List  []*MerchantDTO
	Total int
	Page  int
	Size  int
}

type SaveInput struct {
	Id      int64 // 更新时用
	Name    string
	Contact string
	Phone   string
	Remark  string
}

type IMerchant interface {
	List(ctx context.Context, in ListInput) (*ListDTO, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	SetStatus(ctx context.Context, id int64, status int) error
}

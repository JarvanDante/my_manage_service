// Package controller 总后台商户控制器(M1)。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/merchant/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/service"
)

type Controller struct{ merch service.IMerchant }

func New(svc service.IMerchant) *Controller { return &Controller{merch: svc} }

// List 商户列表。
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.merch.List(ctx, service.ListInput{
		Keyword: req.Keyword, Status: req.Status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.MerchantItem, 0, len(dto.List))
	for _, m := range dto.List {
		items = append(items, v1.MerchantItem{
			Id: m.Id, Name: m.Name, Contact: m.Contact, Phone: m.Phone,
			Status: m.Status, Remark: m.Remark, CreatedAt: m.CreatedAt,
		})
	}
	return &v1.ListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

// Create 创建商户。
func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.merch.Create(ctx, service.SaveInput{
		Name: req.Name, Contact: req.Contact, Phone: req.Phone, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

// Update 更新商户。
func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.merch.Update(ctx, service.SaveInput{
		Id: req.Id, Name: req.Name, Contact: req.Contact, Phone: req.Phone, Remark: req.Remark,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Status 启用/停用。
func (c *Controller) Status(ctx context.Context, req *v1.StatusReq) (res *v1.StatusRes, err error) {
	if err = c.merch.SetStatus(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &v1.StatusRes{}, nil
}

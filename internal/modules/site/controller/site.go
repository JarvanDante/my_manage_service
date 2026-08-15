// Package controller 总后台站点控制器(M2)。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_manage_service/api/site/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/service"
)

type Controller struct{ site service.ISite }

func New(svc service.ISite) *Controller { return &Controller{site: svc} }

// List 站点列表。
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.site.List(ctx, service.ListInput{
		Keyword: req.Keyword, MerchantId: req.MerchantId, Env: req.Env,
		Status: req.Status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.SiteItem, 0, len(dto.List))
	for _, s := range dto.List {
		items = append(items, toItem(s))
	}
	return &v1.ListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

// Detail 站点详情(含域名)。
func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	d, err := c.site.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	domains := make([]v1.DomainItem, 0, len(d.Domains))
	for _, dm := range d.Domains {
		domains = append(domains, v1.DomainItem{
			Id: dm.Id, Domain: dm.Domain, Kind: dm.Kind, IsMain: dm.IsMain, Https: dm.Https, Status: dm.Status,
		})
	}
	return &v1.DetailRes{SiteItem: toItem(&d.SiteDTO), Domains: domains}, nil
}

// Create 创建站点(同时发放 PaaS APPKEY/APPSECRET, secret 仅本次明文返回)。
func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	out, err := c.site.Create(ctx, service.SaveInput{
		MerchantId: req.MerchantId, SiteCode: req.SiteCode, Name: req.Name, Env: req.Env,
		DbHost: req.DbHost, DbPort: req.DbPort, DbName: req.DbName,
		DbUser: req.DbUser, DbPass: req.DbPass, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: out.Id, AppKey: out.AppKey, AppSecret: out.AppSecret}, nil
}

// RevealSecret 查看 APPSECRET(AES 解密回显)。
func (c *Controller) RevealSecret(ctx context.Context, req *v1.RevealSecretReq) (res *v1.RevealSecretRes, err error) {
	out, err := c.site.RevealSecret(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.RevealSecretRes{AppKey: out.AppKey, AppSecret: out.AppSecret}, nil
}

// ResetSecret 重置 APPSECRET。
func (c *Controller) ResetSecret(ctx context.Context, req *v1.ResetSecretReq) (res *v1.ResetSecretRes, err error) {
	out, err := c.site.ResetSecret(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ResetSecretRes{AppKey: out.AppKey, AppSecret: out.AppSecret}, nil
}

// Update 更新站点。
func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.site.Update(ctx, service.SaveInput{
		Id: req.Id, MerchantId: req.MerchantId, Name: req.Name, Env: req.Env,
		DbHost: req.DbHost, DbPort: req.DbPort, DbName: req.DbName,
		DbUser: req.DbUser, DbPass: req.DbPass, Remark: req.Remark,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Status 站点状态。
func (c *Controller) Status(ctx context.Context, req *v1.StatusReq) (res *v1.StatusRes, err error) {
	if err = c.site.SetStatus(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &v1.StatusRes{}, nil
}

// DbCheck 站点库连通性校验。
func (c *Controller) DbCheck(ctx context.Context, req *v1.DbCheckReq) (res *v1.DbCheckRes, err error) {
	d, err := c.site.DbCheck(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DbCheckRes{Ok: d.Ok, LatencyMs: d.LatencyMs, Message: d.Message}, nil
}

// DomainBind 绑定域名。
func (c *Controller) DomainBind(ctx context.Context, req *v1.DomainBindReq) (res *v1.DomainBindRes, err error) {
	id, err := c.site.BindDomain(ctx, req.Id, req.Domain, req.Kind, req.Https)
	if err != nil {
		return nil, err
	}
	return &v1.DomainBindRes{DomainId: id}, nil
}

// DomainUnbind 解绑域名。
func (c *Controller) DomainUnbind(ctx context.Context, req *v1.DomainUnbindReq) (res *v1.DomainUnbindRes, err error) {
	if err = c.site.UnbindDomain(ctx, req.Id, req.DomainId); err != nil {
		return nil, err
	}
	return &v1.DomainUnbindRes{}, nil
}

// DomainMain 设主域名。
func (c *Controller) DomainMain(ctx context.Context, req *v1.DomainMainReq) (res *v1.DomainMainRes, err error) {
	if err = c.site.SetMainDomain(ctx, req.Id, req.DomainId); err != nil {
		return nil, err
	}
	return &v1.DomainMainRes{}, nil
}

func toItem(s *service.SiteDTO) v1.SiteItem {
	return v1.SiteItem{
		Id: s.Id, MerchantId: s.MerchantId, MerchantName: s.MerchantName,
		SiteCode: s.SiteCode, Name: s.Name, Env: s.Env, Status: s.Status,
		DbHost: s.DbHost, DbPort: s.DbPort, DbName: s.DbName, DbUser: s.DbUser,
		DbPass: s.DbPassMask, AppKey: s.AppKey, Remark: s.Remark, CreatedAt: s.CreatedAt,
	}
}

// ProvisionCheck 开站清单校验(M4)。
func (c *Controller) ProvisionCheck(ctx context.Context, req *v1.ProvisionCheckReq) (res *v1.ProvisionCheckRes, err error) {
	d, err := c.site.ProvisionCheck(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	items := make([]v1.ProvisionItem, 0, len(d.Items))
	for _, it := range d.Items {
		items = append(items, v1.ProvisionItem{Key: it.Key, Name: it.Name, Ok: it.Ok, Detail: it.Detail})
	}
	return &v1.ProvisionCheckRes{AllPassed: d.AllPassed, Items: items}, nil
}

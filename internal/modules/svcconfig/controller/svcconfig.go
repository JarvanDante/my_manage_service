// Package controller 平台服务配置发布控制器。
package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_manage_service/api/svcconfig/v1"
	"github.com/JarvanDante/my_manage_service/internal/modules/svcconfig/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/consts"
)

type Controller struct{ sc service.ISvcConfig }

func New(svc service.ISvcConfig) *Controller { return &Controller{sc: svc} }

func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Current(ctx context.Context, req *v1.CurrentReq) (res *v1.CurrentRes, err error) {
	d, err := c.sc.Current(ctx, req.Service, req.Env)
	if err != nil {
		return nil, err
	}
	return &v1.CurrentRes{
		Service: d.Service, Binary: d.Binary, Env: d.Env, NamespaceId: d.NamespaceId,
		DataId: d.DataId, Group: d.Group, Exists: d.Exists, Content: d.Content,
	}, nil
}

func (c *Controller) Publish(ctx context.Context, req *v1.PublishReq) (res *v1.PublishRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.sc.Publish(ctx, service.PublishInput{
		Service: req.Service, Env: req.Env, Content: req.Content, Operator: op, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PublishRes{Version: d.Version, DataId: d.DataId}, nil
}

func (c *Controller) History(ctx context.Context, req *v1.HistoryReq) (res *v1.HistoryRes, err error) {
	d, err := c.sc.History(ctx, req.Service, req.Env, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.HistoryItem, 0, len(d.List))
	for _, h := range d.List {
		items = append(items, v1.HistoryItem{
			Id: h.Id, Version: h.Version, Env: h.Env, Content: h.Content,
			Operator: h.Operator, Remark: h.Remark, CreatedAt: h.CreatedAt,
		})
	}
	return &v1.HistoryRes{List: items, Total: d.Total, Page: d.Page, Size: d.Size}, nil
}

func (c *Controller) Rollback(ctx context.Context, req *v1.RollbackReq) (res *v1.RollbackRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.sc.Rollback(ctx, req.Service, req.Env, req.Version, op)
	if err != nil {
		return nil, err
	}
	return &v1.RollbackRes{Version: d.Version, DataId: d.DataId}, nil
}

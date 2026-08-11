// Package logic 平台服务配置发布业务实现。
// 与站点配置不同: 服务配置是基础设施配置, 原样发布(不注入 database/paas), 回滚=重发历史版本原文。
package logic

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/svcconfig/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/svcconfig/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/nacospub"
)

// catalog 服务标识 -> 二进制名(即 Nacos dataId 前缀)。
var catalog = map[string]string{
	"media":   "my_media",
	"storage": "my_storage",
	"play":    "my_play",
	"ad":      "my_ad",
}

// allowedEnv 允许发布的环境。
var allowedEnv = map[string]bool{"dev": true, "test": true, "prod": true}

type sSvcConfig struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.ISvcConfig { return &sSvcConfig{repo: repo} }

func normEnv(env string) (string, error) {
	if env == "" {
		env = "dev"
	}
	if !allowedEnv[env] {
		return "", gerror.Newf("非法环境: %s(仅 dev/test/prod)", env)
	}
	return env, nil
}

// target 解析发布目标(校验 service/env, 计算 ns/dataId/group)。
func (s *sSvcConfig) target(ctx context.Context, svc, env string) (binary, ns, dataId, group, envOut string, err error) {
	binary = catalog[svc]
	if binary == "" {
		err = gerror.Newf("未知服务: %s", svc)
		return
	}
	envOut, err = normEnv(env)
	if err != nil {
		return
	}
	ns, err = nacospub.NamespaceForEnv(ctx, envOut)
	if err != nil {
		return
	}
	dataId = binary + ".yaml"
	group = nacospub.ServiceGroup(ctx)
	return
}

func (s *sSvcConfig) Current(ctx context.Context, svc, env string) (*service.CurrentConfigDTO, error) {
	binary, ns, dataId, group, envOut, err := s.target(ctx, svc, env)
	if err != nil {
		return nil, err
	}
	content, err := nacospub.Get(ctx, ns, dataId, group)
	if err != nil {
		return nil, err
	}
	return &service.CurrentConfigDTO{
		Service: svc, Binary: binary, Env: envOut, NamespaceId: ns,
		DataId: dataId, Group: group, Content: content, Exists: content != "",
	}, nil
}

func (s *sSvcConfig) Publish(ctx context.Context, in service.PublishInput) (*service.PublishDTO, error) {
	if in.Content == "" {
		return nil, gerror.New("配置内容不能为空")
	}
	return s.doPublish(ctx, in.Service, in.Env, in.Content, in.Operator, in.Remark)
}

func (s *sSvcConfig) History(ctx context.Context, svc, env string, page, size int) (*service.HistoryDTO, error) {
	binary, _, _, _, envOut, err := s.target(ctx, svc, env)
	if err != nil {
		return nil, err
	}
	_ = binary
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := s.repo.History(ctx, svc, envOut, page, size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.HistoryItem, 0, len(list))
	for _, r := range list {
		out = append(out, &service.HistoryItem{
			Id: r.Id, Version: r.Version, Env: r.Env, Content: r.Content,
			Operator: r.Operator, Remark: r.Remark, CreatedAt: fmtTime(r.CreatedAt),
		})
	}
	return &service.HistoryDTO{List: out, Total: total, Page: page, Size: size}, nil
}

func (s *sSvcConfig) Rollback(ctx context.Context, svc, env string, version int, operator int64) (*service.PublishDTO, error) {
	if version <= 0 {
		return nil, gerror.New("版本号无效")
	}
	_, _, _, _, envOut, err := s.target(ctx, svc, env)
	if err != nil {
		return nil, err
	}
	rec, err := s.repo.FindVersion(ctx, svc, envOut, version)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, gerror.Newf("版本 v%d 不存在", version)
	}
	return s.doPublish(ctx, svc, envOut, rec.Content, operator, fmt.Sprintf("回滚自 v%d", version))
}

// doPublish 校验 YAML -> 发布 Nacos -> 落发布记录。
func (s *sSvcConfig) doPublish(ctx context.Context, svc, env, content string, operator int64, remark string) (*service.PublishDTO, error) {
	_, ns, dataId, group, envOut, err := s.target(ctx, svc, env)
	if err != nil {
		return nil, err
	}
	if _, err := gjson.LoadContentType(gjson.ContentTypeYaml, []byte(content)); err != nil {
		return nil, gerror.Wrap(err, "配置内容不是合法的 YAML")
	}
	if err := nacospub.Publish(ctx, ns, dataId, group, content); err != nil {
		return nil, err
	}
	version, err := s.repo.NextVersion(ctx, svc, envOut)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.CreateRecord(ctx, &entity.ServiceConfigPublish{
		Service: svc, Env: envOut, Version: version,
		Content: content, Operator: operator, Remark: remark,
	}); err != nil {
		return nil, err
	}
	return &service.PublishDTO{Version: version, DataId: dataId}, nil
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

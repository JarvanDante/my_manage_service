// Package nacospub Nacos 配置发布客户端(控制面 = 唯一写方)。
// 连接信息读总后台自身配置 nacospub 段; 按 namespaceId 缓存 ConfigClient。
//
// config.yaml 示例:
//
//	nacospub:
//	  addr:  "127.0.0.1:8848"
//	  user:  "nacos"
//	  pass:  "nacos"
//	  group: "DEFAULT_GROUP"
//	  namespaces:            # env -> namespaceId
//	    dev:  "xxxx-xxxx"
//	    test: ""
//	    prod: ""
package nacospub

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	mu          sync.Mutex
	clientsByNs = map[string]config_client.IConfigClient{}
)

// NamespaceForEnv 读取 env 对应的 namespaceId(未配置返回错误)。
func NamespaceForEnv(ctx context.Context, env string) (string, error) {
	v, err := g.Cfg().Get(ctx, "nacospub.namespaces."+env)
	if err != nil {
		return "", err
	}
	ns := v.String()
	if ns == "" {
		return "", gerror.Newf("nacospub.namespaces.%s 未配置(该环境的 Nacos namespaceId)", env)
	}
	return ns, nil
}

// Group 返回发布分组(默认 DEFAULT_GROUP)。
func Group(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "nacospub.group", "DEFAULT_GROUP").String()
}

// ServiceGroup 平台服务配置发布分组(与站点配置区分, 默认 SERVICE)。
func ServiceGroup(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "nacospub.service_group", "SERVICE").String()
}

func getClient(ctx context.Context, namespaceId string) (config_client.IConfigClient, error) {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := clientsByNs[namespaceId]; ok {
		return c, nil
	}
	addr := g.Cfg().MustGet(ctx, "nacospub.addr").String()
	if addr == "" {
		return nil, gerror.New("nacospub.addr 未配置(Nacos 地址)")
	}
	host, port := parseAddr(addr)
	c, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         namespaceId,
			Username:            g.Cfg().MustGet(ctx, "nacospub.user").String(),
			Password:            g.Cfg().MustGet(ctx, "nacospub.pass").String(),
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./temp/nacospub/log",
			CacheDir:            "./temp/nacospub/cache",
			LogLevel:            "warn",
		},
		ServerConfigs: []constant.ServerConfig{{IpAddr: host, Port: port}},
	})
	if err != nil {
		return nil, gerror.Wrap(err, "创建 Nacos 发布客户端失败")
	}
	clientsByNs[namespaceId] = c
	return c, nil
}

// Publish 发布配置。
func Publish(ctx context.Context, namespaceId, dataId, group, content string) error {
	c, err := getClient(ctx, namespaceId)
	if err != nil {
		return err
	}
	ok, err := c.PublishConfig(vo.ConfigParam{
		DataId: dataId, Group: group, Content: content, Type: "yaml",
	})
	if err != nil {
		return gerror.Wrap(err, "Nacos 发布失败")
	}
	if !ok {
		return gerror.New("Nacos 发布返回失败")
	}
	return nil
}

// Get 读取当前配置(不存在返回空串)。
func Get(ctx context.Context, namespaceId, dataId, group string) (string, error) {
	c, err := getClient(ctx, namespaceId)
	if err != nil {
		return "", err
	}
	content, err := c.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	if err != nil {
		if strings.Contains(err.Error(), "config data not exist") {
			return "", nil
		}
		return "", gerror.Wrap(err, "Nacos 读取失败")
	}
	return content, nil
}

func parseAddr(a string) (string, uint64) {
	a = strings.TrimSpace(a)
	host, portStr := a, "8848"
	if i := strings.LastIndex(a, ":"); i >= 0 {
		host, portStr = a[:i], a[i+1:]
	}
	p, _ := strconv.ParseUint(portStr, 10, 64)
	if p == 0 {
		p = 8848
	}
	return host, p
}

// Package boot 启动引导。按需从 Nacos 加载总后台自身配置。
//
// 与站点服务(my_service)的差异 —— 控制面是配置写方, 不能因 Nacos 不可用而无法启动:
//   - 未设 NACOS_ADDR: 用本地 manifest/config/config.yaml(默认)
//   - 设了 NACOS_ADDR 但连接/读取失败: 打告警日志后回退本地配置, 服务照常启动
//
// 环境变量:
//
//	NACOS_ADDR       Nacos 地址, 如 127.0.0.1:8848 (不设则用本地文件)
//	NACOS_NAMESPACE  命名空间 ID (dev/test/prod 命名空间的 ID; public 为空)
//	NACOS_GROUP      分组, 默认 DEFAULT_GROUP
//	NACOS_DATAID     dataId, 默认 my_manage(对应 Nacos 里的 my_manage 配置)
//	NACOS_USER/NACOS_PASS  开启鉴权时的账号密码
package boot

import (
	"strconv"
	"strings"

	nacos "github.com/gogf/gf/contrib/config/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func init() {
	addr := genv.Get("NACOS_ADDR").String()
	if addr == "" {
		return // 本地开发: manifest/config/config.yaml
	}
	ctx := gctx.GetInitCtx()
	host, port := parseAddr(addr)

	namespace := genv.Get("NACOS_NAMESPACE").String()
	group := genv.Get("NACOS_GROUP", "DEFAULT_GROUP").String()
	dataId := genv.Get("NACOS_DATAID", "my_manage").String()
	username := genv.Get("NACOS_USER").String()
	password := genv.Get("NACOS_PASS").String()

	adapter, err := nacos.New(ctx, nacos.Config{
		ServerConfigs: []constant.ServerConfig{{IpAddr: host, Port: port}},
		ClientConfig: constant.ClientConfig{
			NamespaceId:         namespace,
			Username:            username,
			Password:            password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./temp/nacos/log",
			CacheDir:            "./temp/nacos/cache",
			LogLevel:            "warn",
		},
		ConfigParam: vo.ConfigParam{DataId: dataId, Group: group},
		Watch:       true, // 配置变更热更新
	})
	if err != nil {
		// 控制面兜底: 不 Fatal, 回退本地配置, 保证 Nacos 故障时总后台仍可启动修复
		g.Log().Warningf(ctx, "nacos 配置加载失败, 回退本地 config.yaml: %+v", err)
		return
	}
	// 读一次校验配置确实存在, 不存在同样回退
	if _, err = adapter.Get(ctx, "server"); err != nil {
		g.Log().Warningf(ctx, "nacos 配置读取校验失败(dataId=%s), 回退本地 config.yaml: %+v", dataId, err)
		return
	}
	g.Cfg().SetAdapter(adapter)
	g.Log().Infof(ctx, "配置来源=Nacos addr=%s ns=%s group=%s dataId=%s", addr, namespace, group, dataId)
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

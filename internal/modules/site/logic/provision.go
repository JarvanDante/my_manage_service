// Package logic — M4 开站清单(provision-check)。
// 半自动开站的固化校验: 五项全过才允许站点上线(SetStatus→1 的门禁)。
package logic

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_manage_service/internal/dao/sitedb"
	"github.com/JarvanDante/my_manage_service/internal/modules/site/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
	"github.com/JarvanDante/my_manage_service/internal/shared/nacospub"
)

// ProvisionCheck 逐项返回开站清单结果, 不因单项失败而中断。
func (s *sSite) ProvisionCheck(ctx context.Context, id int64) (*service.ProvisionDTO, error) {
	site, err := s.mustSite(ctx, id)
	if err != nil {
		return nil, err
	}
	items := make([]*service.ProvisionItemDTO, 0, 5)
	add := func(key, name string, ok bool, detail string) {
		items = append(items, &service.ProvisionItemDTO{Key: key, Name: name, Ok: ok, Detail: detail})
	}

	// 1) 商户有效
	if ok, err := s.repo.MerchantExists(ctx, site.MerchantId); err != nil {
		add("merchant", "商户有效", false, err.Error())
	} else if !ok {
		add("merchant", "商户有效", false, "商户不存在或已停用")
	} else {
		add("merchant", "商户有效", true, "")
	}

	// 2) 域名: 至少一个且有主域名
	domains, err := s.repo.Domains(ctx, id)
	if err != nil {
		add("domain", "域名配置", false, err.Error())
	} else if len(domains) == 0 {
		add("domain", "域名配置", false, "尚未绑定任何域名")
	} else {
		hasMain := false
		for _, d := range domains {
			if d.IsMain == 1 {
				hasMain = true
				break
			}
		}
		if !hasMain {
			add("domain", "域名配置", false, "已绑定域名但未设置主域名")
		} else {
			add("domain", "域名配置", true, fmt.Sprintf("共 %d 个域名", len(domains)))
		}
	}

	// 3) DB 可连通
	if pass, err := kit.Decrypt(site.DbPass); err != nil {
		add("db", "数据库连通", false, "解密登记密码失败: "+err.Error())
	} else if ms, err := sitedb.Ping(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName); err != nil {
		add("db", "数据库连通", false, err.Error())
	} else {
		add("db", "数据库连通", true, fmt.Sprintf("%dms", ms))

		// 4) 迁移版本(依赖 DB 可连, 放在成功分支里)
		expected := g.Cfg().MustGet(ctx, "provision.gooseVersion", 0).Int64()
		ver, err := sitedb.GooseVersion(ctx, site.DbHost, site.DbPort, site.DbUser, pass, site.DbName)
		switch {
		case err != nil:
			add("migration", "迁移版本", false, "查询 goose_db_version 失败(未跑过迁移?): "+err.Error())
		case ver <= 0:
			add("migration", "迁移版本", false, "站点库尚未执行任何迁移")
		case expected > 0 && ver < expected:
			add("migration", "迁移版本", false,
				fmt.Sprintf("当前 v%d < 期望 v%d, 请先在站点库执行 goose up", ver, expected))
		default:
			add("migration", "迁移版本", true, fmt.Sprintf("当前 v%d", ver))
		}
	}
	if len(items) == 3 { // DB 失败时补迁移项(标记失败)
		add("migration", "迁移版本", false, "数据库不可连, 无法校验")
	}

	// 5) Nacos 配置已发布
	if ns, err := nacospub.NamespaceForEnv(ctx, site.Env); err != nil {
		add("config", "配置已发布", false, err.Error())
	} else if content, err := nacospub.Get(ctx, ns, site.SiteCode+".yaml", nacospub.Group(ctx)); err != nil {
		add("config", "配置已发布", false, err.Error())
	} else if content == "" {
		add("config", "配置已发布", false, "Nacos 上尚无该站点配置, 请先发布")
	} else {
		add("config", "配置已发布", true, "dataId="+site.SiteCode+".yaml")
	}

	all := true
	for _, it := range items {
		if !it.Ok {
			all = false
			break
		}
	}
	return &service.ProvisionDTO{AllPassed: all, Items: items}, nil
}

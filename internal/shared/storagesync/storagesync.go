// Package storagesync 将站点 APPKEY/SECRET 同步到统一存储 paas_client。
package storagesync

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// UpsertClient 同步调用方凭证到 my_storage（best-effort，失败只打日志）。
func UpsertClient(ctx context.Context, appKey, appSecret, siteCode string) {
	enabled := g.Cfg().MustGet(ctx, "paas.storage.enabled", false).Bool()
	if !enabled {
		return
	}
	base := g.Cfg().MustGet(ctx, "paas.storage.base_url", "http://127.0.0.1:8015").String()
	token := g.Cfg().MustGet(ctx, "paas.storage.admin_token", "").String()
	if base == "" || token == "" || appKey == "" || appSecret == "" {
		g.Log().Warning(ctx, "storagesync: skip, missing base_url/admin_token/credentials")
		return
	}
	url := base + "/admin/clients"
	res, err := g.Client().
		Header(map[string]string{
			"X-Admin-Token": token,
			"Content-Type":  "application/json",
		}).
		Put(ctx, url, g.Map{
			"app_key":    appKey,
			"app_secret": appSecret,
			"site_code":  siteCode,
			"status":     1,
			"remark":     "sync from my_manage_service",
		})
	if err != nil {
		g.Log().Warningf(ctx, "storagesync upsert failed app_key=%s: %v", appKey, err)
		return
	}
	defer res.Close()
	if res.StatusCode >= 300 {
		g.Log().Warningf(ctx, "storagesync upsert http=%d app_key=%s body=%s",
			res.StatusCode, appKey, res.ReadAllString())
		return
	}
	g.Log().Infof(ctx, "storagesync upsert ok app_key=%s site=%s", appKey, siteCode)
}

// Package sitedb 站点库动态连接探测(控制面不持有数据面模型, 仅临时连接校验/查询)。
package sitedb

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
)

// Ping 用登记信息临时建连并 ping, 返回耗时(ms)。
func Ping(ctx context.Context, host string, port int, user, pass, dbname string) (int64, error) {
	node := gdb.ConfigNode{
		Type:  "pgsql",
		Host:  host,
		Port:  fmt.Sprintf("%d", port),
		User:  user,
		Pass:  pass,
		Name:  dbname,
		Extra: "sslmode=disable",
	}
	db, err := gdb.New(node)
	if err != nil {
		return 0, err
	}
	defer db.Close(ctx)
	start := time.Now()
	if err = db.PingMaster(); err != nil {
		return 0, err
	}
	// 再跑一个真实查询确认权限
	if _, err = db.GetValue(ctx, "SELECT 1"); err != nil {
		return 0, err
	}
	return time.Since(start).Milliseconds(), nil
}

// GooseVersion 查询站点库当前 goose 迁移版本(M4 开站清单用; 表不存在返回 0)。
func GooseVersion(ctx context.Context, host string, port int, user, pass, dbname string) (int64, error) {
	node := gdb.ConfigNode{
		Type: "pgsql", Host: host, Port: fmt.Sprintf("%d", port),
		User: user, Pass: pass, Name: dbname, Extra: "sslmode=disable",
	}
	db, err := gdb.New(node)
	if err != nil {
		return 0, err
	}
	defer db.Close(ctx)
	v, err := db.GetValue(ctx,
		`SELECT coalesce(max(version_id), 0) FROM goose_db_version WHERE is_applied = true`)
	if err != nil {
		return 0, err
	}
	return v.Int64(), nil
}

// Stats 站点库核心指标(M5 跨站汇总): 用户数/今日新增/今日活跃/累计充值。
type Stats struct {
	UserCount   int
	TodayNew    int
	TodayActive int
	PaidAmount  float64
	PaidOrders  int
}

// QueryStats 连站点库聚合核心指标。today 为 YYYYMMDD。
func QueryStats(ctx context.Context, host string, port int, user, pass, dbname string, today int) (*Stats, error) {
	node := gdb.ConfigNode{
		Type: "pgsql", Host: host, Port: fmt.Sprintf("%d", port),
		User: user, Pass: pass, Name: dbname, Extra: "sslmode=disable",
	}
	db, err := gdb.New(node)
	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)

	st := &Stats{}
	one, err := db.GetOne(ctx, `
		SELECT count(*)                                              AS total,
		       count(*) FILTER (WHERE register_date = ?)             AS today_new,
		       count(*) FILTER (WHERE last_login_at >= current_date) AS today_active
		  FROM users`, today)
	if err != nil {
		return nil, err
	}
	st.UserCount = one["total"].Int()
	st.TodayNew = one["today_new"].Int()
	st.TodayActive = one["today_active"].Int()

	one, err = db.GetOne(ctx, `
		SELECT coalesce(sum(amount), 0) AS amt, count(*) AS cnt
		  FROM recharge_order WHERE status = 1`)
	if err != nil {
		return nil, err
	}
	st.PaidAmount = one["amt"].Float64()
	st.PaidOrders = one["cnt"].Int()
	return st, nil
}

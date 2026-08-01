package cmd

// 注册数据库/缓存适配器。GoFrame v2 将驱动/适配器拆为独立 contrib 模块, 必须显式空导入。
// 注意: 控制面不引 Nacos 配置引导 —— 它是配置写方, 自身配置走本地文件+env。
import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL 驱动
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"   // Redis 适配器
)

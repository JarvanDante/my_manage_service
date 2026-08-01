package cmd

// 注册数据库/缓存适配器与配置引导。
// 控制面的 Nacos 引导带本地兜底: Nacos 不可用时回退 config.yaml, 保证总后台永远可启动。
import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL 驱动
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"   // Redis 适配器

	_ "github.com/JarvanDante/my_manage_service/internal/boot" // Nacos 配置引导(NACOS_ADDR 触发, 失败回退本地)
)

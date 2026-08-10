# Nacos 与 PaaS：职责边界

> 总后台（控制面）同时涉及 **配置下发** 与 **平台能力入口**。二者互补，不要混用。

## 一句话

| 通道 | 角色 | 方向 |
|------|------|------|
| **Nacos** | 配置通道 | 总后台 **下发** → 子站拉取/热更新 |
| **PaaS** | 能力通道 | 子站持凭证 **主动调用** → 平台服务 |

- **Nacos**：下传基本的站点运行配置。
- **PaaS**：提供三方/平台能力的集成调用（支付、广告、媒资、存储、播放等）。

## Nacos 负责什么

子站启动与运行所需的**本机配置**，例如：

- 数据库连接、Redis、站点域名相关项
- 响应字段差异等站级差异（见 my_service `docs/site-features-nacos.md`）
- **调用 PaaS 所需的最小引导配置**：`paas.app_key` / `paas.app_secret`、各能力服务的 base URL 等

发布路径：总后台站点配置发布 → `nacospub` 写入 Nacos → 子站 `boot/nacos.go` 按 `SITE_CODE` 消费（可热更新）。

**不应再塞进 Nacos 的（PaaS 落地后）：**

- 支付渠道、商户号、回调路由等「支付产品配置」
- 广告位、素材、投放策略等「广告产品配置」
- 其它本应落在平台侧、按站点身份管理的业务配置

## PaaS 负责什么

平台级可复用能力（总后台「平台服务」菜单 / TPP 门户入口），例如：

| 能力 | 路径（总后台） | 子站调用方式（目标态） |
|------|----------------|------------------------|
| 媒资中心 | `/paas/media` | HTTP + `app_key` / `app_secret` |
| 统一存储 | `/paas/storage` | 同上 |
| 统一播放 | `/paas/play` | 同上 |
| 支付中台 | `/paas/pay` | 同上 |
| 广告中台 | `/paas/ad` | 同上 |

- 开站即签发站点凭证（`site.app_key` / `app_secret`）；配置发布时注入 Nacos 的 `paas` 段，供子站读取。
- 业务配置（渠道、广告位等）落在 **PaaS 控制面 / 各能力服务**，按 `app_key` 识别调用方站点，而不是再 YAML 下发到每站。

PaaS **不是**「再下一份配置文件」：子站在业务链路里 **同步/异步调 API**（下单、查素材、取播放地址等）。

## 对照（以前 vs 现在）

| 场景 | 以前（偏 Nacos） | 目标态（Nacos + PaaS） |
|------|------------------|------------------------|
| 支付/广告怎么接到子站 | 渠道密钥、开关等写进各站 YAML | 子站用凭证调支付中台 / 广告中台 |
| 站库、Redis、字段差异 | Nacos | **仍走 Nacos** |
| 凭证如何到子站 | 混在业务配置里 | 开站签发；发布时注入 `paas.*` |

## 过渡建议

1. Nacos 继续承载站点基础配置 + PaaS 引导项。
2. 支付/广告等能力按服务逐步迁出：新流量走 PaaS API，旧站可短暂双轨。
3. 总后台 TPP「能力概览」只做目录与跳转；真实服务独立部署，与 HLS 图床等 CDN 配置无关。

## 相关代码 / 文档

- 总后台架构：`docs/architecture.md`
- Nacos 发布：`internal/shared/nacospub`、`internal/modules/siteconfig`（含 `injectPaas`）
- 站点凭证：`manifest/sql/migrations/00011_site_app_credentials.sql`
- 子站字段差异：my_service `docs/site-features-nacos.md`
- 子站 Nacos 模板：my_service `docs/nacos-config-my.yaml`
- 媒资中心服务：`my_media`（中央池 API；转码 Worker 见 `my_transcode`）
- 媒资 M1 已验收链路：`my_media/docs/m1-pipeline.md`（上传→Kafka→转码→HLS→ready→可播）

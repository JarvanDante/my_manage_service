# my_manage_service · 总后台(控制面)架构与目录结构

> 独立 git 仓库,GoFrame v2,单 HTTP 服务(:8003),独立库 `my_manage`。
> 管理对象:商户 → 站点 → 域名 → 配置下发(Nacos)。开发计划见 my_service 仓库 `docs/manage_api_plan.md`(M0~M6)。

## 与 my_service 的关系

| | my_service(数据面) | my_manage_service(控制面) |
|---|---|---|
| 部署 | 一站一部署(front/backend/cron) | 全局一套 |
| 数据库 | 每站一库(my_service, 站B库…) | my_manage 一库 |
| 配置来源 | 启动从 Nacos 拉(SITE_CODE) | 本地 config.yaml + env |
| Nacos 角色 | 只读(消费配置) | 只写(发布配置, nacos-sdk-go) |
| 管理员 | admin_user(站点库) | manage_admin(互不相通) |

## 四个关键决策

1. **单门面**:只有一种使用者,controller 不分 front/backend 子目录;但分层约定(domain 零框架依赖 / module-first / dao 在 cmd 注入)与 my_service 完全一致。
2. **自身配置不走 Nacos**:控制面是配置写方,不能依赖自己管理的组件启动(Nacos 挂了必须还能登上总后台修复)。本地 config.yaml + env;Nacos 仅作为发布客户端(shared/nacospub)。
3. **站点库只动态连接**:db-check / 开站校验 / 跨站汇总时用登记信息(AES 解密)临时连接查询(dao/sitedb/probe.go);不 gf gen dao、不持有数据面模型。
4. **shared 复制不抽库**:kit(token 前缀 manage:token: + AES)/ middleware / rbac 从 my_service 复制精简;两边稳定后再考虑抽公共 pkg 仓库。my_service 里旧的 /manage 门面与 manageapi 二进制在 M0 动工时废弃。

## 目录结构

```
my_manage_service/                          # module github.com/JarvanDante/my_manage_service
├── api/                                    # ① 接口契约(单门面, 按业务模块分)
│   ├── auth/v1/auth.go                     #    M0 登录/退出/信息
│   ├── merchant/v1/merchant.go             #    M1 商户 CRUD
│   ├── site/v1/{site,domain}.go            #    M2 站点/域名/db-check
│   ├── siteconfig/v1/config.go             #    M3 配置发布/历史/回滚
│   └── monitor/v1/monitor.go               #    M5 健康/跨站汇总
│
├── internal/
│   ├── cmd/
│   │   ├── cmd.go                          # 单 server :8003, /manage 分组:
│   │   │                                   #   公开(login) / 仅登录(auth/*) / 权限组(其余)
│   │   └── driver.go                       # pgsql + redis 驱动 blank import
│   │
│   ├── modules/                            # ② module-first(同 my_service 约定)
│   │   ├── auth/
│   │   │   ├── domain/repository.go        #    纯接口, 零框架依赖
│   │   │   ├── service/  logic/  controller/  router.go
│   │   ├── merchant/
│   │   ├── site/
│   │   ├── siteconfig/
│   │   └── monitor/
│   │
│   ├── dao/                                # ③ 仓储实现; g.DB() 默认库 = my_manage
│   │   ├── auth.go  merchant.go  site.go  siteconfig.go
│   │   └── sitedb/probe.go                 #    站点库动态连接探测
│   │
│   ├── model/entity/                       # manage_admin / merchant / site / site_domain /
│   │                                       # site_config_publish ...
│   └── shared/
│       ├── consts/  errcode/
│       ├── kit/                            # token(manage:token:) + AES(站点DB密码)
│       ├── middleware/                     # ManageAuth / ManagePerm / OpLog / CORS / NotFound
│       ├── rbac/                           # Casbin(表在 my_manage.casbin_rule)
│       └── nacospub/                       # nacos-sdk-go 发布客户端封装(唯一写方)
│
├── manifest/
│   ├── config/config.yaml                  # 本地配置(不依赖 Nacos 启动)
│   └── sql/migrations/                     # goose → my_manage
│       ├── 00001_create_manage_admin.sql   # M0: admin/role/casbin_rule/log
│       ├── 00002_create_merchant.sql       # M1
│       ├── 00003_create_site.sql           # M2: site + site_domain
│       └── 00004_create_config_publish.sql # M3
│
├── test/http/                              # .http 端到端测试
├── hack/config.yaml                        # gf gen dao(连 my_manage)
├── Makefile
├── go.mod
└── main.go
```

## 依赖清单

- github.com/gogf/gf/v2 + contrib/drivers/pgsql + contrib/nosql/redis
- github.com/casbin/casbin/v2(权限, M6 前仅超管放行)
- github.com/nacos-group/nacos-sdk-go/v2(配置发布)
- goose(迁移, 命令行工具)

## 路由分组(与 my_service 后台同模式)

```
/manage
├── (公开)      POST /manage/auth/login
├── (仅登录)    POST /manage/auth/logout · GET /manage/auth/info
└── (权限组)    ManageAuth + ManagePerm + ManageOpLog
    ├── /manage/merchants...        (M1)
    ├── /manage/sites...            (M2: 站点/域名/db-check)
    ├── /manage/sites/{id}/config...(M3: 发布/历史/回滚)
    └── /manage/overview, /logs...  (M5)
```

## 开工前置(手动一次)

- [ ] `CREATE DATABASE my_manage;`
- [ ] 环境变量 `MANAGE_SECRET`(AES 密钥, 站点 DB 密码加密)
- [ ] 确认 dev/test/prod 三个 Nacos namespaceId(M3 发布用)

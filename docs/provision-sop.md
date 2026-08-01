# 开站 SOP(半自动)

> 新商户开站的标准流程。人工步骤(①②)+ 总后台校验与下发(③~⑦)。
> 上线是**状态机门禁**: `POST /manage/sites/{id}/status {status:1}` 会自动跑开站清单, 五项全过才放行。

## 流程

| 步 | 操作 | 谁做 | 说明 |
|---|---|---|---|
| ① | 建站点库并跑迁移 | 运维/你 | `CREATE DATABASE <站点库>;` 然后在 my_service 仓库 `goose -dir manifest/sql/migrations postgres "...dbname=<站点库>..." up` |
| ② | (可选)种子数据 | 运维/你 | 站点库的 admin 种子随迁移 00007 自带(admin/admin123, 上线前改密) |
| ③ | 总后台建商户 | 运营 | `POST /manage/merchants` |
| ④ | 总后台建站点 | 运营 | `POST /manage/sites`(site_code 定好不可改; 登记①的库连接信息, 密码 AES 加密存) |
| ⑤ | 绑定域名 | 运营 | `POST /manage/sites/{id}/domains`(首个自动为主域名) |
| ⑥ | 发布配置 | 运营 | `POST /manage/sites/{id}/config/publish`(database 段自动注入, 不用写) |
| ⑦ | 开站清单校验 | 运营 | `POST /manage/sites/{id}/provision-check` 五项全绿 |
| ⑧ | 上线 | 运营 | `POST /manage/sites/{id}/status {status:1}`(清单不过会被拒) |
| ⑨ | 部署站点服务 | 运维 | my_service 二进制 + env: `NACOS_ADDR / NACOS_NAMESPACE=<env的namespaceId> / SITE_CODE=<site_code> / NACOS_USER/PASS` |
| ⑩ | 域名解析/网关 | 运维 | 把⑤的域名指到该部署(Nginx/云解析) |

## 开站清单五项

1. **merchant** 商户存在且启用
2. **domain** 至少绑定一个域名且已设主域名
3. **db** 用登记信息真连站点库成功
4. **migration** 站点库 goose 版本 ≥ 配置的 `provision.gooseVersion`(0=只要求跑过迁移)
5. **config** Nacos 上已发布 `{site_code}.yaml`

## 注意

- `provision.gooseVersion` 在总后台 config.yaml 维护: my_service 每加一个迁移, 这里同步 +1, 防止漏跑迁移的库被上线。
- MANAGE_SECRET 不能变更(变更后旧站点登记密码解不开)。
- 站点停用: `status:2`; 商户停用前须先停其名下全部站点。

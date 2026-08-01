# jh 前端(vben-admin)扫描总结 与 my 项目前端方案

> 扫描对象:
> - `jh_backend` —— 子系统(站点)后台前端
> - `jh_manage_backend` —— 总后台前端
> 目标: 参考两者, 规划 my 项目的「站点后台前端」与「总后台前端」。

---

## 一、两个项目的共同技术底座

| 项 | 结论 |
|---|---|
| 框架 | **Vben Admin v5.5.9** monorepo(pnpm workspace + turbo),应用变体 `apps/web-ele`(Element Plus) |
| 技术栈 | Vue3 + Vite + TypeScript + Pinia + vue-router |
| 访问控制 | `accessMode: "backend"`(后端菜单驱动): 登录后调 `GET /app/menus` 拉**角色菜单树** → `convertMenuToRoute` 转成路由 → `import.meta.glob("../views/**/*.vue")` 按 path 映射页面组件 → `generateAccessible("mixed")` 挂载 |
| 菜单数据 | 后端有 menu 表(id/type/name/path/icon/sort/children), 角色绑定菜单 → **RBAC 是「菜单级」**(能看见哪些页面), 与我们 my 后端的「接口级 Casbin」互补 |
| 请求层 | `request.ts`: baseURL=`VITE_GLOB_API_URL`(如 `/api/backend`), dev 走 Vite proxy(`VITE_PROXY_TARGET`), Authorization 头带 token, 401 统一登出 |

## 二、jh_backend(站点后台)的关键机制

1. **动态品牌(一套前端所有站点复用)**:
   - 登录响应 `userInfo` 携带 `site_code / site_name`
   - 存 `localStorage.VBEN_SITE_NAME` → `preferences.ts` 启动时读取覆盖 `app.name`
   - 退出登录时清除; logo 走 `/logo.png` 可按部署替换
   - 前端不判断站点, 站点身份由后端(token/域名)决定
2. 业务模块: dashboard / userManage / balanceManage(充提) / gameManage / agent / operations / sysSetting / adminManage

## 三、jh_manage_backend(总后台)的关键机制

1. 模块: domainManage(site_code+domain+type+status)、paymentManage+**sitePayment**、gameManage+siteGame、activityManage+siteActivity、customerManage(客户/站点)、withdrawManage、betManage、adminManage
2. **「全局定义 + 站点绑定」模式**: 支付/游戏/活动先做全局资源定义, 再用 site_x 关联表 bind/unbind 到站点(带 status/sort) —— 将来 my 做支付渠道管理可直接借鉴
3. **发布模型与 my 的 M3 高度一致**(相互印证):
   - `is_sync`(是否已发布)、`publish_type`(1发布/2回滚)、`rollback_from`、`operator`、`version`、`is_current`、版本历史列表
   - 区别仅在通道: jh 发 Consul→子站监听落表; my 发 Nacos→站点热更新配置

## 四、my 项目前端方案(建议)

### 4.1 仓库规划
| 新仓库 | 参考 | 对接后端 | baseURL |
|---|---|---|---|
| `my_backend`(站点后台前端) | jh_backend | my_service `:8000` | `/api/backend` → proxy |
| `my_manage_backend`(总后台前端) | jh_manage_backend | my_manage_service `:8003` | `/api/manage` → proxy |

均用 Vben v5 web-ele 起(与 jh 同款, 迁移页面成本最低)。

### 4.2 三个必须适配的差异点
1. **响应格式**: my 后端统一 `{code, message, data}`(code=0 成功)。request.ts 响应拦截器需按此适配(jh 版本已很接近, 改判定字段即可); 401 判定对齐 GoFrame `CodeNotAuthorized`。
2. **菜单来源**: my 后端目前**没有 menu 表**。两条路:
   - A. 前端静态路由(accessMode=frontend) + 后端 Casbin 管接口 —— 最快跑通, 总后台(人少)推荐;
   - B. 增量加 menu 表 + 角色-菜单绑定接口(站点后台 B0.3 / 总后台 M6.5), 复刻 jh 动态菜单体验 —— 站点后台(运营角色多)推荐。
   可以先全走 A, 需要精细控页面时再补 B。
3. **动态品牌**: 复用 jh 方案, 只需 my_service 的 `/backend/auth/info`(或 login 响应)加 `site_code/site_name` 字段(读 SITE_CODE env + site_config/Nacos), 前端照抄 VBEN_SITE_NAME 那套。

### 4.3 页面与已有接口的映射(开工清单)
**my_manage_backend(总后台)** —— 后端接口已全部就绪(M0~M6):
| 页面 | 对接 |
|---|---|
| 登录/管理员/角色权限 | /manage/auth·admins·roles·perms(M0/M6) |
| 商户管理 | /manage/merchants(M1) |
| 站点管理+域名 | /manage/sites·domains·db-check(M2) |
| 配置发布(编辑器+历史+回滚) | /manage/sites/{id}/config*(M3, 页面参考 jh customer 发布/版本历史) |
| 开站清单 | /manage/sites/{id}/provision-check(M4, 清单式 UI) |
| Dashboard/健康/审计 | /manage/overview·health·logs(M5) |

**my_backend(站点后台)** —— 后端接口已全部就绪(B0~B8):
| 页面 | 对接 |
|---|---|
| 登录/管理员/角色权限 | /backend/auth·admins·roles·perms(B0.x) |
| 用户管理 | /backend/users*(B1) |
| 财务(套餐/订单/流水/回调状态) | /backend/recharge-*·vip-*·balance-logs(B2) |
| 兑换码 | /backend/codes*(B3) |
| 用户组 | /backend/user-groups(B4) |
| 任务/签到统计 | /backend/tasks·task-logs·sign-stats(B5) |
| 社交/推广 | /backend/follows·share-*(B6) |
| 消息/公告/客服 | /backend/messages·notices·push·config/customer-url(B7) |
| Dashboard | /backend/stats/*(B8) |

### 4.4 建议开工顺序
1. `my_manage_backend` 先行(模块少、接口全、马上能用): 骨架适配(请求层/登录)→ 商户/站点/域名 → 配置发布(核心页) → 开站清单 → Dashboard
2. `my_backend` 随后: 骨架复用第1步成果 → 用户/财务 → 其余模块
3. 后端小增量(按需): login/info 返回 site_name; (可选)menu 表与角色绑定接口

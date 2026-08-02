# 菜单驱动 RBAC 改造 — 进度与交接(总后台阶段)

> 目标:把权限模型从「casbin 按接口路径配置」改成 **纯 RBAC 菜单驱动**。
> 权限最小单位 = 菜单;每个菜单下挂它的接口权限;给角色配权限 = 勾选菜单树。
> 参考:jh_manage_service(admin_permission 树 + admin_role.permissions)。

## 三个已确认决策
1. 接口权限:服务端强制校验。 2. 做完整的菜单/权限管理页。 3. 总后台先做,验收后照搬子后台。

## 数据模型(my_manage 库)
- 新表 manage_permission:id,parent_id,name,route_url,component,method,icon,
  is_menu(1菜单/0接口),hide_in_menu,affix_tab,active_path,sort,status
- manage_role 加 permissions 字段(逗号分隔权限id,取代casbin)
- 迁移 00005_create_manage_permission.sql —— 已 make migrate 应用(10菜单+33接口权限种子)

## 已完成:总后台后端代码(⚠️未编译,明天先 go build ./... 抓错)
新增:entity/manage_permission.go、dao/auth_permission.go、logic/permission.go、
      controller/permission.go、api/auth/v1/permission.go
改写:shared/middleware/middleware.go(ManagePerm 用权限表替代casbin,无权限返回403不再掉登录)、
      controller/menu.go(/auth/menus 改DB驱动按角色过滤)、service/auth.go、logic/mgmt.go、
      dao/auth.go、domain/repository.go、controller/mgmt.go、api/auth/v1/mgmt.go、router.go
casbin(shared/rbac)留在磁盘但不再引用。

## 新接口(/manage)
GET /auth/menus(按角色菜单树)、GET /permissions(全量权限树)、
POST/PUT/DELETE /permissions/{id}(权限节点CRUD)、角色create/update 带 permissions 字段。

## 明天续接(按序)
1. go build ./... 修编译错误 → 重启 :8003 → 超管验收菜单DB下发。
2. 总后台前端 my_manage_backend(未做):
   - 新建 views/system/permission.vue(菜单/权限管理页,树CRUD 调 /permissions)
     ⚠️ 种子已有「菜单权限」菜单指向 /system/permission,但vue未建,现在点进去会报组件找不到,建完即好。
   - 改 views/system/role.vue:角色弹窗改 ElTree 勾选菜单树(拉 /permissions,勾选存id列表);移除旧按路径配权限弹窗。
   - 角色提交带 permissions;api/core 增加权限树CRUD。
   - 前端在docker容器跑,改完重启 dev(:5778)。
3. 浏览器验收总后台RBAC(#74):超管全量;受限角色只勾部分菜单→只见对应菜单+越权403不掉登录+无死循环。
4. 照搬子后台(my_service + my_backend)。

## 环境/协作要点
- 前端只在 docker 容器跑,dev 需你重启;设备端无Go/node/网络,编译/迁移/push 你在Mac做。
- git 锁:仓库被云端挂载+本地GoLand同时提交会争 .git/index.lock;规避:提交我从云端做你只push,或删 index.lock。
- i18n 已修复根因:两前端 vite.config 加了 application:{ i18n:false }(关掉 unplugin-vue-i18n 的 runtime-only 别名)。别改回去。
- 两前端 accessMode 已=backend;menu.ts 已指向 /auth/menus。

## git 状态(待整理提交)
- my_service:动态菜单接口已提交 3a376de(待push);子后台硬编码菜单明天会被RBAC版替换。
- my_manage_service:本次RBAC后端改动未提交(等编译通过再提交)。
- my_manage_backend/my_backend:accessMode+i18n:false+menu.ts 改动未提交。
- my_backend 是重clone后从 e19ba01 导出的干净仓库。

## 任务
#70✅ #71✅ #72✅ 表+后端完成;#73⬜ 前端(角色勾菜单树+菜单管理页);#74⬜ 验收;子后台#75~待建。

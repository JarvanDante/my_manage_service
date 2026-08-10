# my 系统架构全貌(ARCHITECTURE)

> 站点品牌「漫隐」,site_code `my`。多租户漫画/社区平台:总后台(控制面)+ 子站(业务面)+ PaaS 能力层。
> 本文取代已完成归档的 RBAC_PROGRESS.md,供后续开发会话快速对齐。更新于 2026-08-10。

## 一、服务矩阵

| 仓库 | 端口 | 数据库 | 定位 |
|------|------|--------|------|
| my_manage_service | :8003 | my_manage | 总后台 API(控制面):商户/站点/开站/配置发布(Nacos)/审计/RBAC;站点凭证 APPKEY/APPSECRET 分配 |
| my_service | :8000 | **my** | 子站 API(前台 + 子后台 /backend):用户/财务/兑换码/成长/运营/视频/RBAC;配置经 Nacos 按 SITE_CODE 热更新 |
| my_storage | :8005(Docker 宿主 8015) | my_storage | PaaS 统一存储:MinIO 预签名上传/下载、元数据、配额;对象 key `{site_code}/{biz}/{code}/{filename}` |
| my_media | :8004 | my_media | PaaS 媒资中心:中央媒资池、上传编排、Kafka 投递转码、回写 play_url;M1 链路已验收 |
| my_transcode | :8088(mydocker 编排) | 无(不连业务库) | 转码 Worker:Kafka 消费 → ffprobe/ffmpeg → H.264 HLS + 封面图 → MinIO |
| ~~transcode-worker~~ | — | — | 已废弃(内容在 _to_delete/, 整个目录可删) |

前端(Vben Admin v5.5.9 monorepo, apps/web-ele, Element Plus):

| 仓库 | dev 端口 | 代理 | 说明 |
|------|---------|------|------|
| my_manage_backend | :5778 | /manage → host.docker.internal:8003 | 总后台 UI;含 paas/(storage/media/pay/ad/play) 与 tpp/ 对接文档中心 |
| my_backend | :5779 | /backend → host.docker.internal:8000 | 子后台 UI;含 video/(分片上传) |

## 二、关键数据流

```text
配置(Nacos 通道):   总后台发布 → nacospub → Nacos → 子站 boot/nacos.go 按 SITE_CODE 消费(热更新)
能力(PaaS 通道):    子站持 paas.app_key/secret → 直调 my_storage / my_media 开放门面
媒资:  总后台 /paas/media ─HTTP→ my_media ─Job→ Kafka(media.transcode.jobs) → my_transcode
       my_media ←Result─ Kafka(media.transcode.results);产物入 MinIO(预签名 host.docker.internal:19000)
存储:  Admin 门面 X-Admin-Token;开放门面 X-App-Key/X-App-Secret
```

Nacos 与 PaaS 职责边界详见 `docs/nacos-vs-paas.md`:Nacos 只下发运行配置 + PaaS 引导凭证;支付/广告等产品配置落 PaaS。

## 三、RBAC(两个后台均已上线并验收)

菜单驱动纯 RBAC,弃用 casbin(相关代码为死代码已清理):

- 权限树表:总后台 `manage_permission` / 子后台 `admin_permission`
  字段 `id,parent_id,name,route_url,component,method,icon,is_menu(1菜单/0接口),hide_in_menu,affix_tab,active_path,sort,status`
- 角色表加 `permissions`(勾选的权限 id 逗号串):`manage_role` / `admin_role`
- `GET /auth/menus`:按登录管理员角色从 DB 过滤出菜单树(超管全量),前端 vben accessMode=backend 动态注册路由
- 中间件:超管放行;其余按角色勾选的 is_menu=0 节点做「路径({id} 通配)+方法」强校验;**无权限返回 403(不是 61),避免被前端当掉登录**
- 前端:角色配置=ElTree 勾选菜单树(勾菜单即授菜单+其下接口);系统→菜单权限页可在线维护权限树
- 未勾选的菜单前端路由不注册(访问即 404),接口侧仍有 403 双保险

## 四、开发环境要点(重要约定/坑)

- **前端一律在 docker 的 node 容器里跑**(宿主 /Users/wangdante/D/kugou ↔ 容器 /var/www/html),不在 Mac 宿主机;Go 服务跑宿主机;容器编排在 /Users/wangdante/D/mydocker(网络 jarven)
- Postgres 127.0.0.1:5432(postgres/654321),库:my_manage / my / my_media / my_storage;迁移用 goose(各仓 manifest/sql/migrations,my_manage_service 与各 PaaS 有 make migrate)
- **GoFrame 表结构缓存坑**:进程启动时缓存表字段,新增列后必须「先迁移、后重启服务」,否则新字段被静默丢弃(INSERT/UPDATE/Fields 都会剔除)
- **vben 偏好缓存坑**:preferences 整体持久化在 localStorage(`{app}-preferences`),改 preferences.ts(如 accessMode)后需清掉该 key 才生效
- vue-i18n:两前端 vite.config 均设 `application:{ i18n:false }` 关掉 unplugin-vue-i18n 的 runtime-only 别名(否则显示原始 key),勿改回
- git:仓库同时被云端会话挂载 + 本地 GoLand 提交会争 .git 锁;遇 index.lock/HEAD.lock 报错直接删锁文件重试
- 迁移种子写法:显式 id + 末尾 setval 修序列 + IF NOT EXISTS/DELETE 幂等(避免自增撞显式 id)

## 五、历史里程碑

- 子后台 B0-B8、总后台 M0-M6 后端;两前端全模块开发并浏览器验收
- 动态菜单 method B → 升级为菜单驱动 RBAC(总后台先行验收,子后台移植后亦跑通)
- 统一存储/媒资中心/转码三件套落地,M1 上传→转码→HLS→播放链路验收
- 站点开站分配 APPKEY/APPSECRET;Nacos 与 PaaS 职责切分

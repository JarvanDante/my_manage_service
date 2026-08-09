# my_manage_service

漫隐总管理后台（控制面）API，默认端口 `:8003`。

## 本地开发

```bash
make migrate   # 确保 postgres 已有 my_manage 库
make build
./bin/manageapi
# 或
make dev
```

配置：`manifest/config/config.yaml`（本机 127.0.0.1）。

## Docker 服务（mydocker）

与 `my_transcode` 一样挂在 mydocker 的 `jarven` 网络：

```bash
# 先保证库存在（宿主机执行一次即可）
docker exec postgres psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname='my_manage'" | grep -q 1 \
  || docker exec postgres psql -U postgres -c "CREATE DATABASE my_manage;"
make migrate

make docker-up       # 构建并启动
make docker-logs     # 看日志
make docker-rebuild  # 改代码后重建
make docker-down
```

容器配置：`manifest/config/config.docker.yaml`  
- DB/Redis/Nacos → 容器名 `postgres` / `redis` / `nacos`  
- 媒资同步 → `host.docker.internal:8004`（宿主机 my_media）

探活：`curl -sS http://127.0.0.1:8003/api.json | head`

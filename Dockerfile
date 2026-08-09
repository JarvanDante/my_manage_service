# 构建
FROM golang:1.23-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/manageapi .

# 运行
FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/manageapi /app/manageapi
# 默认 Docker 配置；compose 可挂载覆盖
COPY manifest/config/config.docker.yaml /app/manifest/config/config.yaml
EXPOSE 8003
CMD ["/app/manageapi"]

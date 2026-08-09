APP := my_manage_service
BIN := bin
MYDOCKER := /Users/wangdante/D/mydocker

.PHONY: build dev tidy clean migrate docker-up docker-down docker-logs docker-rebuild

build:                          ## 构建
	go build -o $(BIN)/manageapi .

dev:                            ## 本地热编译运行
	gf run main.go

tidy:
	go mod tidy

migrate:                        ## goose 迁移(my_manage 库)
	goose -dir manifest/sql/migrations postgres "host=127.0.0.1 port=5432 user=postgres password=654321 dbname=my_manage sslmode=disable" up

clean:
	rm -rf $(BIN)

# ---- Docker：走 mydocker 主编排 ----
docker-up: ## 在 mydocker 中构建并启动总后台 API
	cd $(MYDOCKER) && docker compose up -d --build my_manage_service
	@echo "探活: curl -sS http://127.0.0.1:8003/api.json | head"

docker-down:
	cd $(MYDOCKER) && docker compose stop my_manage_service

docker-logs:
	docker logs -f my_manage_service

docker-rebuild:
	cd $(MYDOCKER) && docker compose up -d --build --force-recreate my_manage_service

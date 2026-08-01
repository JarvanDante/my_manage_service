APP := my_manage_service
BIN := bin

.PHONY: build dev tidy clean migrate

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

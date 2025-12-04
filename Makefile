# Makefile for MyS Site

.PHONY: help dev dev-admin build build-admin run run-admin stop clean db-migrate db-init db-sync db-health db-tables db-backup docker-up docker-down

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  make dev          - 启动用户服务（开发模式）"
	@echo "  make dev-admin    - 启动管理后台服务（开发模式）"
	@echo "  make dev-all      - 同时启动用户服务和管理后台"
	@echo "  make build        - 编译用户服务"
	@echo "  make build-admin  - 编译管理后台服务"
	@echo "  make build-all    - 编译所有服务"
	@echo "  make db-migrate   - 运行数据库迁移"
	@echo "  make db-init      - 初始化数据库数据"
	@echo "  make db-sync      - 同步计数器字段"
	@echo "  make db-health    - 数据库健康检查"
	@echo "  make db-tables    - 查看所有表"
	@echo "  make docker-up    - 启动Docker容器"
	@echo "  make docker-down  - 停止Docker容器"
	@echo "  make clean        - 清理编译文件"

# 开发模式 - 用户服务
dev:
	@echo "启动用户服务（端口8081）..."
	go run cmd/server/main.go

# 开发模式 - 管理后台
dev-admin:
	@echo "启动管理后台服务（端口8082）..."
	go run cmd/admin/main.go

# 同时启动两个服务（在后台）
dev-all:
	@echo "启动用户服务和管理后台..."
	@echo "用户服务: http://localhost:8081"
	@echo "管理后台: http://localhost:8082"
	@echo ""
	@echo "按 Ctrl+C 停止所有服务"
	@make -j2 dev dev-admin

# 编译 - 用户服务
build:
	@echo "编译用户服务..."
	go build -o bin/server cmd/server/main.go

# 编译 - 管理后台
build-admin:
	@echo "编译管理后台服务..."
	go build -o bin/admin cmd/admin/main.go

# 编译所有
build-all: build build-admin
	@echo "所有服务编译完成"

# 运行编译后的用户服务
run:
	@echo "运行用户服务..."
	./bin/server

# 运行编译后的管理后台
run-admin:
	@echo "运行管理后台..."
	./bin/admin

# 数据库迁移
db-migrate:
	@echo "运行数据库迁移..."
	go run cmd/migrate/main.go

# 初始化数据库数据
db-init:
	@echo "初始化数据库数据..."
	mysql -h localhost -u root -proot mysite < scripts/init.sql

# 同步计数器
db-sync:
	@echo "同步计数器字段..."
	mysql -h localhost -u root -proot mysite < scripts/sync-counters.sql

# 数据库健康检查
db-health:
	@echo "数据库健康检查..."
	mysql -h localhost -u root -proot -e "SELECT 1" mysite

# 查看所有表
db-tables:
	@echo "查看所有数据库表..."
	mysql -h localhost -u root -proot -e "SHOW TABLES" mysite

# 数据库备份
db-backup:
	@echo "备份数据库..."
	mysqldump -h localhost -u root -proot mysite > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Docker命令
docker-up:
	@echo "启动Docker容器..."
	docker-compose up -d

docker-down:
	@echo "停止Docker容器..."
	docker-compose down

# 清理
clean:
	@echo "清理编译文件..."
	rm -rf bin/
	rm -f mysite
	rm -f admin
	@echo "清理完成"

# Go Demo 项目 Makefile
# 前后端分离项目结构

# 项目信息
PROJECT_NAME := go-demo
VERSION := 1.0.0

# 目录配置
SERVER_DIR := server
WEB_DIR := web
BUILD_DIR := $(SERVER_DIR)/bin

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# 构建标志
.PHONY: all
all: clean deps fmt vet test build

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	$(GOFMT) -s -w .

# 代码检查
.PHONY: vet
vet:
	@echo "🔍 代码检查..."
	$(GOCMD) vet $$(go list ./... | grep -v backup_)

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	$(GOTEST) -v $$(go list ./... | grep -v backup_)

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "📊 生成测试覆盖率报告..."
	$(GOTEST) -coverprofile=coverage.out $$(go list ./... | grep -v backup_)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 构建应用
.PHONY: build
build:
	@echo "🔨 构建应用..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_PATH)

# 构建多平台版本
.PHONY: build-all
build-all:
	@echo "🌍 构建多平台版本..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe $(MAIN_PATH)

# 运行应用
.PHONY: run
run:
	@echo "🚀 运行应用..."
	$(GOCMD) run $(MAIN_PATH)

# 开发模式运行（带热重载）
.PHONY: dev
dev:
	@echo "🔥 开发模式运行..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/air-verse/air@latest"; \
		$(GOCMD) run $(MAIN_PATH) server --config=./configs/config.dev.yaml; \
	fi

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 安装开发工具
# 安装开发工具
.PHONY: install-tools
install-tools:
	@echo "🛠️ 安装开发工具..."
	$(GOCMD) install github.com/air-verse/air@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# 代码质量检查
.PHONY: lint
lint:
	@echo "🔍 代码质量检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "请先安装 golangci-lint: make install-tools"; \
	fi

# ==================== 前端相关命令 ====================

# 检查 npm 是否安装
.PHONY: check-npm
check-npm:
	@command -v npm > /dev/null 2>&1 || { \
		echo "❌ 错误: npm 未安装"; \
		echo "请先安装 Node.js: https://nodejs.org/"; \
		echo "或使用 brew install node (macOS)"; \
		exit 1; \
	}

# 安装前端依赖
.PHONY: web-install
web-install: check-npm
	@echo "📦 安装前端依赖..."
	@cd $(WEB_DIR) && npm install

# 前端开发模式
.PHONY: web-dev
web-dev: check-npm
	@if [ ! -d "$(WEB_DIR)/node_modules" ]; then \
		echo "⚠️  node_modules 不存在，正在安装依赖..."; \
		cd $(WEB_DIR) && npm install; \
	fi
	@echo "🌐 启动前端开发服务器..."
	@cd $(WEB_DIR) && npm run dev

# 前端构建
.PHONY: web-build
web-build: check-npm
	@if [ ! -d "$(WEB_DIR)/node_modules" ]; then \
		echo "⚠️  node_modules 不存在，正在安装依赖..."; \
		cd $(WEB_DIR) && npm install; \
	fi
	@echo "🔨 构建前端..."
	@cd $(WEB_DIR) && npm run build

# 前端预览
.PHONY: web-preview
web-preview: check-npm
	@echo "👀 预览前端构建..."
	@cd $(WEB_DIR) && npm run preview

# 前端代码检查
.PHONY: web-lint
web-lint: check-npm
	@echo "🔍 前端代码检查..."
	@cd $(WEB_DIR) && npm run lint

# 清理前端构建
.PHONY: web-clean
web-clean:
	@echo "🧹 清理前端构建..."
	@rm -rf $(WEB_DIR)/dist
	@rm -rf $(WEB_DIR)/node_modules

# ==================== 全栈开发命令 ====================

# 同时启动前后端（开发模式）
.PHONY: dev-all
dev-all: check-npm
	@echo "🚀 启动全栈开发环境..."
	@echo "📍 后端: http://localhost:8080"
	@echo "📍 前端: http://localhost:3000"
	@echo ""
	@if [ ! -d "$(WEB_DIR)/node_modules" ]; then \
		echo "⚠️  node_modules 不存在，正在安装依赖..."; \
		cd $(WEB_DIR) && npm install; \
	fi
	@$(MAKE) -j2 dev web-dev

# 安装所有依赖
.PHONY: install-all
install-all: deps web-install
	@echo "✅ 所有依赖安装完成"

# 构建所有
.PHONY: build-all-stack
build-all-stack: build web-build
	@echo "✅ 前后端构建完成"

# ==================== Docker 相关命令 ====================

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "🐳 构建 Docker 镜像..."
	docker build -f deployments/Dockerfile -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest

# Docker 运行
.PHONY: docker-run
docker-run:
	@echo "🐳 运行 Docker 容器..."
	docker run -p 8080:8080 $(PROJECT_NAME):latest

# 完整部署 (应用 + MySQL + Redis + Nginx)
.PHONY: docker-deploy
docker-deploy:
	@echo "🚀 开始完整部署..."
	@cd deployments && docker-compose pull
	@cd deployments && docker-compose build
	@cd deployments && docker-compose up -d
	@echo "⏳ 等待服务启动..."
	@sleep 10
	@echo "✅ 服务启动完成！"
	@$(MAKE) docker-info

# 简化部署 (应用 + MySQL + Redis)
.PHONY: docker-deploy-simple
docker-deploy-simple:
	@echo "🚀 开始简化部署..."
	@cd deployments && docker-compose -f docker-compose.simple.yml pull
	@cd deployments && docker-compose -f docker-compose.simple.yml build
	@cd deployments && docker-compose -f docker-compose.simple.yml up -d
	@echo "⏳ 等待服务启动..."
	@sleep 10
	@echo "✅ 服务启动完成！"
	@$(MAKE) docker-info

# 仅启动依赖服务 (MySQL + Redis)
.PHONY: docker-deps
docker-deps:
	@echo "🔧 启动依赖服务 (MySQL + Redis)..."
	@cd deployments && docker-compose up -d mysql redis
	@echo "⏳ 等待服务启动..."
	@sleep 5
	@echo "✅ 依赖服务启动完成！"
	@echo ""
	@echo "📍 MySQL: localhost:3306"
	@echo "📍 Redis: localhost:6379"
	@echo ""
	@echo "💡 现在可以在本地运行应用："
	@echo "   go run main.go server --config=./configs/config.dev.yaml"

# Docker Compose 启动
.PHONY: docker-up
docker-up:
	@echo "🐳 启动 Docker Compose..."
	cd deployments && docker-compose up -d

# Docker Compose 停止
.PHONY: docker-down
docker-down:
	@echo "🐳 停止 Docker Compose..."
	cd deployments && docker-compose down

# 停止所有服务
.PHONY: docker-stop
docker-stop:
	@echo "🛑 停止所有服务..."
	@cd deployments && docker-compose down 2>/dev/null || true
	@cd deployments && docker-compose -f docker-compose.simple.yml down 2>/dev/null || true
	@echo "✅ 所有服务已停止"

# 查看服务状态
.PHONY: docker-status
docker-status:
	@echo "📊 服务状态："
	@cd deployments && docker-compose ps

# 查看应用日志
.PHONY: docker-logs
docker-logs:
	@echo "📋 应用日志："
	@cd deployments && docker-compose logs -f app

# 查看所有日志
.PHONY: docker-logs-all
docker-logs-all:
	@echo "📋 所有服务日志："
	@cd deployments && docker-compose logs -f

# 重启应用
.PHONY: docker-restart
docker-restart:
	@echo "🔄 重启应用..."
	@cd deployments && docker-compose restart app
	@echo "✅ 重启完成"

# 重启所有服务
.PHONY: docker-restart-all
docker-restart-all:
	@echo "🔄 重启所有服务..."
	@cd deployments && docker-compose restart
	@echo "✅ 重启完成"

# 清理所有数据（危险操作）
.PHONY: docker-clean
docker-clean:
	@echo "⚠️  警告：此操作将删除所有容器、镜像和数据卷！"
	@echo "⚠️  所有数据库数据将被永久删除！"
	@read -p "确定要继续吗？(输入 'yes' 确认): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🧹 停止并删除所有容器..."; \
		cd deployments && docker-compose down -v; \
		cd deployments && docker-compose -f docker-compose.simple.yml down -v 2>/dev/null || true; \
		echo "🗑️  删除应用镜像..."; \
		docker rmi go-demo:latest 2>/dev/null || true; \
		docker rmi deployments-app 2>/dev/null || true; \
		docker rmi deployments_app 2>/dev/null || true; \
		echo "🧹 清理未使用的资源..."; \
		docker system prune -f; \
		echo "✅ 清理完成"; \
	else \
		echo "❌ 操作已取消"; \
	fi

# ==================== Podman 相关命令 ====================

# Podman 构建
.PHONY: podman-build
podman-build:
	@echo "🦭 构建 Podman 镜像..."
	podman build -f deployments/Dockerfile -t $(PROJECT_NAME):$(VERSION) .
	podman tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest

# Podman 运行
.PHONY: podman-run
podman-run:
	@echo "🦭 运行 Podman 容器..."
	podman run -p 8080:8080 $(PROJECT_NAME):latest

# Podman Compose 完整部署 (应用 + MySQL + Redis + Nginx)
.PHONY: podman-deploy
podman-deploy:
	@echo "🚀 开始 Podman 完整部署..."
	@cd deployments && podman compose pull
	@cd deployments && podman compose build
	@cd deployments && podman compose up -d
	@echo "⏳ 等待服务启动..."
	@sleep 10
	@echo "✅ 服务启动完成！"
	@$(MAKE) podman-info

# Podman Compose 简化部署 (应用 + MySQL + Redis)
.PHONY: podman-deploy-simple
podman-deploy-simple:
	@echo "🚀 开始 Podman 简化部署..."
	@cd deployments && podman compose -f docker-compose.simple.yml pull
	@cd deployments && podman compose -f docker-compose.simple.yml build
	@cd deployments && podman compose -f docker-compose.simple.yml up -d
	@echo "⏳ 等待服务启动..."
	@sleep 10
	@echo "✅ 服务启动完成！"
	@$(MAKE) podman-info

# Podman 仅启动依赖服务 (MySQL + Redis)
.PHONY: podman-deps
podman-deps:
	@echo "🔧 Podman 启动依赖服务 (MySQL + Redis)..."
	@cd deployments && podman compose up -d mysql redis
	@echo "⏳ 等待服务启动..."
	@sleep 5
	@echo "✅ 依赖服务启动完成！"
	@echo ""
	@echo "📍 MySQL: localhost:3306"
	@echo "📍 Redis: localhost:6379"
	@echo ""
	@echo "💡 现在可以在本地运行应用："
	@echo "   go run main.go server --config=./configs/config.dev.yaml"

# Podman Compose 启动
.PHONY: podman-up
podman-up:
	@echo "🦭 启动 Podman Compose..."
	cd deployments && podman compose up -d

# Podman Compose 停止
.PHONY: podman-down
podman-down:
	@echo "🦭 停止 Podman Compose..."
	cd deployments && podman compose down

# Podman 停止所有服务
.PHONY: podman-stop
podman-stop:
	@echo "🛑 Podman 停止所有服务..."
	@cd deployments && podman compose down 2>/dev/null || true
	@cd deployments && podman compose -f docker-compose.simple.yml down 2>/dev/null || true
	@echo "✅ 所有服务已停止"

# Podman 查看服务状态
.PHONY: podman-status
podman-status:
	@echo "📊 Podman 服务状态："
	@cd deployments && podman compose ps

# Podman 查看应用日志
.PHONY: podman-logs
podman-logs:
	@echo "📋 Podman 应用日志："
	@cd deployments && podman compose logs -f app

# Podman 查看所有日志
.PHONY: podman-logs-all
podman-logs-all:
	@echo "📋 Podman 所有服务日志："
	@cd deployments && podman compose logs -f

# Podman 重启应用
.PHONY: podman-restart
podman-restart:
	@echo "🔄 Podman 重启应用..."
	@cd deployments && podman compose restart app
	@echo "✅ 重启完成"

# Podman 重启所有服务
.PHONY: podman-restart-all
podman-restart-all:
	@echo "🔄 Podman 重启所有服务..."
	@cd deployments && podman compose restart
	@echo "✅ 重启完成"

# Podman 清理所有数据（危险操作）
.PHONY: podman-clean
podman-clean:
	@echo "⚠️  警告：此操作将删除所有容器、镜像和数据卷！"
	@echo "⚠️  所有数据库数据将被永久删除！"
	@read -p "确定要继续吗？(输入 'yes' 确认): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🧹 停止并删除所有容器..."; \
		cd deployments && podman compose down -v; \
		cd deployments && podman compose -f docker-compose.simple.yml down -v 2>/dev/null || true; \
		echo "🗑️  删除应用镜像..."; \
		podman rmi go-demo:latest 2>/dev/null || true; \
		podman rmi deployments-app 2>/dev/null || true; \
		podman rmi deployments_app 2>/dev/null || true; \
		echo "🧹 清理未使用的资源..."; \
		podman system prune -f; \
		echo "✅ 清理完成"; \
	else \
		echo "❌ 操作已取消"; \
	fi

# Podman 显示服务信息
.PHONY: podman-info
podman-info:
	@echo ""
	@echo "════════════════════════════════════════════════════════"
	@echo "✅ Podman 服务访问地址："
	@echo "  • 应用 API:      http://localhost:8080"
	@echo "  • Nginx 代理:    http://localhost"
	@echo "  • Swagger 文档:  http://localhost:8080/swagger/index.html"
	@echo "  • 健康检查:      http://localhost:8080/health"
	@echo ""
	@echo "✅ 数据库连接信息："
	@echo "  • MySQL:         localhost:3306"
	@echo "  • Redis:         localhost:6379"
	@echo "════════════════════════════════════════════════════════"
	@echo ""
	@echo "💡 测试服务："
	@echo "   curl http://localhost:8080/health"
	@echo ""
	@echo "💡 查看日志："
	@echo "   make podman-logs"
	@echo ""

# 显示服务信息
.PHONY: docker-info
docker-info:
	@echo ""
	@echo "════════════════════════════════════════════════════════"
	@echo "✅ 服务访问地址："
	@echo "  • 应用 API:      http://localhost:8080"
	@echo "  • Nginx 代理:    http://localhost"
	@echo "  • Swagger 文档:  http://localhost:8080/swagger/index.html"
	@echo "  • 健康检查:      http://localhost:8080/health"
	@echo ""
	@echo "✅ 数据库连接信息："
	@echo "  • MySQL:         localhost:3306"
	@echo "  • Redis:         localhost:6379"
	@echo "════════════════════════════════════════════════════════"
	@echo ""
	@echo "💡 测试服务："
	@echo "   curl http://localhost:8080/health"
	@echo ""
	@echo "💡 查看日志："
	@echo "   make docker-logs"
	@echo ""

# 生成 API 文档
.PHONY: docs
docs:
	@echo "📚 生成 API 文档..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH); \
	else \
		echo "请先安装 swag: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 数据库迁移
.PHONY: migrate
migrate:
	@echo "🗄️ 数据库迁移..."
	$(GOCMD) run $(MAIN_PATH) -migrate

# 健康检查
.PHONY: health
health:
	@echo "🏥 健康检查..."
	@curl -f http://localhost:8080/health || echo "服务未运行"

# 显示帮助信息
.PHONY: help
help:
	@echo "════════════════════════════════════════════════════════"
	@echo "           Go Demo 项目 Makefile 帮助文档"
	@echo "════════════════════════════════════════════════════════"
	@echo ""
	@echo "📦 后端基础命令:"
	@echo "  all              - 执行完整的构建流程 (clean + deps + fmt + vet + test + build)"
	@echo "  deps             - 安装后端依赖"
	@echo "  fmt              - 格式化代码"
	@echo "  vet              - 代码检查"
	@echo "  test             - 运行测试"
	@echo "  test-coverage    - 运行测试并生成覆盖率报告"
	@echo "  build            - 构建后端应用"
	@echo "  build-all        - 构建多平台版本"
	@echo "  run              - 运行后端应用"
	@echo "  dev              - 后端开发模式运行（热重载）"
	@echo "  clean            - 清理构建文件"
	@echo ""
	@echo "🌐 前端命令:"
	@echo "  web-install      - 安装前端依赖"
	@echo "  web-dev          - 前端开发模式运行"
	@echo "  web-build        - 构建前端"
	@echo "  web-preview      - 预览前端构建"
	@echo "  web-lint         - 前端代码检查"
	@echo "  web-clean        - 清理前端构建"
	@echo ""
	@echo "🚀 全栈开发命令:"
	@echo "  dev-all          - 同时启动前后端开发服务器"
	@echo "  install-all      - 安装所有依赖（前端+后端）"
	@echo "  build-all-stack  - 构建前后端"
	@echo ""
	@echo "🛠️  开发工具:"
	@echo "  install-tools    - 安装开发工具"
	@echo "  lint             - 代码质量检查"
	@echo "  docs             - 生成 API 文档"
	@echo "  migrate          - 数据库迁移"
	@echo "  health           - 健康检查"
	@echo ""
	@echo "🐳 Docker 基础命令:"
	@echo "  docker-build     - 构建 Docker 镜像"
	@echo "  docker-run       - 运行 Docker 容器"
	@echo "  docker-up        - 启动 Docker Compose"
	@echo "  docker-down      - 停止 Docker Compose"
	@echo ""
	@echo "🚀 Docker 快速部署:"
	@echo "  docker-deploy         - 完整部署 (应用 + MySQL + Redis + Nginx)"
	@echo "  docker-deploy-simple  - 简化部署 (应用 + MySQL + Redis)"
	@echo "  docker-deps           - 仅启动依赖服务 (MySQL + Redis)"
	@echo ""
	@echo "🔧 Docker 管理命令:"
	@echo "  docker-stop           - 停止所有服务"
	@echo "  docker-status         - 查看服务状态"
	@echo "  docker-logs           - 查看应用日志"
	@echo "  docker-logs-all       - 查看所有服务日志"
	@echo "  docker-restart        - 重启应用"
	@echo "  docker-restart-all    - 重启所有服务"
	@echo "  docker-clean          - 清理所有数据（危险操作）"
	@echo "  docker-info           - 显示服务信息"
	@echo ""
	@echo "🦭 Podman 基础命令:"
	@echo "  podman-build     - 构建 Podman 镜像"
	@echo "  podman-run       - 运行 Podman 容器"
	@echo "  podman-up        - 启动 Podman Compose"
	@echo "  podman-down      - 停止 Podman Compose"
	@echo ""
	@echo "🚀 Podman 快速部署:"
	@echo "  podman-deploy         - 完整部署 (应用 + MySQL + Redis + Nginx)"
	@echo "  podman-deploy-simple  - 简化部署 (应用 + MySQL + Redis)"
	@echo "  podman-deps           - 仅启动依赖服务 (MySQL + Redis)"
	@echo ""
	@echo "🔧 Podman 管理命令:"
	@echo "  podman-stop           - 停止所有服务"
	@echo "  podman-status         - 查看服务状态"
	@echo "  podman-logs           - 查看应用日志"
	@echo "  podman-logs-all       - 查看所有服务日志"
	@echo "  podman-restart        - 重启应用"
	@echo "  podman-restart-all    - 重启所有服务"
	@echo "  podman-clean          - 清理所有数据（危险操作）"
	@echo "  podman-info           - 显示服务信息"
	@echo ""
	@echo "💡 快速开始 (Docker):"
	@echo "  1. 完整部署:    make docker-deploy"
	@echo "  2. 查看状态:    make docker-status"
	@echo "  3. 查看日志:    make docker-logs"
	@echo "  4. 健康检查:    make health"
	@echo "  5. 停止服务:    make docker-stop"
	@echo ""
	@echo "💡 快速开始 (Podman):"
	@echo "  1. 完整部署:    make podman-deploy"
	@echo "  2. 查看状态:    make podman-status"
	@echo "  3. 查看日志:    make podman-logs"
	@echo "  4. 健康检查:    make health"
	@echo "  5. 停止服务:    make podman-stop"
	@echo ""
	@echo "  help             - 显示此帮助信息"
	@echo "════════════════════════════════════════════════════════"
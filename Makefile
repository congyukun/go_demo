# Go Demo 项目 Makefile

# 项目信息
PROJECT_NAME := go-demo
VERSION := 1.0.0
BUILD_DIR := bin
MAIN_PATH := cmd/server/main.go

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# 构建标志
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(shell date '+%Y-%m-%d %H:%M:%S')' -X 'main.GitCommit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)'"

# 默认目标
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
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		$(GOCMD) run $(MAIN_PATH); \
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
	@echo "Go Demo 项目 Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  all           - 执行完整的构建流程 (clean + deps + fmt + vet + test + build)"
	@echo "  deps          - 安装依赖"
	@echo "  fmt           - 格式化代码"
	@echo "  vet           - 代码检查"
	@echo "  test          - 运行测试"
	@echo "  test-coverage - 运行测试并生成覆盖率报告"
	@echo "  build         - 构建应用"
	@echo "  build-all     - 构建多平台版本"
	@echo "  run           - 运行应用"
	@echo "  dev           - 开发模式运行（热重载）"
	@echo "  clean         - 清理构建文件"
	@echo "  install-tools - 安装开发工具"
	@echo "  lint          - 代码质量检查"
	@echo "  docker-build  - 构建 Docker 镜像"
	@echo "  docker-run    - 运行 Docker 容器"
	@echo "  docker-up     - 启动 Docker Compose"
	@echo "  docker-down   - 停止 Docker Compose"
	@echo "  docs          - 生成 API 文档"
	@echo "  migrate       - 数据库迁移"
	@echo "  health        - 健康检查"
	@echo "  help          - 显示此帮助信息"
#!/bin/bash

# 项目重构迁移脚本
set -e

echo "开始项目结构重构迁移..."

# 创建备份目录
BACKUP_DIR="backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "创建备份目录: $BACKUP_DIR"

# 备份旧文件
echo "备份旧的项目文件..."

# 备份旧的目录结构
if [ -d "controllers" ]; then
    cp -r controllers "$BACKUP_DIR/"
    echo "已备份 controllers/ 目录"
fi

if [ -d "services" ]; then
    cp -r services "$BACKUP_DIR/"
    echo "已备份 services/ 目录"
fi

if [ -d "models" ]; then
    cp -r models "$BACKUP_DIR/"
    echo "已备份 models/ 目录"
fi

if [ -d "routes" ]; then
    cp -r routes "$BACKUP_DIR/"
    echo "已备份 routes/ 目录"
fi

if [ -d "config" ]; then
    cp -r config "$BACKUP_DIR/"
    echo "已备份 config/ 目录"
fi

if [ -d "db" ]; then
    cp -r db "$BACKUP_DIR/"
    echo "已备份 db/ 目录"
fi

if [ -d "logger" ]; then
    cp -r logger "$BACKUP_DIR/"
    echo "已备份 logger/ 目录"
fi

if [ -d "utils" ]; then
    cp -r utils "$BACKUP_DIR/"
    echo "已备份 utils/ 目录"
fi

if [ -d "registry" ]; then
    cp -r registry "$BACKUP_DIR/"
    echo "已备份 registry/ 目录"
fi

if [ -f "main.go" ]; then
    cp main.go "$BACKUP_DIR/"
    echo "已备份 main.go"
fi

if [ -f "go_demo" ]; then
    cp go_demo "$BACKUP_DIR/"
    echo "已备份 go_demo 二进制文件"
fi

# 删除旧的目录结构
echo "删除旧的目录结构..."

# 删除旧目录
rm -rf controllers/
rm -rf services/
rm -rf models/
rm -rf routes/
rm -rf config/
rm -rf db/
rm -rf logger/
rm -rf utils/
rm -rf registry/
rm -rf tools/

# 删除旧文件
rm -f main.go
rm -f go_demo

echo "已删除旧的目录和文件"

# 创建缺少的新目录
echo "创建标准项目目录结构..."

# 创建 pkg 目录结构
mkdir -p pkg/database
mkdir -p pkg/logger

# 创建其他必要目录
mkdir -p data
mkdir -p bin

# 设置脚本执行权限
chmod +x scripts/*.sh

echo "项目结构重构完成！"
echo ""
echo "📁 新的项目结构："
echo "├── cmd/                    # 应用程序入口"
echo "├── internal/              # 内部应用代码"
echo "├── pkg/                   # 可重用的库代码"
echo "├── configs/              # 配置文件"
echo "├── api/                  # API 文档"
echo "├── docs/                 # 项目文档"
echo "├── scripts/              # 脚本文件"
echo "├── tests/                # 测试文件"
echo "├── deployments/          # 部署配置"
echo "├── logs/                 # 日志文件"
echo "├── data/                 # 数据文件"
echo "└── bin/                  # 编译输出"
echo ""
echo "🔄 旧文件已备份到: $BACKUP_DIR"
echo "📝 请检查新的配置文件: configs/config.yaml"
echo "🚀 运行新的应用: go run cmd/server/main.go"
echo ""
echo "重构完成！"
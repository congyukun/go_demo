#!/bin/bash

# 构建脚本
set -e

# 项目信息
PROJECT_NAME="go-demo"
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建目录
BUILD_DIR="bin"
mkdir -p ${BUILD_DIR}

# 构建标志
LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}'"

echo "开始构建 ${PROJECT_NAME}..."
echo "版本: ${VERSION}"
echo "构建时间: ${BUILD_TIME}"
echo "Git提交: ${GIT_COMMIT}"

# 构建不同平台的二进制文件
echo "构建 Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-linux-amd64 cmd/server/main.go

echo "构建 Darwin amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-darwin-amd64 cmd/server/main.go

echo "构建 Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME}-windows-amd64.exe cmd/server/main.go

echo "构建完成！"
echo "输出目录: ${BUILD_DIR}/"
ls -la ${BUILD_DIR}/
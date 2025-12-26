# 多阶段构建
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 代理（国内镜像加速）
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata make

# 复制 go mod 文件（利用 Docker 缓存层）
COPY go.mod go.sum ./

# 下载依赖（这一层会被缓存，除非 go.mod/go.sum 改变）
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 单独复制启动脚本（因为 deployments 在 .dockerignore 中）
COPY deployments/wait-for-services.sh ./wait-for-services.sh

# 构建应用（添加编译优化）
RUN go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o main .

# 运行阶段 - 使用更小的基础镜像
FROM alpine:3.19

# 安装运行时依赖（包括健康检查和等待脚本所需工具）
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    wget \
    netcat-openbsd \
    mysql-client \
    redis

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和启动脚本
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/wait-for-services.sh /usr/local/bin/wait-for-services.sh

# 创建日志目录并设置权限
RUN mkdir -p logs && \
    chmod +x /usr/local/bin/wait-for-services.sh && \
    chown -R appuser:appgroup /app && \
    chmod -R 755 /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 优化的健康检查（减少启动等待时间）
HEALTHCHECK --interval=15s --timeout=5s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 默认命令（使用 shell 形式确保参数正确传递）
CMD /usr/local/bin/wait-for-services.sh ./main server --config=/app/configs/config.docker.yaml
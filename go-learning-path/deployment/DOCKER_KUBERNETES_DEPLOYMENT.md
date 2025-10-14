# 🐳 Docker 和 Kubernetes 部署指南

## 📦 Docker 容器化

### 1. 多阶段构建 Dockerfile
```dockerfile
# 第一阶段：构建阶段
FROM golang:1.21-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server

# 第二阶段：运行阶段
FROM alpine:3.18

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder --chown=appuser:appuser /app/main .

# 复制配置文件
COPY --chown=appuser:appuser configs/ configs/

# 切换用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动命令
CMD ["./main"]
```

### 2. Docker Compose 开发环境
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=appuser
      - DB_PASSWORD=apppassword
      - DB_NAME=appdb
    depends_on:
      mysql:
        condition: service_healthy
    volumes:
      - .:/app
      - go-modules:/go/pkg/mod
    develop:
      watch:
        - action: sync
          path: .
          target: /app
        - action: rebuild
          path: go.mod

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=appdb
      - MYSQL_USER=appuser
      - MYSQL_PASSWORD=apppassword
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
      - ./deployments/mysql/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes

volumes:
  mysql-data:
  redis-data:
  go-modules:
```

### 3. 生产环境 Dockerfile
```dockerfile
# 使用distroless镜像提供更小的攻击面
FROM gcr.io/distroless/base-debian11

# 设置元数据
LABEL maintainer="your-team@company.com"
LABEL version="1.0.0"
LABEL description="Go application with GORM"

# 设置工作目录
WORKDIR /app

# 复制二进制文件和配置文件
COPY --chown=nonroot:nonroot main .
COPY --chown=nonroot:nonroot configs/ configs/

# 使用非root用户
USER nonroot

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["/app/main"]
```

## ☸️ Kubernetes 部署

### 1. Deployment 配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
  namespace: production
  labels:
    app: go-app
    environment: production
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
        environment: production
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
    spec:
      containers:
      - name: go-app
        image: your-registry/go-app:1.0.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
          protocol: TCP
        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: host
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: password
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
        volumeMounts:
        - name: config-volume
          mountPath: /app/configs
      volumes:
      - name: config-volume
        configMap:
          name: app-config
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - go-app
              topologyKey: kubernetes.io/hostname
---
apiVersion: v1
kind: Service
metadata:
  name: go-app-service
  namespace: production
spec:
  selector:
    app: go-app
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
  type: ClusterIP
```

### 2. ConfigMap 和 Secret
```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: production
data:
  config.yaml: |
    server:
      port: 8080
      mode: release
    database:
      host: mysql.production.svc.cluster.local
      port: 3306
      name: appdb
    redis:
      host: redis.production.svc.cluster.local
      port: 6379

# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-secret
  namespace: production
type: Opaque
data:
  host: bXlzcWwucHJvZHVjdGlvbi5zdmMuY2x1c3Rlci5sb2NhbA==  # base64 encoded
  username: YXBwdXNlcg==                                # appuser
  password: cGFzc3dvcmQxMjM=                            # password123
```

### 3. Horizontal Pod Autoscaler
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-app-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      policies:
      - type: Pods
        value: 2
        periodSeconds: 60
      - type: Percent
        value: 50
        periodSeconds: 60
      selectPolicy: Max
      stabilizationWindowSeconds: 0
    scaleDown:
      policies:
      - type: Pods
        value: 1
        periodSeconds: 60
      selectPolicy: Max
      stabilizationWindowSeconds: 300
```

## 🔧 部署脚本和工具

### 1. 部署脚本
```bash
#!/bin/bash
# deploy.sh

set -e

# 环境变量
ENV=${1:-staging}
VERSION=${2:-latest}
NAMESPACE=${3:-$ENV}

echo "🚀 开始部署到 $ENV 环境，版本: $VERSION"

# 构建Docker镜像
echo "📦 构建Docker镜像..."
docker build -t your-registry/go-app:$VERSION .

# 推送镜像
echo "📤 推送镜像到仓库..."
docker push your-registry/go-app:$VERSION

# 更新Kubernetes部署
echo "🔄 更新Kubernetes部署..."
kubectl set image deployment/go-app go-app=your-registry/go-app:$VERSION -n $NAMESPACE

# 等待部署完成
echo "⏳ 等待部署完成..."
kubectl rollout status deployment/go-app -n $NAMESPACE --timeout=300s

echo "✅ 部署完成!"
```

### 2. GitLab CI/CD 配置
```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  APP_NAME: go-app
  DOCKER_REGISTRY: your-registry

.test-template: &test-template
  image: golang:1.21
  before_script:
    - go version
    - go mod download
  script:
    - go test -v ./... -coverprofile=coverage.out
    - go tool cover -func=coverage.out

test:
  <<: *test-template
  stage: test

build:
  stage: build
  image: docker:20.10
  services:
    - docker:20.10-dind
  script:
    - docker build -t $DOCKER_REGISTRY/$APP_NAME:$CI_COMMIT_SHA .
    - docker push $DOCKER_REGISTRY/$APP_NAME:$CI_COMMIT_SHA
  only:
    - main
    - develop

deploy-staging:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context staging
    - kubectl set image deployment/$APP_NAME $APP_NAME=$DOCKER_REGISTRY/$APP_NAME:$CI_COMMIT_SHA -n staging
    - kubectl rollout status deployment/$APP_NAME -n staging --timeout=300s
  environment:
    name: staging
    url: https://staging.example.com
  only:
    - develop

deploy-production:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context production
    - kubectl set image deployment/$APP_NAME $APP_NAME=$DOCKER_REGISTRY/$APP_NAME:$CI_COMMIT_SHA -n production
    - kubectl rollout status deployment/$APP_NAME -n production --timeout=300s
  environment:
    name: production
    url: https://example.com
  only:
    - main
  when: manual
```

## 📊 监控和日志

### 1. Prometheus 监控配置
```yaml
# prometheus-rules.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: go-app-rules
  namespace: monitoring
spec:
  groups:
  - name: go-app
    rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "高错误率警告"
        description: "应用错误率超过5%"

    - alert: HighLatency
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
      for: 10m
      labels:
        severity: warning
      annotations:
        summary: "高延迟警告"
        description: "95%的请求延迟超过1秒"
```

### 2. Grafana 仪表板
```json
{
  "dashboard": {
    "title": "Go应用监控",
    "panels": [
      {
        "title": "请求率",
        "type": "graph",
        "targets": [{
          "expr": "rate(http_requests_total[5m])",
          "legendFormat": "{{handler}}"
        }]
      },
      {
        "title": "错误率",
        "type": "graph",
        "targets": [{
          "expr": "rate(http_requests_total{status=~'5..'}[5m]) / rate(http_requests_total[5m])",
          "legendFormat": "错误率"
        }]
      }
    ]
  }
}
```

## 🚀 最佳实践

### 1. 安全最佳实践
- 使用非root用户运行容器
- 定期更新基础镜像和安全补丁
- 扫描镜像中的漏洞
- 使用网络策略限制网络访问

### 2. 性能最佳实践
- 合理设置资源请求和限制
- 使用就绪性和存活性探针
- 实现优雅关闭
- 使用连接池和缓存

### 3. 可观察性最佳实践
- 实现完整的日志记录
- 设置监控和告警
- 使用分布式追踪
- 收集业务指标

## 📋 部署检查清单

- [ ] Docker镜像构建和测试
- [ ] Kubernetes资源配置验证
- [ ] 监控和告警设置
- [ ] 备份和恢复策略
- [ ] 安全扫描完成
- [ ] 性能测试通过

下一步学习：**深入研究Go设计模式和最佳实践**
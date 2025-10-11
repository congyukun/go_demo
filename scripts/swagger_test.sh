#!/bin/bash

# Swagger文档测试脚本

echo "=== Go Demo API Swagger文档测试 ==="
echo

# 检查服务是否运行
if ! nc -z localhost 8080 2>/dev/null; then
    echo "⚠️  服务未在8080端口运行，尝试启动..."
    
    # 检查端口占用
    PORT_IN_USE=$(lsof -ti:8080 || echo "")
    if [ -n "$PORT_IN_USE" ]; then
        echo "❌ 端口8080被占用 (PID: $PORT_IN_USE)"
        echo "请执行: kill -9 $PORT_IN_USE"
        exit 1
    fi
    
    echo "请手动启动服务: make run 或 go run cmd/server/main.go"
    exit 1
fi

echo "✅ 服务正在运行"

# 测试Swagger文档访问
echo
echo "=== 测试Swagger文档访问 ==="

# 测试Swagger UI
echo "测试 Swagger UI 访问..."
SWAGGER_UI=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/swagger/index.html)
if [ "$SWAGGER_UI" -eq 200 ]; then
    echo "✅ Swagger UI 可访问: http://localhost:8080/swagger/index.html"
else
    echo "❌ Swagger UI 访问失败 (HTTP $SWAGGER_UI)"
fi

# 测试Swagger JSON
echo "测试 Swagger JSON 文档..."
SWAGGER_JSON=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/swagger/doc.json)
if [ "$SWAGGER_JSON" -eq 200 ]; then
    echo "✅ Swagger JSON 文档可访问: http://localhost:8080/swagger/doc.json"
else
    echo "❌ Swagger JSON 文档访问失败 (HTTP $SWAGGER_JSON)"
fi

# 测试API端点
echo
echo "=== 测试API端点 ==="

# 测试健康检查
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$HEALTH" -eq 200 ]; then
    echo "✅ 健康检查端点可访问: http://localhost:8080/health"
else
    echo "❌ 健康检查端点访问失败 (HTTP $HEALTH)"
fi

# 测试认证端点
AUTH_LOGIN=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"test","password":"test"}')
echo "认证登录端点状态: HTTP $AUTH_LOGIN"

echo
echo "=== 文档访问链接 ==="
echo "📖 Swagger UI: http://localhost:8080/swagger/index.html"
echo "📋 Swagger JSON: http://localhost:8080/swagger/doc.json"
echo "🏥 健康检查: http://localhost:8080/health"
echo
echo "=== 使用说明 ==="
echo "1. 打开浏览器访问: http://localhost:8080/swagger/index.html"
echo "2. 点击 'Authorize' 按钮，输入: Bearer <your_jwt_token>"
echo "3. 测试各个API接口"
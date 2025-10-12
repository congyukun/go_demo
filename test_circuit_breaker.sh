#!/bin/bash

echo "开始熔断器测试..."

# 1. 测试基本熔断功能
echo "=== 测试基本熔断功能 ==="
echo "发送失败请求触发熔断..."

# 首先发送一些失败请求来触发熔断
for i in {1..5}; do
  # 模拟失败请求（使用不存在的用户登录）
  curl -s --noproxy "*" -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "username": "nonexistent_user_'$i'",
      "password": "wrong_password"
    }' > /dev/null
done

# 等待一小段时间让请求被记录
sleep 1

# 发送正常请求，应该被熔断
echo "发送正常请求，应该被熔断..."
circuit_response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)

if [[ "$circuit_response" == *"服务暂时不可用"* ]]; then
  echo "✓ 熔断器正常工作，请求被正确拒绝"
elif [[ "$circuit_response" == *"ok"* ]]; then
  echo "✗ 熔断器可能未正常工作，请求被允许通过"
else
  echo "? 无法确定熔断器状态，响应: $circuit_response"
fi

# 等待超时时间，让熔断器进入半开状态
echo "等待30秒超时时间..."
sleep 30

# 2. 测试半开状态
echo "=== 测试半开状态 ==="
echo "发送请求到半开状态的熔断器..."

half_open_response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)

if [[ "$half_open_response" == *"ok"* ]]; then
  echo "✓ 半开状态允许请求通过"
elif [[ "$half_open_response" == *"服务暂时不可用"* ]]; then
  echo "✗ 半开状态仍然拒绝请求"
else
  echo "? 无法确定半开状态，响应: $half_open_response"
fi

# 3. 测试熔断器恢复
echo "=== 测试熔断器恢复 ==="
echo "发送成功请求让熔断器恢复..."

# 发送一些成功请求
for i in {1..5}; do
  curl -s --noproxy "*" -X GET http://localhost:8080/health > /dev/null
  sleep 0.5
done

# 检查熔断器是否恢复
recovery_response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)

if [[ "$recovery_response" == *"ok"* ]]; then
  echo "✓ 熔断器已恢复正常工作"
else
  echo "✗ 熔断器恢复可能有问题"
fi

# 4. 测试并发场景
echo "=== 测试并发场景 ==="
echo "并发发送多个请求..."

# 并发发送10个请求
for i in {1..10}; do
  (
    response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)
    if [[ "$response" == *"ok"* ]]; then
      echo "请求 $i: 成功"
    elif [[ "$response" == *"服务暂时不可用"* ]]; then
      echo "请求 $i: 被熔断"
    else
      echo "请求 $i: 未知响应"
    fi
  ) &
done

wait # 等待所有并发请求完成

# 5. 测试API级别熔断器
echo "=== 测试API级别熔断器 ==="
echo "测试特定API的熔断功能..."

# 测试认证API
for i in {1..5}; do
  curl -s --noproxy "*" -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "username": "test_user_'$i'",
      "password": "wrong_password"
    }' > /dev/null
done

# 测试用户API
api_response=$(curl -s --noproxy "*" -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer invalid_token")

if [[ "$api_response" == *"服务暂时不可用"* ]]; then
  echo "✓ API级别熔断器正常工作"
else
  echo "✗ API级别熔断器可能未生效"
fi

echo "=== 熔断器测试完成 ==="
echo "测试包括："
echo "1. 基本熔断功能"
echo "2. 半开状态处理"
echo "3. 熔断器恢复"
echo "4. 并发场景"
echo "5. API级别熔断"
echo ""
echo "注意：熔断器配置为错误率50%时触发，超时时间30秒"
echo "半开状态最大请求数：10个"
echo "触发熔断检查的最小请求数：100个"
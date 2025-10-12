#!/bin/bash

echo "开始综合限流测试..."

# 1. 测试IP级限流
echo "=== 测试IP级限流 ==="
echo "发送210个请求到/health端点（IP限流：200个/分钟）..."

ip_success_count=0
ip_rate_limited_count=0

for i in {1..210}; do
  response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)
  
  if [[ "$response" == *"请求过于频繁"* ]]; then
    ip_rate_limited_count=$((ip_rate_limited_count + 1))
  elif [[ "$response" == *"ok"* ]] || [[ "$response" == *"status"* ]]; then
    ip_success_count=$((ip_success_count + 1))
  fi
  
  # 添加小延迟以避免过快请求
  sleep 0.05
done

echo "IP级限流测试结果："
echo "成功请求数: $ip_success_count"
echo "被限流请求数: $ip_rate_limited_count"

# 等待一段时间让限流重置
echo "等待60秒让限流重置..."
sleep 60

# 2. 测试用户级限流
echo "=== 测试用户级限流 ==="
echo "注册测试用户..."
register_response=$(curl -s --noproxy "*" -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "comprehensive_user_'$(date +%s)'",
    "password": "password123",
    "email": "comprehensive_user_'$(date +%s)'@example.com",
    "name": "综合测试用户",
    "mobile": "13800138777"
  }')

# 提取用户名
username=$(echo "$register_response" | grep -o '"username":"[^"]*' | sed 's/"username":"//')

if [ -z "$username" ]; then
  username="comprehensive_user_$(date +%s)"
fi

echo "登录获取token..."
login_response=$(curl -s --noproxy "*" -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'$username'",
    "password": "password123"
  }')

# 提取token
token=$(echo "$login_response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$token" ]; then
  echo "无法获取token，跳过用户级限流测试"
  user_success_count=0
  user_rate_limited_count=0
else
  echo "发送110个请求到/auth/profile端点（用户限流：100个/分钟）..."
  
  user_success_count=0
  user_rate_limited_count=0
  
  for i in {1..110}; do
    response=$(curl -s --noproxy "*" -X GET http://localhost:8080/api/v1/auth/profile \
      -H "Authorization: Bearer $token")
    
    if [[ "$response" == *"请求过于频繁"* ]]; then
      user_rate_limited_count=$((user_rate_limited_count + 1))
    elif [[ "$response" == *"获取成功"* ]]; then
      user_success_count=$((user_success_count + 1))
    fi
    
    # 添加小延迟以避免过快请求
    sleep 0.1
  done
fi

echo "用户级限流测试结果："
echo "成功请求数: $user_success_count"
echo "被限流请求数: $user_rate_limited_count"

# 等待一段时间让限流重置
echo "等待60秒让限流重置..."
sleep 60

# 3. 测试全局限流
echo "=== 测试全局限流 ==="
echo "发送1010个请求到/health端点（全局限流：1000个/秒）..."

global_success_count=0
global_rate_limited_count=0

for i in {1..1010}; do
  response=$(curl -s --noproxy "*" -X GET http://localhost:8080/health)
  
  if [[ "$response" == *"请求过于频繁"* ]]; then
    global_rate_limited_count=$((global_rate_limited_count + 1))
  elif [[ "$response" == *"ok"* ]] || [[ "$response" == *"status"* ]]; then
    global_success_count=$((global_success_count + 1))
  fi
  
  # 不添加延迟，快速发送请求以测试每秒1000的限制
done

echo "全局限流测试结果："
echo "成功请求数: $global_success_count"
echo "被限流请求数: $global_rate_limited_count"

# 总结
echo "=== 综合测试总结 ==="
echo "IP级限流: $ip_success_count 成功, $ip_rate_limited_count 被限流"
echo "用户级限流: $user_success_count 成功, $user_rate_limited_count 被限流"
echo "全局限流: $global_success_count 成功, $global_rate_limited_count 被限流"

if [ $ip_rate_limited_count -gt 0 ] && [ $user_rate_limited_count -gt 0 ] && [ $global_rate_limited_count -gt 0 ]; then
  echo "所有限流测试均成功！"
else
  echo "部分限流测试可能未成功"
fi
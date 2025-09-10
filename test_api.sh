#!/bin/bash

echo "测试注册API验证功能..."

# 测试1: 密码太短
echo "1. 测试密码太短:"
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "12345",
    "name": "测试用户",
    "mobile": "13812345678"
  }' | jq .

echo -e "\n"

# 测试2: 手机号格式错误
echo "2. 测试手机号格式错误:"
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "name": "测试用户",
    "mobile": "12345"
  }' | jq .

echo -e "\n"

# 测试3: 用户名太短
echo "3. 测试用户名太短:"
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ab",
    "password": "123456",
    "name": "测试用户",
    "mobile": "13812345678"
  }' | jq .

echo -e "\n"

# 测试4: 邮箱格式错误
echo "4. 测试邮箱格式错误:"
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "name": "测试用户",
    "email": "invalid-email",
    "mobile": "13812345678"
  }' | jq .

echo -e "\n"

# 测试5: 正确的请求
echo "5. 测试正确的请求:"
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "password": "123456",
    "name": "测试用户",
    "email": "test@example.com",
    "mobile": "13812345678"
  }' | jq .
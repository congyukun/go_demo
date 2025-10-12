#!/bin/bash

# 获取JWT令牌
TOKEN=$(curl -s -H "Content-Type: application/json" -d '{"username":"testuser","password":"newpassword123"}' http://localhost:8080/api/v1/auth/login | jq -r '.data.token')

echo "获取到的JWT令牌: $TOKEN"

# 测试删除用户接口
echo "测试删除用户接口..."
curl -X DELETE -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/users/3

echo ""
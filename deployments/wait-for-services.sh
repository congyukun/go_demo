#!/bin/sh
# 等待服务就绪脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
MYSQL_HOST="${MYSQL_HOST:-mysql}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_demo}"

REDIS_HOST="${REDIS_HOST:-redis}"
REDIS_PORT="${REDIS_PORT:-6379}"

MAX_RETRIES="${MAX_RETRIES:-30}"
RETRY_INTERVAL="${RETRY_INTERVAL:-2}"

echo "${YELLOW}等待服务就绪...${NC}"

# 等待 MySQL
echo "${YELLOW}检查 MySQL 连接: ${MYSQL_HOST}:${MYSQL_PORT}${NC}"
retry_count=0
until nc -z "$MYSQL_HOST" "$MYSQL_PORT" 2>/dev/null; do
  retry_count=$((retry_count + 1))
  if [ $retry_count -ge $MAX_RETRIES ]; then
    echo "${RED}MySQL 连接超时${NC}"
    exit 1
  fi
  echo "${YELLOW}等待 MySQL 启动... (${retry_count}/${MAX_RETRIES})${NC}"
  sleep $RETRY_INTERVAL
done

# 验证 MySQL 数据库可用
echo "${YELLOW}验证 MySQL 数据库...${NC}"
retry_count=0
until wget --spider -q "http://${MYSQL_HOST}:${MYSQL_PORT}" 2>/dev/null || \
      mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1" >/dev/null 2>&1 || \
      mysqladmin ping -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" >/dev/null 2>&1; do
  retry_count=$((retry_count + 1))
  if [ $retry_count -ge $MAX_RETRIES ]; then
    echo "${RED}MySQL 数据库验证超时${NC}"
    exit 1
  fi
  echo "${YELLOW}等待 MySQL 数据库就绪... (${retry_count}/${MAX_RETRIES})${NC}"
  sleep $RETRY_INTERVAL
done
echo "${GREEN}✓ MySQL 已就绪${NC}"

# 等待 Redis
echo "${YELLOW}检查 Redis 连接: ${REDIS_HOST}:${REDIS_PORT}${NC}"
retry_count=0
until nc -z "$REDIS_HOST" "$REDIS_PORT" 2>/dev/null; do
  retry_count=$((retry_count + 1))
  if [ $retry_count -ge $MAX_RETRIES ]; then
    echo "${RED}Redis 连接超时${NC}"
    exit 1
  fi
  echo "${YELLOW}等待 Redis 启动... (${retry_count}/${MAX_RETRIES})${NC}"
  sleep $RETRY_INTERVAL
done

# 验证 Redis 可用
echo "${YELLOW}验证 Redis 服务...${NC}"
retry_count=0
until redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping >/dev/null 2>&1; do
  retry_count=$((retry_count + 1))
  if [ $retry_count -ge $MAX_RETRIES ]; then
    echo "${RED}Redis 服务验证超时${NC}"
    exit 1
  fi
  echo "${YELLOW}等待 Redis 服务就绪... (${retry_count}/${MAX_RETRIES})${NC}"
  sleep $RETRY_INTERVAL
done
echo "${GREEN}✓ Redis 已就绪${NC}"

echo "${GREEN}所有服务已就绪，启动应用...${NC}"

# 执行传入的命令
exec "$@"

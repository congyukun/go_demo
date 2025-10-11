package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	MaxRetries   int    `yaml:"max_retries"`
}

// RedisCache Redis缓存客户端
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache 创建Redis缓存客户端
func NewRedisCache(config RedisConfig) (*RedisCache, error) {
	// 设置默认值
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set 设置缓存
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	// 设置缓存
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("缓存不存在")
		}
		return "", err
	}
	return val, nil
}

// GetObject 获取缓存并反序列化为对象
func (r *RedisCache) GetObject(key string, dest interface{}) error {
	val, err := r.Get(key)
	if err != nil {
		return err
	}

	// 反序列化
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("反序列化失败: %w", err)
	}

	return nil
}

// Delete 删除缓存
func (r *RedisCache) Delete(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Expire 设置过期时间
func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (r *RedisCache) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Increment 自增
func (r *RedisCache) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// IncrementBy 自增指定值
func (r *RedisCache) IncrementBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, value).Result()
}

// Decrement 自减
func (r *RedisCache) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// DecrementBy 自减指定值
func (r *RedisCache) DecrementBy(key string, value int64) (int64, error) {
	return r.client.DecrBy(r.ctx, key, value).Result()
}

// SetNX 仅在键不存在时设置（用于分布式锁）
func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化失败: %w", err)
	}

	return r.client.SetNX(r.ctx, key, data, expiration).Result()
}

// HSet 设置哈希字段
func (r *RedisCache) HSet(key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return r.client.HSet(r.ctx, key, field, data).Err()
}

// HGet 获取哈希字段
func (r *RedisCache) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HDelete 删除哈希字段
func (r *RedisCache) HDelete(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// LPush 列表左侧插入
func (r *RedisCache) LPush(key string, values ...interface{}) error {
	return r.client.LPush(r.ctx, key, values...).Err()
}

// RPush 列表右侧插入
func (r *RedisCache) RPush(key string, values ...interface{}) error {
	return r.client.RPush(r.ctx, key, values...).Err()
}

// LPop 列表左侧弹出
func (r *RedisCache) LPop(key string) (string, error) {
	return r.client.LPop(r.ctx, key).Result()
}

// RPop 列表右侧弹出
func (r *RedisCache) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

// LRange 获取列表范围
func (r *RedisCache) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// SAdd 集合添加成员
func (r *RedisCache) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisCache) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func (r *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SRemove 移除集合成员
func (r *RedisCache) SRemove(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

// ZAdd 有序集合添加成员
func (r *RedisCache) ZAdd(key string, score float64, member interface{}) error {
	z := redis.Z{
		Score:  score,
		Member: member,
	}
	return r.client.ZAdd(r.ctx, key, &z).Err()
}

// ZRange 获取有序集合范围
func (r *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围（包含分数）
func (r *RedisCache) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(r.ctx, key, start, stop).Result()
}

// FlushDB 清空当前数据库
func (r *RedisCache) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Ping 测试连接
func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// GetClient 获取原始Redis客户端（用于高级操作）
func (r *RedisCache) GetClient() *redis.Client {
	return r.client
}

// 缓存键生成辅助函数

// UserCacheKey 生成用户缓存键
func UserCacheKey(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

// UserListCacheKey 生成用户列表缓存键
func UserListCacheKey(page, size int) string {
	return fmt.Sprintf("users:page:%d:size:%d", page, size)
}

// TokenCacheKey 生成令牌缓存键
func TokenCacheKey(token string) string {
	return fmt.Sprintf("token:%s", token)
}

// SessionCacheKey 生成会话缓存键
func SessionCacheKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

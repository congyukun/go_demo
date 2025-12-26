package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Password     string `mapstructure:"password" yaml:"password"`
	DB           int    `mapstructure:"db" yaml:"db"`
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries" yaml:"max_retries"`
}

// RedisCache Redis缓存客户端
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存客户端
func NewRedisCache(config RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()

	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()
	return c.client.Get(ctx, key).Result()
}

// GetObject 获取对象缓存
func (c *RedisCache) GetObject(key string, dest interface{}) error {
	ctx := context.Background()

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (c *RedisCache) Delete(keys ...string) error {
	ctx := context.Background()
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *RedisCache) Exists(key string) (bool, error) {
	ctx := context.Background()
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Expire 设置过期时间
func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	ctx := context.Background()
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *RedisCache) TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	return c.client.TTL(ctx, key).Result()
}

// Increment 自增（实现 CacheInterface）
func (c *RedisCache) Increment(key string) (int64, error) {
	ctx := context.Background()
	return c.client.Incr(ctx, key).Result()
}

// IncrementBy 增加指定值（实现 CacheInterface）
func (c *RedisCache) IncrementBy(key string, value int64) (int64, error) {
	ctx := context.Background()
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decrement 自减（实现 CacheInterface）
func (c *RedisCache) Decrement(key string) (int64, error) {
	ctx := context.Background()
	return c.client.Decr(ctx, key).Result()
}

// DecrementBy 减少指定值（实现 CacheInterface）
func (c *RedisCache) DecrementBy(key string, value int64) (int64, error) {
	ctx := context.Background()
	return c.client.DecrBy(ctx, key, value).Result()
}

// Incr 自增（向后兼容）
func (c *RedisCache) Incr(key string) (int64, error) {
	return c.Increment(key)
}

// Decr 自减（向后兼容）
func (c *RedisCache) Decr(key string) (int64, error) {
	return c.Decrement(key)
}

// IncrBy 增加指定值（向后兼容）
func (c *RedisCache) IncrBy(key string, value int64) (int64, error) {
	return c.IncrementBy(key, value)
}

// SetNX 设置键值（仅当键不存在时）
func (c *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	ctx := context.Background()

	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化失败: %w", err)
	}

	return c.client.SetNX(ctx, key, data, expiration).Result()
}

// HSet 设置哈希字段
func (c *RedisCache) HSet(key string, field string, value interface{}) error {
	ctx := context.Background()
	return c.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func (c *RedisCache) HGet(key string, field string) (string, error) {
	ctx := context.Background()
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (c *RedisCache) HGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	return c.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (c *RedisCache) HDel(key string, fields ...string) error {
	ctx := context.Background()
	return c.client.HDel(ctx, key, fields...).Err()
}

// HDelete 删除哈希字段（实现 CacheInterface）
func (c *RedisCache) HDelete(key string, fields ...string) error {
	return c.HDel(key, fields...)
}

// SAdd 添加集合成员
func (c *RedisCache) SAdd(key string, members ...interface{}) error {
	ctx := context.Background()
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (c *RedisCache) SMembers(key string) ([]string, error) {
	ctx := context.Background()
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember 判断是否是集合成员
func (c *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	ctx := context.Background()
	return c.client.SIsMember(ctx, key, member).Result()
}

// SRem 移除集合成员
func (c *RedisCache) SRem(key string, members ...interface{}) error {
	ctx := context.Background()
	return c.client.SRem(ctx, key, members...).Err()
}

// LPush 从列表左侧推入元素
func (c *RedisCache) LPush(key string, values ...interface{}) error {
	ctx := context.Background()
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧推入元素
func (c *RedisCache) RPush(key string, values ...interface{}) error {
	ctx := context.Background()
	return c.client.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出元素
func (c *RedisCache) LPop(key string) (string, error) {
	ctx := context.Background()
	return c.client.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func (c *RedisCache) RPop(key string) (string, error) {
	ctx := context.Background()
	return c.client.RPop(ctx, key).Result()
}

// LRange 获取列表指定范围元素
func (c *RedisCache) LRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	return c.client.LRange(ctx, key, start, stop).Result()
}

// ZAdd 添加有序集合成员（实现 CacheInterface）
func (c *RedisCache) ZAdd(key string, members ...*ZMember) error {
	ctx := context.Background()
	zMembers := make([]*redis.Z, len(members))
	for i, m := range members {
		zMembers[i] = &redis.Z{Score: m.Score, Member: m.Member}
	}
	return c.client.ZAdd(ctx, key, zMembers...).Err()
}

// ZRange 获取有序集合指定范围成员
func (c *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合指定范围成员（带分数）（实现 CacheInterface）
func (c *RedisCache) ZRangeWithScores(key string, start, stop int64) ([]*ZMember, error) {
	ctx := context.Background()
	results, err := c.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	members := make([]*ZMember, len(results))
	for i, r := range results {
		members[i] = &ZMember{
			Score:  r.Score,
			Member: r.Member.(string),
		}
	}
	return members, nil
}

// ZRem 移除有序集合成员（实现 CacheInterface）
func (c *RedisCache) ZRem(key string, members ...string) error {
	ctx := context.Background()
	iMembers := make([]interface{}, len(members))
	for i, m := range members {
		iMembers[i] = m
	}
	return c.client.ZRem(ctx, key, iMembers...).Err()
}

// ZRemRangeByScore 按分数范围移除有序集合成员
func (c *RedisCache) ZRemRangeByScore(ctx interface{}, key, min, max string) (int64, error) {
	bgCtx := context.Background()
	return c.client.ZRemRangeByScore(bgCtx, key, min, max).Result()
}

// ZCard 获取有序集合成员数量
func (c *RedisCache) ZCard(key string) (int64, error) {
	ctx := context.Background()
	return c.client.ZCard(ctx, key).Result()
}

// SetExpire 设置键的过期时间
func (c *RedisCache) SetExpire(ctx interface{}, key string, expiration time.Duration) error {
	bgCtx := context.Background()
	return c.client.Expire(bgCtx, key, expiration).Err()
}

// FlushDB 清空当前数据库
func (c *RedisCache) FlushDB() error {
	ctx := context.Background()
	return c.client.FlushDB(ctx).Err()
}

// Ping 测试连接
func (c *RedisCache) Ping() error {
	ctx := context.Background()
	return c.client.Ping(ctx).Err()
}

// Close 关闭Redis连接
func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// HealthCheck Redis健康检查
func (c *RedisCache) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.client.Ping(ctx).Err()
}

// GetClient 获取原始Redis客户端
func (c *RedisCache) GetClient() *redis.Client {
	return c.client
}

// GetStats 获取Redis统计信息
func (c *RedisCache) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()

	info, err := c.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	poolStats := c.client.PoolStats()

	return map[string]interface{}{
		"hits":        poolStats.Hits,
		"misses":      poolStats.Misses,
		"timeouts":    poolStats.Timeouts,
		"total_conns": poolStats.TotalConns,
		"idle_conns":  poolStats.IdleConns,
		"stale_conns": poolStats.StaleConns,
		"info":        info,
	}, nil
}

package cache

import (
	"time"
)

// CacheInterface 缓存接口
// 定义了缓存操作的基本方法，支持多种缓存实现
type CacheInterface interface {
	// 基本操作
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	GetObject(key string, dest interface{}) error
	Delete(keys ...string) error
	Exists(key string) (bool, error)
	
	// 过期时间操作
	Expire(key string, expiration time.Duration) error
	TTL(key string) (time.Duration, error)
	
	// 计数器操作
	Increment(key string) (int64, error)
	IncrementBy(key string, value int64) (int64, error)
	Decrement(key string) (int64, error)
	DecrementBy(key string, value int64) (int64, error)
	
	// 高级操作
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	
	// 哈希操作
	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDelete(key string, fields ...string) error
	
	// 列表操作
	LPush(key string, values ...interface{}) error
	RPush(key string, values ...interface{}) error
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	
	// 有序集合操作
	ZAdd(key string, members ...*ZMember) error
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]*ZMember, error)
	ZRem(key string, members ...string) error
	
	// 管理操作
	FlushDB() error
	Ping() error
}

// ZMember 有序集合成员
type ZMember struct {
	Score  float64
	Member string
}
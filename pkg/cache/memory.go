package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
// 使用内存作为缓存存储，适用于开发环境或小型应用
type MemoryCache struct {
	data     map[string]*cacheItem
	mu       sync.RWMutex
	stopChan chan struct{}
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() CacheInterface {
	cache := &MemoryCache{
		data:     make(map[string]*cacheItem),
		stopChan: make(chan struct{}),
	}
	
	// 启动清理过期项的协程
	go cache.startExpirationChecker()
	
	return cache
}

// startExpirationChecker 启动过期项检查器
func (m *MemoryCache) startExpirationChecker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.cleanExpiredItems()
		case <-m.stopChan:
			return
		}
	}
}

// cleanExpiredItems 清理过期项
func (m *MemoryCache) cleanExpiredItems() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	for key, item := range m.data {
		if !item.expiration.IsZero() && item.expiration.Before(now) {
			delete(m.data, key)
		}
	}
}

// Set 设置缓存
func (m *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var expTime time.Time
	if expiration > 0 {
		expTime = time.Now().Add(expiration)
	}
	
	m.data[key] = &cacheItem{
		value:      value,
		expiration: expTime,
	}
	
	return nil
}

// Get 获取缓存
func (m *MemoryCache) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return "", errors.New("缓存不存在")
	}
	
	if !item.expiration.IsZero() && item.expiration.Before(time.Now()) {
		return "", errors.New("缓存已过期")
	}
	
	// 序列化值
	data, err := json.Marshal(item.value)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}
	
	return string(data), nil
}

// GetObject 获取缓存并反序列化为对象
func (m *MemoryCache) GetObject(key string, dest interface{}) error {
	val, err := m.Get(key)
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
func (m *MemoryCache) Delete(keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, key := range keys {
		delete(m.data, key)
	}
	
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCache) Exists(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return false, nil
	}
	
	if !item.expiration.IsZero() && item.expiration.Before(time.Now()) {
		return false, nil
	}
	
	return true, nil
}

// Expire 设置过期时间
func (m *MemoryCache) Expire(key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	if !exists {
		return errors.New("缓存不存在")
	}
	
	item.expiration = time.Now().Add(expiration)
	return nil
}

// TTL 获取剩余过期时间
func (m *MemoryCache) TTL(key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return 0, errors.New("缓存不存在")
	}
	
	if item.expiration.IsZero() {
		return -1, nil // 永不过期
	}
	
	remaining := time.Until(item.expiration)
	if remaining <= 0 {
		return 0, nil // 已过期
	}
	
	return remaining, nil
}

// Increment 自增
func (m *MemoryCache) Increment(key string) (int64, error) {
	return m.IncrementBy(key, 1)
}

// IncrementBy 自增指定值
func (m *MemoryCache) IncrementBy(key string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	var current int64
	
	if exists {
		// 尝试将值转换为int64
		switch v := item.value.(type) {
		case int64:
			current = v
		case int:
			current = int64(v)
		case float64:
			current = int64(v)
		case string:
			var num int64
			if err := json.Unmarshal([]byte(v), &num); err == nil {
				current = num
			}
		default:
			return 0, fmt.Errorf("无法将值转换为数字")
		}
	}
	
	current += value
	item.value = current
	
	// 如果项不存在，创建新项
	if !exists {
		m.data[key] = &cacheItem{
			value: current,
		}
	}
	
	return current, nil
}

// Decrement 自减
func (m *MemoryCache) Decrement(key string) (int64, error) {
	return m.IncrementBy(key, -1)
}

// DecrementBy 自减指定值
func (m *MemoryCache) DecrementBy(key string, value int64) (int64, error) {
	return m.IncrementBy(key, -value)
}

// SetNX 仅在键不存在时设置（用于分布式锁）
func (m *MemoryCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.data[key]; exists {
		return false, nil
	}
	
	var expTime time.Time
	if expiration > 0 {
		expTime = time.Now().Add(expiration)
	}
	
	m.data[key] = &cacheItem{
		value:      value,
		expiration: expTime,
	}
	
	return true, nil
}

// HSet 设置哈希字段
func (m *MemoryCache) HSet(key, field string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var item *cacheItem
	var hash map[string]interface{}
	
	if existingItem, exists := m.data[key]; exists {
		item = existingItem
		// 尝试将值转换为map
		if h, ok := item.value.(map[string]interface{}); ok {
			hash = h
		} else {
			hash = make(map[string]interface{})
		}
	} else {
		hash = make(map[string]interface{})
		item = &cacheItem{
			value: hash,
		}
		m.data[key] = item
	}
	
	hash[field] = value
	return nil
}

// HGet 获取哈希字段
func (m *MemoryCache) HGet(key, field string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return "", errors.New("哈希表不存在")
	}
	
	hash, ok := item.value.(map[string]interface{})
	if !ok {
		return "", errors.New("值不是哈希表")
	}
	
	value, exists := hash[field]
	if !exists {
		return "", errors.New("字段不存在")
	}
	
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}
	
	return string(data), nil
}

// HGetAll 获取所有哈希字段
func (m *MemoryCache) HGetAll(key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return nil, errors.New("哈希表不存在")
	}
	
	hash, ok := item.value.(map[string]interface{})
	if !ok {
		return nil, errors.New("值不是哈希表")
	}
	
	result := make(map[string]string)
	for field, value := range hash {
		// 序列化值
		data, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("序列化失败: %w", err)
		}
		result[field] = string(data)
	}
	
	return result, nil
}

// HDelete 删除哈希字段
func (m *MemoryCache) HDelete(key string, fields ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	if !exists {
		return errors.New("哈希表不存在")
	}
	
	hash, ok := item.value.(map[string]interface{})
	if !ok {
		return errors.New("值不是哈希表")
	}
	
	for _, field := range fields {
		delete(hash, field)
	}
	
	return nil
}

// LPush 列表左侧插入
func (m *MemoryCache) LPush(key string, values ...interface{}) error {
	return m.listPush(key, true, values...)
}

// RPush 列表右侧插入
func (m *MemoryCache) RPush(key string, values ...interface{}) error {
	return m.listPush(key, false, values...)
}

// listPush 列表插入实现
func (m *MemoryCache) listPush(key string, left bool, values ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var item *cacheItem
	var list []interface{}
	
	if existingItem, exists := m.data[key]; exists {
		item = existingItem
		// 尝试将值转换为slice
		if l, ok := item.value.([]interface{}); ok {
			list = l
		} else {
			list = make([]interface{}, 0)
		}
	} else {
		list = make([]interface{}, 0)
		item = &cacheItem{
			value: list,
		}
		m.data[key] = item
	}
	
	if left {
		// 左侧插入
		list = append(values, list...)
	} else {
		// 右侧插入
		list = append(list, values...)
	}
	
	item.value = list
	return nil
}

// LPop 列表左侧弹出
func (m *MemoryCache) LPop(key string) (string, error) {
	return m.listPop(key, true)
}

// RPop 列表右侧弹出
func (m *MemoryCache) RPop(key string) (string, error) {
	return m.listPop(key, false)
}

// listPop 列表弹出实现
func (m *MemoryCache) listPop(key string, left bool) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	if !exists {
		return "", errors.New("列表不存在")
	}
	
	list, ok := item.value.([]interface{})
	if !ok {
		return "", errors.New("值不是列表")
	}
	
	if len(list) == 0 {
		return "", errors.New("列表为空")
	}
	
	var value interface{}
	if left {
		// 左侧弹出
		value = list[0]
		list = list[1:]
	} else {
		// 右侧弹出
		value = list[len(list)-1]
		list = list[:len(list)-1]
	}
	
	item.value = list
	
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}
	
	return string(data), nil
}

// LRange 获取列表范围
func (m *MemoryCache) LRange(key string, start, stop int64) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return nil, errors.New("列表不存在")
	}
	
	list, ok := item.value.([]interface{})
	if !ok {
		return nil, errors.New("值不是列表")
	}
	
	// 处理负数索引
	length := int64(len(list))
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	
	// 边界检查
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	
	if start > stop || start >= length {
		return []string{}, nil
	}
	
	result := make([]string, 0, stop-start+1)
	for i := start; i <= stop; i++ {
		// 序列化值
		data, err := json.Marshal(list[int(i)])
		if err != nil {
			return nil, fmt.Errorf("序列化失败: %w", err)
		}
		result = append(result, string(data))
	}
	
	return result, nil
}

// ZAdd 添加有序集合成员
func (m *MemoryCache) ZAdd(key string, members ...*ZMember) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var item *cacheItem
	var zset []*ZMember
	
	if existingItem, exists := m.data[key]; exists {
		item = existingItem
		// 尝试将值转换为有序集合
		if z, ok := item.value.([]*ZMember); ok {
			zset = z
		} else {
			zset = make([]*ZMember, 0)
		}
	} else {
		zset = make([]*ZMember, 0)
		item = &cacheItem{
			value: zset,
		}
		m.data[key] = item
	}
	
	for _, newMember := range members {
		// 检查成员是否已存在
		found := false
		for i, member := range zset {
			if member.Member == newMember.Member {
				// 更新分数
				zset[i].Score = newMember.Score
				found = true
				break
			}
		}
		
		if !found {
			// 添加新成员
			zset = append(zset, newMember)
		}
	}
	
	item.value = zset
	return nil
}

// ZRange 获取有序集合范围
func (m *MemoryCache) ZRange(key string, start, stop int64) ([]string, error) {
	members, err := m.ZRangeWithScores(key, start, stop)
	if err != nil {
		return nil, err
	}
	
	result := make([]string, len(members))
	for i, member := range members {
		result[i] = member.Member
	}
	
	return result, nil
}

// ZRangeWithScores 获取有序集合范围（包含分数）
func (m *MemoryCache) ZRangeWithScores(key string, start, stop int64) ([]*ZMember, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return nil, errors.New("有序集合不存在")
	}
	
	zset, ok := item.value.([]*ZMember)
	if !ok {
		return nil, errors.New("值不是有序集合")
	}
	
	// 按分数排序
	sort.Slice(zset, func(i, j int) bool {
		return zset[i].Score < zset[j].Score
	})
	
	// 处理负数索引
	length := int64(len(zset))
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	
	// 边界检查
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	
	if start > stop || start >= length {
		return []*ZMember{}, nil
	}
	
	result := make([]*ZMember, 0, stop-start+1)
	for i := start; i <= stop; i++ {
		result = append(result, zset[int(i)])
	}
	
	return result, nil
}

// ZRem 删除有序集合成员
func (m *MemoryCache) ZRem(key string, members ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	if !exists {
		return errors.New("有序集合不存在")
	}
	
	zset, ok := item.value.([]*ZMember)
	if !ok {
		return errors.New("值不是有序集合")
	}
	
	for _, memberToRemove := range members {
		for i, member := range zset {
			if member.Member == memberToRemove {
				// 删除成员
				zset = append(zset[:i], zset[i+1:]...)
				break
			}
		}
	}
	
	item.value = zset
	return nil
}

// FlushDB 清空当前数据库
func (m *MemoryCache) FlushDB() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data = make(map[string]*cacheItem)
	return nil
}

// Ping 测试连接
func (m *MemoryCache) Ping() error {
	return nil
}

// Close 关闭缓存
func (m *MemoryCache) Close() error {
	close(m.stopChan)
	return nil
}

// SetExpire 设置过期时间（兼容接口）
func (m *MemoryCache) SetExpire(ctx interface{}, key string, expiration time.Duration) error {
	return m.Expire(key, expiration)
}

// ZCard 获取有序集合大小
func (m *MemoryCache) ZCard(key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[key]
	if !exists {
		return 0, nil
	}
	
	zset, ok := item.value.([]*ZMember)
	if !ok {
		return 0, errors.New("值不是有序集合")
	}
	
	return int64(len(zset)), nil
}

// ZRemRangeByScore 按分数范围删除有序集合成员
func (m *MemoryCache) ZRemRangeByScore(ctx interface{}, key, min, max string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, exists := m.data[key]
	if !exists {
		return 0, errors.New("有序集合不存在")
	}
	
	zset, ok := item.value.([]*ZMember)
	if !ok {
		return 0, errors.New("值不是有序集合")
	}
	
	// 解析分数范围
	minScore, err := parseFloat(min)
	if err != nil {
		return 0, fmt.Errorf("无效的最小分数: %w", err)
	}
	
	maxScore, err := parseFloat(max)
	if err != nil {
		return 0, fmt.Errorf("无效的最大分数: %w", err)
	}
	
	// 过滤出不在范围内的成员
	newZset := make([]*ZMember, 0)
	removed := 0
	
	for _, member := range zset {
		if member.Score < minScore || member.Score > maxScore {
			newZset = append(newZset, member)
		} else {
			removed++
		}
	}
	
	item.value = newZset
	return int64(removed), nil
}

// parseFloat 解析浮点数，支持"-inf"和"+inf"
func parseFloat(s string) (float64, error) {
	switch s {
	case "-inf":
		return -1.7976931348623157e+308, nil // 最小浮点数
	case "+inf", "inf":
		return 1.7976931348623157e+308, nil // 最大浮点数
	default:
		// 尝试解析为浮点数
		var f float64
		_, err := fmt.Sscanf(s, "%f", &f)
		if err != nil {
			return 0, err
		}
		return f, nil
	}
}
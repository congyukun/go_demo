// Package captcha 提供验证码生成和验证功能
package captcha

import (
	"fmt"
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
)

// CaptchaService 验证码服务接口
type CaptchaService interface {
	// Generate 生成验证码，返回验证码ID和Base64编码的图片
	Generate() (id string, b64s string, err error)
	// Verify 验证验证码，验证后自动删除
	Verify(id, answer string) bool
	// VerifyWithoutClear 验证验证码但不删除（用于调试）
	VerifyWithoutClear(id, answer string) bool
}

// Config 验证码配置
type Config struct {
	Height     int           // 图片高度
	Width      int           // 图片宽度
	Length     int           // 验证码长度
	MaxSkew    float64       // 最大倾斜度
	DotCount   int           // 干扰点数量
	ExpireTime time.Duration // 过期时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Height:     80,
		Width:      240,
		Length:     4,
		MaxSkew:    0.7,
		DotCount:   80,
		ExpireTime: 5 * time.Minute,
	}
}

// captchaService 验证码服务实现
type captchaService struct {
	store  base64Captcha.Store
	driver *base64Captcha.DriverString
	config Config
}

// memoryStore 内存存储实现（带过期时间）
type memoryStore struct {
	sync.RWMutex
	data       map[string]storeItem
	expireTime time.Duration
}

type storeItem struct {
	value     string
	expiredAt time.Time
}

// newMemoryStore 创建内存存储
func newMemoryStore(expireTime time.Duration) *memoryStore {
	store := &memoryStore{
		data:       make(map[string]storeItem),
		expireTime: expireTime,
	}
	// 启动清理协程
	go store.cleanupLoop()
	return store
}

// Set 存储验证码
func (s *memoryStore) Set(id string, value string) error {
	s.Lock()
	defer s.Unlock()
	s.data[id] = storeItem{
		value:     value,
		expiredAt: time.Now().Add(s.expireTime),
	}
	return nil
}

// Get 获取验证码（验证后删除）
func (s *memoryStore) Get(id string, clear bool) string {
	s.Lock()
	defer s.Unlock()

	item, ok := s.data[id]
	if !ok {
		return ""
	}

	// 检查是否过期
	if time.Now().After(item.expiredAt) {
		delete(s.data, id)
		return ""
	}

	if clear {
		delete(s.data, id)
	}
	return item.value
}

// Verify 验证验证码
func (s *memoryStore) Verify(id, answer string, clear bool) bool {
	value := s.Get(id, clear)
	return value != "" && value == answer
}

// cleanupLoop 定期清理过期验证码
func (s *memoryStore) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

// cleanup 清理过期验证码
func (s *memoryStore) cleanup() {
	s.Lock()
	defer s.Unlock()

	now := time.Now()
	for id, item := range s.data {
		if now.After(item.expiredAt) {
			delete(s.data, id)
		}
	}
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(cfg Config) CaptchaService {
	store := newMemoryStore(cfg.ExpireTime)

	driver := &base64Captcha.DriverString{
		Height:          cfg.Height,
		Width:           cfg.Width,
		NoiseCount:      cfg.DotCount,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          cfg.Length,
		Source:          "0123456789abcdefghijklmnopqrstuvwxyz",
		Fonts:           []string{"wqy-microhei.ttc"},
	}

	return &captchaService{
		store:  store,
		driver: driver,
		config: cfg,
	}
}

// NewDefaultCaptchaService 使用默认配置创建验证码服务
func NewDefaultCaptchaService() CaptchaService {
	return NewCaptchaService(DefaultConfig())
}

// Generate 生成验证码
func (s *captchaService) Generate() (id string, b64s string, err error) {
	captcha := base64Captcha.NewCaptcha(s.driver, s.store)
	id, b64s, _, err = captcha.Generate()
	if err != nil {
		return "", "", fmt.Errorf("生成验证码失败: %w", err)
	}
	return id, b64s, nil
}

// Verify 验证验证码
func (s *captchaService) Verify(id, answer string) bool {
	return s.store.Verify(id, answer, true)
}

// VerifyWithoutClear 验证验证码但不删除
func (s *captchaService) VerifyWithoutClear(id, answer string) bool {
	return s.store.Verify(id, answer, false)
}

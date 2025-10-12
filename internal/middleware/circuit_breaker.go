package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"go_demo/internal/utils"
	"go_demo/pkg/logger"
)

// CircuitState 熔断器状态
type CircuitState int

const (
	StateClosed CircuitState = iota // 关闭状态，正常工作
	StateOpen                       // 打开状态，熔断中
	StateHalfOpen                   // 半开状态，尝试恢复
)

// String 返回状态的字符串表示
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	// 熔断器名称
	Name string
	// 最大请求数，达到此数量后开始计算错误率
	MaxRequests uint32
	// 错误率阈值，超过此阈值则熔断
	ErrorThreshold float64
	// 熔断器打开后的超时时间
	Timeout time.Duration
	// 半开状态下的最大请求数
	MaxHalfOpenRequests uint32
	// 是否启用熔断器
	Enabled bool
	// 熔断时的回调函数
	OnCircuitOpen func(*gin.Context)
	// 熔断器状态变更回调函数
	OnStateChange func(name string, from, to CircuitState)
}

// DefaultCircuitBreakerConfig 返回默认的熔断器配置
func DefaultCircuitBreakerConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:                name,
		MaxRequests:         100,
		ErrorThreshold:      0.5, // 50%错误率
		Timeout:             time.Minute,
		MaxHalfOpenRequests: 10,
		Enabled:             true,
		OnCircuitOpen: func(c *gin.Context) {
			utils.ResponseError(c, http.StatusServiceUnavailable, "服务暂时不可用，请稍后再试")
		},
		OnStateChange: func(name string, from, to CircuitState) {
			// 只在日志初始化后才记录状态变更
			if logger.Logger != nil {
			logger.Info("熔断器状态变更",
				logger.String("name", name),
				logger.String("from", from.String()),
				logger.String("to", to.String()),
			)
		}
		},
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	config     CircuitBreakerConfig
	state      CircuitState
	generation uint64 // 代数，每次状态变更时递增
	counts     Counts
	expiry     time.Time
	mutex      sync.Mutex
}

// Counts 计数器
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// NewCircuitBreaker 创建新的熔断器
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute 执行函数，如果熔断器打开则返回错误
func (cb *CircuitBreaker) Execute(req func() error) error {
	if !cb.config.Enabled {
		return req()
	}

	generation, err := cb.beforeRequest()
	if err != nil {
		return err
	}

	err = req()
	cb.afterRequest(generation, err)
	return err
}

// ExecuteGin 在Gin中间件中执行熔断器
func (cb *CircuitBreaker) ExecuteGin(c *gin.Context, handler gin.HandlerFunc) {
	if !cb.config.Enabled {
		handler(c)
		return
	}

	generation, err := cb.beforeRequest()
	if err != nil {
		// 只在日志初始化后才记录警告
		if logger.Logger != nil {
			logger.Warn("熔断器阻止请求",
				logger.String("name", cb.config.Name),
				logger.String("path", c.Request.URL.Path),
				logger.String("state", cb.state.String()),
			)
		}
		cb.config.OnCircuitOpen(c)
		c.Abort()
		return
	}

	// 执行处理器
	handler(c)

	// 检查是否有错误
	if c.Writer.Status() >= 500 {
		cb.afterRequest(generation, errors.New("server error"))
	} else {
		cb.afterRequest(generation, nil)
	}
}

// beforeRequest 请求前的处理
func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, errors.New("circuit breaker is open")
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.config.MaxHalfOpenRequests {
		return generation, errors.New("circuit breaker is half open and max requests reached")
	}

	cb.counts.Requests++
	return generation, nil
}

// afterRequest 请求后的处理
func (cb *CircuitBreaker) afterRequest(before uint64, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if err != nil {
		cb.onFailure(state, now)
	} else {
		cb.onSuccess(state, now)
	}
}

// onSuccess 成功时的处理
func (cb *CircuitBreaker) onSuccess(state CircuitState, now time.Time) {
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	if state == StateHalfOpen && cb.counts.ConsecutiveSuccesses >= cb.config.MaxHalfOpenRequests {
		cb.setState(StateClosed, now)
	}
}

// onFailure 失败时的处理
func (cb *CircuitBreaker) onFailure(state CircuitState, now time.Time) {
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0

	if state == StateClosed && cb.counts.Requests >= cb.config.MaxRequests &&
		float64(cb.counts.TotalFailures)/float64(cb.counts.Requests) >= cb.config.ErrorThreshold {
		cb.setState(StateOpen, now)
	} else if state == StateHalfOpen {
		cb.setState(StateOpen, now)
	}
}

// currentState 获取当前状态
func (cb *CircuitBreaker) currentState(now time.Time) (CircuitState, uint64) {
	switch cb.state {
	case StateClosed:
		// 关闭状态不需要检查过期时间
		return cb.state, cb.generation
	case StateOpen:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

// setState 设置状态
func (cb *CircuitBreaker) setState(state CircuitState, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.config.Name, prev, state)
	}
}

// toNewGeneration 切换到新一代
func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts = Counts{}

	var zero time.Time
	switch cb.state {
	case StateClosed:
		// 关闭状态不需要过期时间
		cb.expiry = zero
	case StateOpen:
		// 打开状态需要设置超时时间，超时后转为半开状态
		cb.expiry = now.Add(cb.config.Timeout)
	default: // StateHalfOpen
		// 半开状态不需要过期时间
		cb.expiry = zero
	}
}

// CircuitBreakerMiddleware 创建熔断器中间件
func CircuitBreakerMiddleware(config CircuitBreakerConfig) gin.HandlerFunc {
	cb := NewCircuitBreaker(config)
	
	return func(c *gin.Context) {
		cb.ExecuteGin(c, func(c *gin.Context) {
			c.Next()
		})
	}
}

// APICircuitBreaker API级别的熔断器中间件
func APICircuitBreaker(name string) gin.HandlerFunc {
	config := DefaultCircuitBreakerConfig(name)
	return CircuitBreakerMiddleware(config)
}

// 全局熔断器管理器
var (
	circuitBreakers = make(map[string]*CircuitBreaker)
	cbMutex         sync.RWMutex
)

// GetCircuitBreaker 获取或创建熔断器
func GetCircuitBreaker(name string, config CircuitBreakerConfig) *CircuitBreaker {
	cbMutex.RLock()
	cb, exists := circuitBreakers[name]
	cbMutex.RUnlock()
	
	if exists {
		return cb
	}
	
	cbMutex.Lock()
	defer cbMutex.Unlock()
	
	// 再次检查，防止并发创建
	if cb, exists := circuitBreakers[name]; exists {
		return cb
	}
	
	cb = NewCircuitBreaker(config)
	circuitBreakers[name] = cb
	return cb
}

// ServiceCircuitBreaker 服务级别的熔断器中间件
func ServiceCircuitBreaker(serviceName string) gin.HandlerFunc {
	config := DefaultCircuitBreakerConfig(serviceName)
	cb := GetCircuitBreaker(serviceName, config)
	
	return func(c *gin.Context) {
		cb.ExecuteGin(c, func(c *gin.Context) {
			c.Next()
		})
	}
}

// CircuitBreakerFactory 熔断器工厂
type CircuitBreakerFactory struct {
	config CircuitBreakerConfig
	cb     *CircuitBreaker
}

// NewCircuitBreakerFactory 创建熔断器工厂
func NewCircuitBreakerFactory(config CircuitBreakerConfig) *CircuitBreakerFactory {
	cb := NewCircuitBreaker(config)
	
	return &CircuitBreakerFactory{
		config: config,
		cb:     cb,
	}
}

// Middleware 返回熔断器中间件
func (f *CircuitBreakerFactory) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		f.cb.ExecuteGin(c, func(c *gin.Context) {
			c.Next()
		})
	}
}
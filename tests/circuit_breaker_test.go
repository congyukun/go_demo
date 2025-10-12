package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go_demo/internal/middleware"
)

// TestCircuitBreakerBasic 测试熔断器基本功能
func TestCircuitBreakerBasic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := middleware.DefaultCircuitBreakerConfig("test")
	config.MaxRequests = 5
	config.ErrorThreshold = 0.5
	config.MaxHalfOpenRequests = 3
	config.Timeout = 1 * time.Second

	cb := middleware.NewCircuitBreaker(config)

	// 测试1: 正常请求应该通过
	t.Run("NormalRequests", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			err := cb.Execute(func() error {
				return nil
			})
			assert.NoError(t, err)
		}
	})

	// 测试2: 失败请求达到阈值后应该熔断
	t.Run("CircuitBreakerShouldOpen", func(t *testing.T) {
		// 发送失败请求，达到错误率阈值
		for i := 0; i < 3; i++ {
			cb.Execute(func() error {
				return errors.New("test error")
			})
		}

		// 此时错误率应该达到阈值，熔断器应该打开
		// 注意：此时总请求数为 3(成功) + 3(失败) = 6，失败率为 3/6 = 50%，达到阈值
		err := cb.Execute(func() error {
			return nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit breaker is open")
	})
}

// TestCircuitBreakerGinMiddleware 测试Gin中间件中的熔断器
func TestCircuitBreakerGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := middleware.DefaultCircuitBreakerConfig("test-api")
	config.MaxRequests = 3
	config.ErrorThreshold = 0.5
	config.MaxHalfOpenRequests = 2
	config.Timeout = 1 * time.Second

	router := gin.New()

	// 添加熔断器中间件
	router.Use(middleware.CircuitBreakerMiddleware(config))

	// 模拟成功处理函数
	successHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}

	// 模拟失败处理函数
	failHandler := func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
	}

	router.GET("/success", successHandler)
	router.GET("/fail", failHandler)

	// 测试1: 正常请求
	t.Run("SuccessfulRequests", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/success", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	// 测试2: 失败请求导致熔断
	t.Run("CircuitBreakerOpens", func(t *testing.T) {
		// 发送失败请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/fail", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 再发送一个失败请求，达到错误率阈值
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/fail", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 此时熔断器应该打开，成功请求也会被拒绝
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/success", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

// TestCircuitBreakerConcurrent 测试并发场景下的熔断器
func TestCircuitBreakerConcurrent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := middleware.DefaultCircuitBreakerConfig("test-concurrent")
	config.MaxRequests = 10
	config.ErrorThreshold = 0.5
	config.MaxHalfOpenRequests = 5
	config.Timeout = 2 * time.Second

	cb := middleware.NewCircuitBreaker(config)

	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	// 并发发送请求
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			err := cb.Execute(func() error {
				if index%3 == 0 {
					return errors.New("simulated error")
				}
				return nil
			})

			mu.Lock()
			if err != nil {
				failCount++
			} else {
				successCount++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	t.Logf("并发测试结果 - 成功: %d, 失败: %d", successCount, failCount)

	// 验证结果
	assert.Greater(t, successCount, 0, "应该有一些成功的请求")
	assert.Greater(t, failCount, 0, "应该有一些失败的请求")
}

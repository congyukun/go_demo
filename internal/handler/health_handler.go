package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_demo/pkg/cache"
	"go_demo/pkg/database"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB, cache *cache.RedisCache) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services"`
	Version   string                 `json:"version,omitempty"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Health 基础健康检查
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags 健康检查
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// HealthCheck 详细健康检查
// @Summary 详细健康检查
// @Description 检查服务及其依赖的健康状态
// @Tags 健康检查
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health/check [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]ServiceInfo),
		Version:   "1.0.0",
	}

	allHealthy := true

	// 检查数据库
	dbInfo := h.checkDatabase()
	response.Services["database"] = dbInfo
	if dbInfo.Status != "healthy" {
		allHealthy = false
	}

	// 检查Redis
	redisInfo := h.checkRedis()
	response.Services["redis"] = redisInfo
	if redisInfo.Status != "healthy" {
		allHealthy = false
	}

	// 设置整体状态
	if !allHealthy {
		response.Status = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// checkDatabase 检查数据库健康状态
func (h *HealthHandler) checkDatabase() ServiceInfo {
	if h.db == nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: "数据库连接未初始化",
		}
	}

	err := database.HealthCheck(h.db)
	if err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	// 获取数据库统计信息
	stats, err := database.GetStats(h.db)
	if err != nil {
		return ServiceInfo{
			Status:  "healthy",
			Message: "数据库连接正常，但无法获取统计信息",
		}
	}

	return ServiceInfo{
		Status:  "healthy",
		Message: "数据库连接正常",
		Details: stats,
	}
}

// checkRedis 检查Redis健康状态
func (h *HealthHandler) checkRedis() ServiceInfo {
	if h.cache == nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: "Redis连接未初始化",
		}
	}

	err := h.cache.HealthCheck()
	if err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	// 获取Redis统计信息
	stats, err := h.cache.GetStats()
	if err != nil {
		return ServiceInfo{
			Status:  "healthy",
			Message: "Redis连接正常，但无法获取统计信息",
		}
	}

	return ServiceInfo{
		Status:  "healthy",
		Message: "Redis连接正常",
		Details: stats,
	}
}

// Readiness 就绪检查（用于Kubernetes）
// @Summary 就绪检查
// @Description 检查服务是否准备好接收流量
// @Tags 健康检查
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// 检查关键依赖
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"message": "数据库未初始化",
		})
		return
	}

	if err := database.HealthCheck(h.db); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"message": "数据库连接失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Liveness 存活检查（用于Kubernetes）
// @Summary 存活检查
// @Description 检查服务是否存活
// @Tags 健康检查
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

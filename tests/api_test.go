package tests

import (
	"bytes"
	"encoding/json"
	"go_demo/internal/handler"
	"go_demo/internal/models"
	"go_demo/internal/repository"
	"go_demo/internal/router"
	"go_demo/internal/service"
	"go_demo/pkg/logger"
	"go_demo/pkg/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupAPITestRouter 设置API测试路由
func setupAPITestRouter(t *testing.T) *gin.Engine {
	// 初始化验证器
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	// 初始化日志
	logConfig := logger.LogConfig{
		Level:      "error", // 测试时使用error级别减少日志输出
		Format:     "console",
		OutputPath: "/tmp/api_test.log",
	}
	if err := logger.Init(logConfig); err != nil {
		t.Fatalf("日志初始化失败: %v", err)
	}

	// 使用现有的测试数据库设置
	db := setupTestDB(t)

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)

	// 初始化服务层
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)

	// 设置路由
	r := router.NewRouter(authHandler, userHandler)
	engine := r.Setup()

	return engine
}

func TestHealthCheck(t *testing.T) {
	router := setupAPITestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 实际 %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("期望状态 ok, 实际 %v", response["status"])
	}
}

func TestUserRegistration(t *testing.T) {
	router := setupAPITestRouter(t)

	// 测试用户注册
	registerReq := models.RegisterRequest{
		Username: "testuser",
		Password: "123456",
		Name:     "测试用户",
		Email:    "test@example.com",
		Mobile:   "13812345678",
	}

	jsonData, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 实际 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["message"] != "注册成功" {
		t.Errorf("期望消息 '注册成功', 实际 %v", response["message"])
	}
}

func TestUserLogin(t *testing.T) {
	router := setupAPITestRouter(t)

	// 先注册用户
	registerReq := models.RegisterRequest{
		Username: "loginuser",
		Password: "123456",
		Name:     "登录测试用户",
		Email:    "login@example.com",
		Mobile:   "15987654321",
	}

	jsonData, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("注册失败: %d, %s", w.Code, w.Body.String())
	}

	// 测试登录
	loginReq := models.LoginRequest{
		Username: "loginuser",
		Password: "123456",
	}

	jsonData, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 实际 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["message"] != "登录成功" {
		t.Errorf("期望消息 '登录成功', 实际 %v", response["message"])
	}

	// 检查返回的数据中是否包含token
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应数据格式错误")
	}

	if data["access_token"] == nil || data["access_token"] == "" {
		t.Error("access_token 不应该为空")
	}

	if data["refresh_token"] == nil || data["refresh_token"] == "" {
		t.Error("refresh_token 不应该为空")
	}
}

func TestInvalidLogin(t *testing.T) {
	router := setupAPITestRouter(t)

	// 测试无效登录
	loginReq := models.LoginRequest{
		Username: "nonexistent",
		Password: "wrongpassword",
	}

	jsonData, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期望状态码 401, 实际 %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok {
		t.Fatal("响应消息格式错误")
	}

	if !strings.Contains(message, "用户名或密码错误") {
		t.Errorf("期望错误消息包含 '用户名或密码错误', 实际: %s", message)
	}
}

func TestValidationErrors(t *testing.T) {
	router := setupAPITestRouter(t)

	// 测试验证错误 - 用户名太短
	registerReq := models.RegisterRequest{
		Username: "ab", // 太短
		Password: "123456",
		Name:     "测试用户",
		Email:    "test@example.com",
	}

	jsonData, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("期望状态码 400, 实际 %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok {
		t.Fatal("响应消息格式错误")
	}

	if !strings.Contains(message, "用户名") {
		t.Errorf("期望错误消息包含 '用户名', 实际: %s", message)
	}
}

func TestDuplicateRegistration(t *testing.T) {
	router := setupAPITestRouter(t)

	// 注册第一个用户
	registerReq := models.RegisterRequest{
		Username: "duplicateuser",
		Password: "123456",
		Name:     "重复测试用户",
		Email:    "duplicate@example.com",
		Mobile:   "18666666666",
	}

	jsonData, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("第一次注册失败: %d, %s", w.Code, w.Body.String())
	}

	// 尝试注册相同用户名的用户
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 409 { // Conflict
		t.Errorf("期望状态码 409, 实际 %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok {
		t.Fatal("响应消息格式错误")
	}

	if !strings.Contains(message, "已存在") {
		t.Errorf("期望错误消息包含 '已存在', 实际: %s", message)
	}
}

// getAuthToken 获取认证token的辅助函数
func getAuthToken(t *testing.T, router *gin.Engine) string {
	// 先注册用户
	registerReq := models.RegisterRequest{
		Username: "authuser",
		Password: "123456",
		Name:     "认证测试用户",
		Email:    "auth@example.com",
		Mobile:   "10.27.0",
	}

	jsonData, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 登录获取token
	loginReq := models.LoginRequest{
		Username: "authuser",
		Password: "123456",
	}

	jsonData, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	return data["access_token"].(string)
}

func TestAuthenticatedEndpoints(t *testing.T) {
	router := setupAPITestRouter(t)

	// 获取认证token
	token := getAuthToken(t, router)

	// 测试获取用户列表（需要认证）
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 实际 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["message"] != "获取成功" {
		t.Errorf("期望消息 '获取成功', 实际 %v", response["message"])
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	router := setupAPITestRouter(t)

	// 测试未认证访问受保护的端点
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期望状态码 401, 实际 %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok {
		t.Fatal("响应消息格式错误")
	}

	if !strings.Contains(message, "未认证") {
		t.Errorf("期望错误消息包含 '未认证', 实际: %s", message)
	}
}

func TestGetProfile(t *testing.T) {
	router := setupAPITestRouter(t)

	// 获取认证token
	token := getAuthToken(t, router)

	// 测试获取用户资料
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期望状态码 200, 实际 %d, 响应: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["message"] != "获取成功" {
		t.Errorf("期望消息 '获取成功', 实际 %v", response["message"])
	}
}

package tests

import (
	"go_demo/internal/models"
	"go_demo/internal/repository"
	"go_demo/internal/service"
	"go_demo/pkg/logger"
	"go_demo/pkg/validator"
	"testing"
)

func TestAuthServiceLayer(t *testing.T) {
	// 初始化依赖
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	logConfig := logger.LogConfig{
		Level:      "error",
		Format:     "console",
		OutputPath: "/tmp/service_test.log",
	}
	if err := logger.Init(logConfig); err != nil {
		t.Fatalf("日志初始化失败: %v", err)
	}

	// 设置测试数据库
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)

	t.Run("用户注册", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "servicetest",
			Password: "123456",
			Name:     "服务测试用户",
			Email:    "service@example.com",
			Mobile:   "13812345678",
		}

		user, err := authService.Register(req)
		if err != nil {
			t.Fatalf("注册失败: %v", err)
		}

		if user == nil {
			t.Fatal("返回的用户不应该为空")
		}

		if user.Username != req.Username {
			t.Errorf("期望用户名 %s, 实际 %s", req.Username, user.Username)
		}

		if user.Email != req.Email {
			t.Errorf("期望邮箱 %s, 实际 %s", req.Email, user.Email)
		}
	})

	t.Run("重复注册", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "duplicate",
			Password: "123456",
			Name:     "重复用户",
			Email:    "duplicate@example.com",
			Mobile:   "15987654321",
		}

		// 第一次注册
		_, err := authService.Register(req)
		if err != nil {
			t.Fatalf("第一次注册失败: %v", err)
		}

		// 第二次注册应该失败
		_, err = authService.Register(req)
		if err == nil {
			t.Error("重复注册应该失败")
		}

		if err.Error() != "用户名已存在" {
			t.Errorf("期望错误 '用户名已存在', 实际 %s", err.Error())
		}
	})

	t.Run("用户登录", func(t *testing.T) {
		// 先注册用户
		registerReq := models.RegisterRequest{
			Username: "logintest",
			Password: "123456",
			Name:     "登录测试用户",
			Email:    "logintest@example.com",
			Mobile:   "18666666666",
		}

		_, err := authService.Register(registerReq)
		if err != nil {
			t.Fatalf("注册失败: %v", err)
		}

		// 测试登录
		loginReq := models.LoginRequest{
			Username: "logintest",
			Password: "123456",
		}

		response, err := authService.Login(loginReq)
		if err != nil {
			t.Fatalf("登录失败: %v", err)
		}

		if response == nil {
			t.Fatal("登录响应不应该为空")
		}

		if response.Token == "" {
			t.Error("访问令牌不应该为空")
		}

		if response.User.ID == 0 {
			t.Error("用户信息不应该为空")
		}
	})

	t.Run("错误登录", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Username: "nonexistent",
			Password: "wrongpassword",
		}

		_, err := authService.Login(loginReq)
		if err == nil {
			t.Error("错误的登录信息应该失败")
		}

		if err.Error() != "用户名或密码错误" {
			t.Errorf("期望错误 '用户名或密码错误', 实际 %s", err.Error())
		}
	})
}

func TestUserService(t *testing.T) {
	// 初始化依赖
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	logConfig := logger.LogConfig{
		Level:      "error",
		Format:     "console",
		OutputPath: "/tmp/service_test.log",
	}
	if err := logger.Init(logConfig); err != nil {
		t.Fatalf("日志初始化失败: %v", err)
	}

	// 设置测试数据库
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 创建测试用户
	testUser := &models.User{
		Username: "userservicetest",
		Email:    "userservice@example.com",
		Name:     "用户服务测试",
		Password: "hashed_Password",
		Mobile:   "10.27.0",
		Status:   1,
	}
	db.Create(testUser)

	t.Run("获取用户列表", func(t *testing.T) {
		users, total, err := userService.GetUsers(1, 10)
		if err != nil {
			t.Fatalf("获取用户列表失败: %v", err)
		}

		if total == 0 {
			t.Error("用户总数不应该为0")
		}

		if len(users) == 0 {
			t.Error("用户列表不应该为空")
		}
	})

	t.Run("根据ID获取用户", func(t *testing.T) {
		user, err := userService.GetUserByID(int(testUser.ID))
		if err != nil {
			t.Fatalf("获取用户失败: %v", err)
		}

		if user == nil {
			t.Fatal("用户不应该为空")
		}

		if user.Username != testUser.Username {
			t.Errorf("期望用户名 %s, 实际 %s", testUser.Username, user.Username)
		}
	})

	t.Run("获取不存在的用户", func(t *testing.T) {
		_, err := userService.GetUserByID(99999)
		if err == nil {
			t.Error("获取不存在的用户应该失败")
		}

		if err.Error() != "用户不存在" {
			t.Errorf("期望错误 '用户不存在', 实际 %s", err.Error())
		}
	})

	t.Run("创建用户", func(t *testing.T) {
		req := models.UserCreateRequest{
			Username: "newuser",
			Password: "123456",
			Email:    "newuser@example.com",
		}

		user, err := userService.CreateUser(req)
		if err != nil {
			t.Fatalf("创建用户失败: %v", err)
		}

		if user == nil {
			t.Fatal("创建的用户不应该为空")
		}

		if user.Username != req.Username {
			t.Errorf("期望用户名 %s, 实际 %s", req.Username, user.Username)
		}
	})

	t.Run("更新用户", func(t *testing.T) {
		newEmail := "updated@example.com"
		newName := "更新后的名称"
		status := 0

		req := models.UpdateUserRequest{
			Email:  newEmail,
			Name:   newName,
			Status: &status,
		}

		user, err := userService.UpdateUser(int(testUser.ID), req)
		if err != nil {
			t.Fatalf("更新用户失败: %v", err)
		}

		if user == nil {
			t.Fatal("更新后的用户不应该为空")
		}

		if user.Email != newEmail {
			t.Errorf("期望邮箱 %s, 实际 %s", newEmail, user.Email)
		}

		if user.Status != status {
			t.Errorf("期望状态 %d, 实际 %d", status, user.Status)
		}
	})

	t.Run("删除用户", func(t *testing.T) {
		// 创建一个用于删除的用户
		deleteUser := &models.User{
			Username: "deletetest",
			Email:    "delete@example.com",
			Name:     "删除测试用户",
			Password: "hashed_Password",
			Mobile:   "Password",
			Status:   1,
		}
		db.Create(deleteUser)

		err := userService.DeleteUser(int(deleteUser.ID))
		if err != nil {
			t.Fatalf("删除用户失败: %v", err)
		}

		// 验证用户已被软删除
		_, err = userService.GetUserByID(int(deleteUser.ID))
		if err == nil {
			t.Error("删除的用户不应该被找到")
		}
	})

	t.Run("获取用户统计", func(t *testing.T) {
		stats, err := userService.GetUserStats()
		if err != nil {
			t.Fatalf("获取用户统计失败: %v", err)
		}

		if stats == nil {
			t.Fatal("统计信息不应该为空")
		}

		// 检查统计信息的基本字段
		if _, ok := stats["total"]; !ok {
			t.Error("统计信息应该包含total字段")
		}

		if _, ok := stats["active"]; !ok {
			t.Error("统计信息应该包含active字段")
		}

		if _, ok := stats["inactive"]; !ok {
			t.Error("统计信息应该包含inactive字段")
		}
	})
}

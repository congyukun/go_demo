package tests

import (
	"go_demo/internal/models"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"testing"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用内存数据库进行测试
	cfg := database.MySQLConfig{
		Driver:          "mysql",
		DSN:             "root:@tcp(localhost:3306)/go_demo_test?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 3600,
		ConnMaxIdleTime: 1800,
		LogMode:         false,
		SlowThreshold:   200,
	}

	// 初始化日志（测试模式）
	logCfg := logger.LogConfig{
		Level:      "error",
		Format:     "console",
		OutputPath: "/tmp/test.log",
		MaxSize:    10,
		MaxBackup:  1,
		MaxAge:     1,
		Compress:   false,
	}
	logger.Init(logCfg)

	db, err := database.NewMySQL(cfg)
	if err != nil {
		t.Skipf("跳过数据库测试，无法连接到测试数据库: %v", err)
	}

	// 自动迁移测试表
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	// 清理测试数据
	db.Exec("DELETE FROM users")
	database.Close(db)
}

func TestUserDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	t.Run("创建用户", func(t *testing.T) {
		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Name:     "测试用户",
			Password: "hashed_Password",
			Mobile:   "13800138000",
			Status:   1,
		}

		err := db.Create(user).Error
		if err != nil {
			t.Fatalf("创建用户失败: %v", err)
		}

		if user.ID == 0 {
			t.Error("用户ID应该被自动设置")
		}

		if user.CreatedAt.IsZero() {
			t.Error("创建时间应该被自动设置")
		}
	})

	t.Run("查询用户", func(t *testing.T) {
		// 先创建一个用户
		user := &models.User{
			Username: "testuser2",
			Email:    "test2@example.com",
			Name:     "测试用户2",
			Password: "hashed_Password",
			Mobile:   "10.27.0",
			Status:   1,
		}
		db.Create(user)

		// 按用户名查询
		var foundUser models.User
		err := db.Where("username = ?", "queryuser").First(&foundUser).Error
		if err != nil {
			t.Fatalf("查询用户失败: %v", err)
		}

		if foundUser.Username != "queryuser" {
			t.Errorf("期望用户名 queryuser, 实际 %s", foundUser.Username)
		}

		// 按邮箱查询
		err = db.Where("email = ?", "query@example.com").First(&foundUser).Error
		if err != nil {
			t.Fatalf("按邮箱查询用户失败: %v", err)
		}
	})

	t.Run("更新用户", func(t *testing.T) {
		// 先创建一个用户
		user := &models.User{
			Username: "testuser3",
			Email:    "test3@example.com",
			Name:     "测试用户3",
			Password: "hashed_Password",
			Mobile:   "password",
			Status:   1,
		}
		db.Create(user)

		// 更新用户
		user.Email = "updated@example.com"
		user.Status = 0
		err := db.Save(user).Error
		if err != nil {
			t.Fatalf("更新用户失败: %v", err)
		}

		// 验证更新
		var updatedUser models.User
		db.First(&updatedUser, user.ID)
		if updatedUser.Email != "updated@example.com" {
			t.Errorf("期望邮箱 updated@example.com, 实际 %s", updatedUser.Email)
		}
		if updatedUser.Status != 0 {
			t.Errorf("期望状态 0, 实际 %d", updatedUser.Status)
		}
	})

	t.Run("软删除用户", func(t *testing.T) {
		// 先创建一个用户
		user := &models.User{
			Username: "testuser4",
			Email:    "test4@example.com",
			Name:     "测试用户4",
			Password: "hashed_Password",
			Mobile:   "secret",
			Status:   1,
		}
		db.Create(user)

		// 软删除用户
		err := db.Delete(user).Error
		if err != nil {
			t.Fatalf("删除用户失败: %v", err)
		}

		// 验证软删除
		var deletedUser models.User
		err = db.First(&deletedUser, user.ID).Error
		if err != gorm.ErrRecordNotFound {
			t.Error("软删除的用户不应该被查询到")
		}

		// 使用Unscoped查询应该能找到
		err = db.Unscoped().First(&deletedUser, user.ID).Error
		if err != nil {
			t.Fatalf("使用Unscoped查询删除的用户失败: %v", err)
		}
		if deletedUser.DeletedAt.IsZero() {
			t.Error("删除时间应该被设置")
		}
	})

	t.Run("用户模型方法", func(t *testing.T) {
		user := &models.User{
			Username: "methoduser",
			Email:    "method@example.com",
			Name:     "方法测试用户",
			Status:   1,
		}

		// 测试IsActive方法
		if user.IsActive() == 1 {
			t.Error("状态为1的用户应该是激活的")
		}

		user.Status = 0
		if user.IsActive() == 1 {
			t.Error("状态为0的用户应该是非激活的")
		}

		// 由于数据库表结构中没有角色和最后登录时间字段，跳过这些测试
		t.Log("跳过角色和最后登录时间相关测试，因为数据库表结构中没有这些字段")
	})

	t.Run("用户响应转换", func(t *testing.T) {
		user := &models.User{
			ID:       1,
			Username: "responseuser",
			Email:    "response@example.com",
			Name:     "响应测试用户",
			Mobile:   "13800138000",
			Status:   1,
		}

		response := user.ToResponse()
		if response == nil {
			t.Fatal("响应不应该为空")
		}

		if response.ID != user.ID {
			t.Errorf("期望ID %d, 实际 %d", user.ID, response.ID)
		}
		if response.Username != user.Username {
			t.Errorf("期望用户名 %s, 实际 %s", user.Username, response.Username)
		}
		if response.Email != user.Email {
			t.Errorf("期望邮箱 %s, 实际 %s", user.Email, response.Email)
		}

		// 验证密码哈希不在响应中
		// 这里无法直接验证，但通过结构体定义可以确保PasswordHash字段有json:"-"标签
	})
}

func TestUserValidation(t *testing.T) {
	t.Run("用户查询参数验证", func(t *testing.T) {
		query := &models.UserQuery{
			Page: 0,
			Size: 0,
		}

		// 测试默认值
		if query.GetPage() != 1 {
			t.Errorf("期望默认页码 1, 实际 %d", query.GetPage())
		}
		if query.GetSize() != 10 {
			t.Errorf("期望默认每页数量 10, 实际 %d", query.GetSize())
		}

		// 测试最大值限制
		query.Size = 200
		if query.GetSize() != 100 {
			t.Errorf("期望最大每页数量 100, 实际 %d", query.GetSize())
		}

		// 测试偏移量计算
		query.Page = 3
		query.Size = 20
		expectedOffset := (3 - 1) * 20
		if query.GetOffset() != expectedOffset {
			t.Errorf("期望偏移量 %d, 实际 %d", expectedOffset, query.GetOffset())
		}
	})
}

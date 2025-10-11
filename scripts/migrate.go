package main

import (
	"flag"
	"fmt"
	"go_demo/internal/config"
	"go_demo/internal/models"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"log"
	"os"

	"gorm.io/gorm"
)

func main() {
	var (
		configPath = flag.String("config", "./configs/config.yaml", "配置文件路径")
		action     = flag.String("action", "migrate", "操作类型: migrate, rollback, seed")
		help       = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// 连接数据库
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		logger.Fatal("数据库连接失败", logger.Err(err))
	}
	defer database.Close(db)

	// 执行操作
	switch *action {
	case "migrate":
		if err := migrate(db); err != nil {
			logger.Fatal("数据库迁移失败", logger.Err(err))
		}
		logger.Info("数据库迁移完成")
	case "rollback":
		if err := rollback(db); err != nil {
			logger.Fatal("数据库回滚失败", logger.Err(err))
		}
		logger.Info("数据库回滚完成")
	case "seed":
		if err := seed(db); err != nil {
			logger.Fatal("数据库种子数据创建失败", logger.Err(err))
		}
		logger.Info("数据库种子数据创建完成")
	default:
		fmt.Printf("未知操作: %s\n", *action)
		printHelp()
		os.Exit(1)
	}
}

// migrate 执行数据库迁移
func migrate(db *gorm.DB) error {
	logger.Info("开始执行数据库迁移...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// 创建索引
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}

// rollback 回滚数据库（删除表）
func rollback(db *gorm.DB) error {
	logger.Info("开始执行数据库回滚...")

	// 删除表（注意顺序，先删除有外键依赖的表）
	tables := []interface{}{
		&models.User{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			logger.Warn("删除表失败", logger.Err(err))
		}
	}

	return nil
}

// seed 创建种子数据
func seed(db *gorm.DB) error {
	logger.Info("开始创建种子数据...")

	// 创建管理员用户
	adminUser := &models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Name:     "管理员",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // Password
		Mobile:   "13800138000",
		Status:   1,
	}

	// 检查管理员是否已存在
	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		if err := db.Create(adminUser).Error; err != nil {
			return fmt.Errorf("创建管理员用户失败: %w", err)
		}
		logger.Info("管理员用户创建成功", logger.String("username", "admin"))
	} else {
		logger.Info("管理员用户已存在，跳过创建")
	}

	// 创建测试用户
	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "测试用户",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // Password
		Mobile:   "10.27.0",
		Status:   1,
	}

	// 检查测试用户是否已存在
	db.Model(&models.User{}).Where("username = ?", "testuser").Count(&count)
	if count == 0 {
		if err := db.Create(testUser).Error; err != nil {
			return fmt.Errorf("创建测试用户失败: %w", err)
		}
		logger.Info("测试用户创建成功", logger.String("username", "testuser"))
	} else {
		logger.Info("测试用户已存在，跳过创建")
	}

	return nil
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) error {
	logger.Info("开始创建数据库索引...")

	// 用户表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)").Error; err != nil {
		return fmt.Errorf("创建用户状态索引失败: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)").Error; err != nil {
		return fmt.Errorf("创建用户角色索引失败: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)").Error; err != nil {
		return fmt.Errorf("创建用户创建时间索引失败: %w", err)
	}

	logger.Info("数据库索引创建完成")
	return nil
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("数据库迁移工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run scripts/migrate.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -config string    配置文件路径 (默认: ./configs/config.yaml)")
	fmt.Println("  -action string    操作类型: migrate, rollback, seed (默认: migrate)")
	fmt.Println("  -help            显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run scripts/migrate.go                           # 执行迁移")
	fmt.Println("  go run scripts/migrate.go -action=seed              # 创建种子数据")
	fmt.Println("  go run scripts/migrate.go -action=rollback          # 回滚数据库")
	fmt.Println("  go run scripts/migrate.go -config=./config/prod.yaml # 使用指定配置文件")
}

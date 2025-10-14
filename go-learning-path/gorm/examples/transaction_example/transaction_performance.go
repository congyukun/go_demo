package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Email     string    `gorm:"uniqueIndex;size:100;not null"`
	Balance   float64   `gorm:"type:decimal(10,2);default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// 订单模型
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	Status    string    `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// 事务日志模型
type TransactionLog struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"`
	OrderID     uint      `gorm:"index"`
	Amount      float64   `gorm:"type:decimal(10,2)"`
	Type        string    `gorm:"type:varchar(20)"` // debit, credit
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

var db *gorm.DB

func initDB() {
	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Order{}, &TransactionLog{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	fmt.Println("数据库连接成功")
}

// 基础事务示例
func BasicTransaction(userID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查用户余额
		var user User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}

		if user.Balance < amount {
			return fmt.Errorf("余额不足: 当前余额 %.2f, 需要 %.2f", user.Balance, amount)
		}

		// 2. 扣减余额
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return fmt.Errorf("扣减余额失败: %w", err)
		}

		// 3. 创建订单
		order := Order{
			UserID: userID,
			Amount: amount,
			Status: "completed",
		}
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}

		// 4. 记录事务日志
		logEntry := TransactionLog{
			UserID:      userID,
			OrderID:     order.ID,
			Amount:      amount,
			Type:        "debit",
			Description: fmt.Sprintf("消费订单 #%d", order.ID),
		}
		if err := tx.Create(&logEntry).Error; err != nil {
			return fmt.Errorf("记录日志失败: %w", err)
		}

		fmt.Printf("事务执行成功: 用户 %d 消费 %.2f\n", userID, amount)
		return nil
	})
}

// 嵌套事务示例
func NestedTransaction(userID uint, amounts []float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for i, amount := range amounts {
			// 在每个操作中使用嵌套事务
			err := tx.Transaction(func(subTx *gorm.DB) error {
				return processSingleTransaction(subTx, userID, amount, i+1)
			})
			
			if err != nil {
				return err // 任何一个失败都会回滚整个事务
			}
		}
		return nil
	})
}

func processSingleTransaction(tx *gorm.DB, userID uint, amount float64, seq int) error {
	var user User
	if err := tx.First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	if user.Balance < amount {
		return fmt.Errorf("余额不足: 当前余额 %.2f, 需要 %.2f", user.Balance, amount)
	}

	// 扣减余额
	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		return fmt.Errorf("扣减余额失败: %w", err)
	}

	// 创建订单
	order := Order{
		UserID: userID,
		Amount: amount,
		Status: "completed",
	}
	if err := tx.Create(&order).Error; err != nil {
		return fmt.Errorf("创建订单失败: %w", err)
	}

	// 记录日志
	logEntry := TransactionLog{
		UserID:      userID,
		OrderID:     order.ID,
		Amount:      amount,
		Type:        "debit",
		Description: fmt.Sprintf("第%d笔消费订单 #%d", seq, order.ID),
	}
	if err := tx.Create(&logEntry).Error; err != nil {
		return fmt.Errorf("记录日志失败: %w", err)
	}

	fmt.Printf("子事务 %d 执行成功: 消费 %.2f\n", seq, amount)
	return nil
}

// 手动事务控制
func ManualTransactionControl(userID uint, amount float64) error {
	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("事务回滚 due to panic: %v", r)
		}
	}()

	// 业务逻辑
	var user User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("用户不存在: %w", err)
	}

	if user.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("余额不足")
	}

	// 更新余额
	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新余额失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	fmt.Printf("手动事务执行成功: 用户 %d 消费 %.2f\n", userID, amount)
	return nil
}

// 带上下文的事务
func TransactionWithContext(ctx context.Context, userID uint, amount float64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var user User
		if err := tx.First(&user, userID).Error; err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}

		if user.Balance < amount {
			return fmt.Errorf("余额不足")
		}

		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return fmt.Errorf("更新余额失败: %w", err)
		}

		fmt.Printf("带上下文事务执行成功: 用户 %d 消费 %.2f\n", userID, amount)
		return nil
	})
}

// 性能优化：批量插入
func BatchInsertUsers(users []*User) error {
	batchSize := 100 // 每批处理100条记录
	return db.CreateInBatches(users, batchSize).Error
}

// 性能优化：批量更新
func BatchUpdateUserBalance(userIDs []uint, amount float64) error {
	return db.Model(&User{}).
		Where("id IN ?", userIDs).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// 性能监控：查询执行时间
func QueryWithTimeout(ctx context.Context, timeout time.Duration) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("查询超时")
		}
		return nil, err
	}
	return users, nil
}

func main() {
	initDB()

	// 创建测试用户
	testUser := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Balance:  1000.00,
	}
	
	if err := db.Create(testUser).Error; err != nil {
		log.Fatal("创建测试用户失败:", err)
	}

	fmt.Println("=== 基础事务示例 ===")
	if err := BasicTransaction(testUser.ID, 100.00); err != nil {
		log.Printf("基础事务失败: %v", err)
	}

	fmt.Println("\n=== 嵌套事务示例 ===")
	if err := NestedTransaction(testUser.ID, []float64{50.00, 30.00, 20.00}); err != nil {
		log.Printf("嵌套事务失败: %v", err)
	}

	fmt.Println("\n=== 手动事务控制示例 ===")
	if err := ManualTransactionControl(testUser.ID, 40.00); err != nil {
		log.Printf("手动事务失败: %v", err)
	}

	fmt.Println("\n=== 带上下文的事务示例 ===")
	ctx := context.Background()
	if err := TransactionWithContext(ctx, testUser.ID, 25.00); err != nil {
		log.Printf("带上下文事务失败: %v", err)
	}

	// 查询最终余额
	var finalUser User
	db.First(&finalUser, testUser.ID)
	fmt.Printf("\n最终余额: %.2f\n", finalUser.Balance)
}
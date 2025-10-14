# 🗃️ GORM 数据库迁移和版本控制

## 📚 迁移策略

### 1. 自动迁移 vs 手动迁移

**自动迁移 (AutoMigrate)**
```go
// 简单自动迁移
err := db.AutoMigrate(&User{}, &Order{}, &Product{})
if err != nil {
    log.Fatal("自动迁移失败:", err)
}

// 带选项的自动迁移
err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").
    AutoMigrate(&User{})
```

**手动迁移 (推荐生产环境使用)**
```go
// 手动执行DDL语句
func migrateUsersTable(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(50) NOT NULL UNIQUE,
            email VARCHAR(100) NOT NULL UNIQUE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
    `).Error
}
```

### 2. 迁移版本控制

```go
// 迁移版本记录模型
type Migration struct {
    ID        uint      `gorm:"primaryKey"`
    Version   string    `gorm:"size:100;uniqueIndex;not null"`
    Name      string    `gorm:"size:200;not null"`
    AppliedAt time.Time `gorm:"autoCreateTime"`
}

// 迁移管理器
type MigrationManager struct {
    db *gorm.DB
}

func NewMigrationManager(db *gorm.DB) *MigrationManager {
    return &MigrationManager{db: db}
}

// 检查是否已应用迁移
func (m *MigrationManager) IsApplied(version string) (bool, error) {
    var count int64
    err := m.db.Model(&Migration{}).Where("version = ?", version).Count(&count).Error
    return count > 0, err
}

// 记录迁移
func (m *MigrationManager) RecordMigration(version, name string) error {
    migration := &Migration{
        Version: version,
        Name:    name,
    }
    return m.db.Create(migration).Error
}
```

## 🚀 迁移示例

### 1. 基础表结构迁移
```go
// 初始迁移
func InitialMigration(db *gorm.DB) error {
    migrations := []func(*gorm.DB) error{
        createUsersTable,
        createOrdersTable,
        createProductsTable,
    }

    for _, migration := range migrations {
        if err := migration(db); err != nil {
            return err
        }
    }
    return nil
}

func createUsersTable(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(50) NOT NULL UNIQUE,
            email VARCHAR(100) NOT NULL UNIQUE,
            status ENUM('active', 'inactive', 'banned') DEFAULT 'active',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            INDEX idx_users_status (status),
            INDEX idx_users_created_at (created_at)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `).Error
}
```

### 2. 结构变更迁移
```go
// 添加新字段
func AddUserProfileFields(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 添加新字段
        if err := tx.Exec(`
            ALTER TABLE users 
            ADD COLUMN avatar VARCHAR(255) AFTER email,
            ADD COLUMN bio TEXT AFTER avatar,
            ADD COLUMN last_login DATETIME AFTER bio
        `).Error; err != nil {
            return err
        }

        // 添加索引
        if err := tx.Exec(`
            ALTER TABLE users 
            ADD INDEX idx_users_last_login (last_login)
        `).Error; err != nil {
            return err
        }

        return nil
    })
}

// 修改字段类型
func ModifyUserEmailField(db *gorm.DB) error {
    return db.Exec(`
        ALTER TABLE users 
        MODIFY COLUMN email VARCHAR(150) NOT NULL UNIQUE
    `).Error
}
```

### 3. 数据迁移
```go
// 数据迁移示例
func MigrateUserData(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 备份旧数据
        if err := tx.Exec(`
            CREATE TABLE IF NOT EXISTS users_backup 
            AS SELECT * FROM users
        `).Error; err != nil {
            return err
        }

        // 转换数据
        if err := tx.Exec(`
            UPDATE users 
            SET status = CASE 
                WHEN status = 1 THEN 'active'
                WHEN status = 0 THEN 'inactive'
                ELSE 'banned'
            END
        `).Error; err != nil {
            return err
        }

        return nil
    })
}
```

## 🔄 版本控制最佳实践

### 1. 迁移文件组织
```
migrations/
├── 20240101000000_initial_schema.go
├── 20240102000000_add_user_profile.go
├── 20240103000000_create_orders_table.go
└── 20240104000000_add_indexes.go
```

### 2. 时间戳版本控制
```go
// 迁移文件模板
package migrations

import "gorm.io/gorm"

func init() {
    RegisterMigration("20240101000000", "initial_schema", InitialSchema)
}

func InitialSchema(db *gorm.DB) error {
    // 迁移逻辑
    return nil
}
```

### 3. 迁移注册表
```go
var migrationRegistry = make(map[string]MigrationFunc)

type MigrationFunc func(*gorm.DB) error

func RegisterMigration(version, name string, fn MigrationFunc) {
    migrationRegistry[version] = fn
}

func GetMigrations() map[string]MigrationFunc {
    return migrationRegistry
}
```

## 🛡️ 安全迁移实践

### 1. 事务性迁移
```go
func SafeMigration(db *gorm.DB, migrationFunc func(*gorm.DB) error) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := migrationFunc(tx); err != nil {
            return fmt.Errorf("迁移失败: %w", err)
        }
        return nil
    })
}
```

### 2. 回滚策略
```go
// 带有回滚的迁移
func MigrationWithRollback(db *gorm.DB, migrate, rollback func(*gorm.DB) error) error {
    // 执行迁移
    if err := migrate(db); err != nil {
        // 执行回滚
        if rollbackErr := rollback(db); rollbackErr != nil {
            return fmt.Errorf("迁移失败: %w, 回滚也失败: %v", err, rollbackErr)
        }
        return err
    }
    return nil
}

// 回滚函数示例
func rollbackAddColumn(db *gorm.DB) error {
    return db.Exec("ALTER TABLE users DROP COLUMN IF EXISTS new_column").Error
}
```

### 3. 预检查和生产验证
```go
func PreFlightCheck(db *gorm.DB) error {
    // 检查数据库连接
    if err := db.Exec("SELECT 1").Error; err != nil {
        return fmt.Errorf("数据库连接检查失败: %w", err)
    }

    // 检查必要表是否存在
    var tableExists bool
    if err := db.Raw(`
        SELECT COUNT(*) > 0 
        FROM information_schema.tables 
        WHERE table_schema = DATABASE() AND table_name = 'migrations'
    `).Scan(&tableExists).Error; err != nil {
        return err
    }

    if !tableExists {
        return fmt.Errorf("迁移记录表不存在")
    }

    return nil
}
```

## 📊 迁移工具集成

### 1. 命令行迁移工具
```go
// 迁移命令
type MigrateCommand struct {
    db *gorm.DB
}

func (cmd *MigrateCommand) Run(args []string) error {
    switch args[0] {
    case "up":
        return cmd.MigrateUp()
    case "down":
        return cmd.MigrateDown()
    case "status":
        return cmd.MigrationStatus()
    default:
        return fmt.Errorf("未知命令: %s", args[0])
    }
}

func (cmd *MigrateCommand) MigrateUp() error {
    migrations := GetMigrations()
    for version, migrationFunc := range migrations {
        applied, err := IsMigrationApplied(cmd.db, version)
        if err != nil {
            return err
        }
        
        if !applied {
            if err := SafeMigration(cmd.db, migrationFunc); err != nil {
                return err
            }
            if err := RecordMigration(cmd.db, version, "applied"); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### 2. 集成到Go项目
```go
// main.go
func main() {
    // 初始化数据库
    db, err := initDatabase()
    if err != nil {
        log.Fatal("数据库初始化失败:", err)
    }

    // 执行迁移
    if err := runMigrations(db); err != nil {
        log.Fatal("数据库迁移失败:", err)
    }

    // 启动应用
    startApplication(db)
}

func runMigrations(db *gorm.DB) error {
    migrator := NewMigrationManager(db)
    
    // 检查是否需要迁移
    needsMigration, err := migrator.NeedsMigration()
    if err != nil {
        return err
    }

    if needsMigration {
        log.Println("开始执行数据库迁移...")
        if err := migrator.Migrate(); err != nil {
            return err
        }
        log.Println("数据库迁移完成")
    }
    
    return nil
}
```

## 🚨 生产环境注意事项

1. **备份优先**: 执行迁移前务必备份数据库
2. **低峰期执行**: 在业务低峰期执行迁移操作
3. **监控性能**: 监控迁移过程中的数据库性能
4. **回滚计划**: 准备好回滚方案
5. **测试验证**: 在生产环境执行前充分测试

## 📋 迁移检查清单

- [ ] 数据库备份完成
- [ ] 迁移脚本经过测试
- [ ] 回滚方案准备就绪
- [ ] 在低峰期执行
- [ ] 监控系统就绪
- [ ] 团队通知完成

下一步学习：**性能监控和调试技巧**
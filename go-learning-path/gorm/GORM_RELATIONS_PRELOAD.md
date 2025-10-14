# 🔗 GORM 关联关系和预加载优化

## 📚 关联关系类型

### 1. 一对一关系 (One-to-One)
```go
// 用户个人资料模型
type UserProfile struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"uniqueIndex"` // 外键
    Bio       string `gorm:"type:text"`
    Website   string `gorm:"type:varchar(255)"`
    Location  string `gorm:"type:varchar(100)"`
    
    // 属于关系
    User      User   `gorm:"foreignKey:UserID"`
}

// 用户模型添加关联
type User struct {
    ID        uint         `gorm:"primaryKey"`
    // ... 其他字段
    
    // 拥有一个关系
    Profile   UserProfile  `gorm:"foreignKey:UserID"`
}
```

### 2. 一对多关系 (One-to-Many)
```go
// 订单模型
type Order struct {
    ID         uint      `gorm:"primaryKey"`
    UserID     uint      `gorm:"index"` // 外键
    Amount     float64   `gorm:"type:decimal(10,2)"`
    Status     string    `gorm:"type:varchar(20)"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    
    // 属于关系
    User       User      `gorm:"foreignKey:UserID"`
}

// 用户模型添加关联
type User struct {
    ID        uint         `gorm:"primaryKey"`
    // ... 其他字段
    
    // 拥有多个关系
    Orders    []Order      `gorm:"foreignKey:UserID"`
}
```

### 3. 多对多关系 (Many-to-Many)
```go
// 角色模型
type Role struct {
    ID          uint       `gorm:"primaryKey"`
    Name        string     `gorm:"type:varchar(50);uniqueIndex"`
    Description string     `gorm:"type:text"`
    CreatedAt   time.Time  `gorm:"autoCreateTime"`
    
    // 多对多关系
    Users       []User     `gorm:"many2many:user_roles;"`
}

// 用户模型添加多对多关联
type User struct {
    ID        uint         `gorm:"primaryKey"`
    // ... 其他字段
    
    // 多对多关系
    Roles     []Role       `gorm:"many2many:user_roles;"`
}

// 连接表（自动创建，也可自定义）
// type UserRole struct {
//     UserID    uint      `gorm:"primaryKey"`
//     RoleID    uint      `gorm:"primaryKey"`
//     CreatedAt time.Time `gorm:"autoCreateTime"`
// }
```

## 🚀 预加载优化（解决N+1问题）

### 1. 基础预加载
```go
// 错误的N+1查询方式
func (r *userRepository) GetUsersWithOrders() ([]User, error) {
    var users []User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    
    // 这里会产生N+1查询问题！
    for i := range users {
        if err := r.db.Model(&users[i]).Association("Orders").Find(&users[i].Orders); err != nil {
            return nil, err
        }
    }
    
    return users, nil
}

// 正确的预加载方式
func (r *userRepository) GetUsersWithOrders() ([]User, error) {
    var users []User
    err := r.db.Preload("Orders").Find(&users).Error
    return users, err
}
```

### 2. 多级预加载
```go
// 多级关联预加载
func (r *userRepository) GetUsersWithDetails() ([]User, error) {
    var users []User
    
    err := r.db.
        Preload("Profile").                    // 一级预加载
        Preload("Orders").                    // 一级预加载
        Preload("Orders.OrderItems").         // 二级预加载
        Preload("Roles").                     // 一级预加载
        Preload("Roles.Permissions").         // 二级预加载
        Find(&users).Error
    
    return users, err
}
```

### 3. 条件预加载
```go
// 带条件的预加载
func (r *userRepository) GetUsersWithActiveOrders() ([]User, error) {
    var users []User
    
    err := r.db.
        Preload("Orders", "status = ?", "completed"). // 只预加载已完成订单
        Preload("Profile").
        Where("status = ?", 1).
        Find(&users).Error
    
    return users, err
}

// 使用函数进行复杂条件预加载
func (r *userRepository) GetUsersWithRecentOrders() ([]User, error) {
    var users []User
    
    err := r.db.
        Preload("Orders", func(db *gorm.DB) *gorm.DB {
            return db.Where("created_at > ?", time.Now().AddDate(0, -1, 0)). // 最近一个月的订单
                     Order("created_at DESC")
        }).
        Find(&users).Error
    
    return users, err
}
```

### 4. 选择性预加载
```go
// 只预加载特定字段（减少数据传输）
func (r *userRepository) GetUsersWithOrderSummary() ([]User, error) {
    var users []User
    
    err := r.db.
        Preload("Orders", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "user_id", "amount", "status", "created_at")
        }).
        Find(&users).Error
    
    return users, err
}
```

## 🎯 关联操作

### 1. 创建关联
```go
// 创建用户和关联数据
func (r *userRepository) CreateUserWithProfile(user *User, profile *UserProfile) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // 设置关联
        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 2. 更新关联
```go
// 更新关联数据
func (r *userRepository) UpdateUserRoles(userID uint, roleIDs []uint) error {
    var user User
    var roles []Role
    
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 查找用户
        if err := tx.First(&user, userID).Error; err != nil {
            return err
        }
        
        // 查找角色
        if err := tx.Find(&roles, roleIDs).Error; err != nil {
            return err
        }
        
        // 替换关联（清空原有关联，添加新关联）
        return tx.Model(&user).Association("Roles").Replace(&roles)
    })
}
```

### 3. 删除关联
```go
// 清空关联
func (r *userRepository) ClearUserRoles(userID uint) error {
    var user User
    if err := r.db.First(&user, userID).Error; err != nil {
        return err
    }
    
    // 清空所有角色关联
    return r.db.Model(&user).Association("Roles").Clear()
}

// 删除特定关联
func (r *userRepository) RemoveUserRole(userID uint, roleID uint) error {
    var user User
    var role Role
    
    if err := r.db.First(&user, userID).Error; err != nil {
        return err
    }
    
    if err := r.db.First(&role, roleID).Error; err != nil {
        return err
    }
    
    // 删除特定关联
    return r.db.Model(&user).Association("Roles").Delete(&role)
}
```

## ⚡ 性能优化技巧

### 1. 批量预加载
```go
// 批量处理预加载，减少查询次数
func (r *userRepository) BatchPreloadUsers(userIDs []uint) ([]User, error) {
    var users []User
    
    err := r.db.
        Where("id IN ?", userIDs).
        Preload("Orders", func(db *gorm.DB) *gorm.DB {
            return db.Where("status = ?", "completed").Select("id", "user_id", "amount")
        }).
        Preload("Profile", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "user_id", "bio")
        }).
        Find(&users).Error
    
    return users, err
}
```

### 2. 延迟加载控制
```go
// 根据需要动态加载关联
func (r *userRepository) GetUserWithOptionalRelations(userID uint, preloads []string) (*User, error) {
    var user User
    
    db := r.db.Model(&User{})
    
    // 动态添加预加载
    for _, preload := range preloads {
        db = db.Preload(preload)
    }
    
    err := db.First(&user, userID).Error
    return &user, err
}
```

## 🚨 常见陷阱与解决方案

### 1. N+1查询问题
**问题**：循环中查询关联数据
**解决方案**：使用Preload一次性加载

### 2. 循环引用问题
**问题**：模型间相互引用导致序列化问题
**解决方案**：使用DTO或忽略JSON序列化

### 3. 性能问题
**问题**：预加载过多不必要的数据
**解决方案**：选择性预加载特定字段

### 4. 事务一致性
**问题**：关联操作缺乏事务保护
**解决方案**：使用事务包装关联操作

## 📊 预加载性能对比

| 场景 | 查询次数 | 性能 | 推荐 |
|------|---------|------|------|
| N+1查询 | N+1 | 差 | ❌ 避免 |
| 基础预加载 | 2 | 好 | ✅ 推荐 |
| 条件预加载 | 2 | 好 | ✅ 推荐 |
| 多级预加载 | 3+ | 中 | ⚠️ 谨慎使用 |

下一步学习：**事务管理和性能优化**
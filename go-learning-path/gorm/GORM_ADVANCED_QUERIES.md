# 🔍 GORM 查询构建器和复杂查询技巧

## 📚 高级查询技术

### 1. 链式查询构建器
```go
// 基础链式调用
func (r *userRepository) FindActiveUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.
        Where("status = ?", 1).
        Where("last_login > ?", time.Now().AddDate(0, -1, 0)). // 最近一个月有登录
        Order("last_login DESC").
        Limit(100).
        Find(&users).Error
    return users, err
}
```

### 2. 动态查询构建
```go
// 使用Scopes实现可复用的查询条件
func WithStatus(status int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", status)
    }
}

func WithUsernameLike(username string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if username != "" {
            return db.Where("username LIKE ?", "%"+username+"%")
        }
        return db
    }
}

// 使用Scopes
func (r *userRepository) FindUsers(criteria map[string]interface{}) ([]models.User, error) {
    var users []models.User
    
    db := r.db.Model(&models.User{})
    
    if status, ok := criteria["status"]; ok {
        db = db.Scopes(WithStatus(status.(int)))
    }
    
    if username, ok := criteria["username"]; ok {
        db = db.Scopes(WithUsernameLike(username.(string)))
    }
    
    err := db.Find(&users).Error
    return users, err
}
```

### 3. 复杂条件查询
```go
// 多条件组合查询
func (r *userRepository) AdvancedSearch(params map[string]interface{}) ([]models.User, error) {
    var users []models.User
    
    query := r.db.Model(&models.User{})
    
    // 字符串条件
    if username, ok := params["username"]; ok {
        query = query.Where("username LIKE ?", "%"+username.(string)+"%")
    }
    
    // 范围查询
    if createdStart, ok := params["created_start"]; ok {
        query = query.Where("created_at >= ?", createdStart.(time.Time))
    }
    if createdEnd, ok := params["created_end"]; ok {
        query = query.Where("created_at <= ?", createdEnd.(time.Time))
    }
    
    // IN 查询
    if statuses, ok := params["statuses"]; ok {
        query = query.Where("status IN ?", statuses.([]int))
    }
    
    // OR 条件
    if keyword, ok := params["keyword"]; ok {
        query = query.Where(
            r.db.Where("username LIKE ?", "%"+keyword.(string)+"%").
                Or("email LIKE ?", "%"+keyword.(string)+"%").
                Or("name LIKE ?", "%"+keyword.(string)+"%"),
        )
    }
    
    err := query.Find(&users).Error
    return users, err
}
```

### 4. 子查询
```go
// 使用子查询
func (r *userRepository) FindUsersWithRecentLogin() ([]models.User, error) {
    var users []models.User
    
    // 子查询：最近7天登录的用户
    subQuery := r.db.Model(&models.User{}).
        Select("id").
        Where("last_login > ?", time.Now().AddDate(0, 0, -7))
    
    err := r.db.
        Where("id IN (?)", subQuery).
        Find(&users).Error
    
    return users, err
}

// 关联子查询
func (r *userRepository) FindUsersWithOrderCount(minOrders int) ([]models.User, error) {
    var users []models.User
    
    subQuery := r.db.Model(&Order{}).
        Select("user_id, COUNT(*) as order_count").
        Group("user_id").
        Having("order_count >= ?", minOrders)
    
    err := r.db.
        Joins("INNER JOIN (?) AS oc ON users.id = oc.user_id", subQuery).
        Find(&users).Error
    
    return users, err
}
```

### 5. 原生SQL与GORM混合
```go
// 混合使用原生SQL和GORM
func (r *userRepository) ComplexQuery() ([]models.User, error) {
    var users []models.User
    
    err := r.db.
        Where("status = ?", 1).
        Where("EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND orders.status = 'completed')").
        Where("created_at > ?", time.Now().AddDate(-1, 0, 0)). // 一年内注册
        Order("(SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) DESC"). // 按订单数排序
        Find(&users).Error
    
    return users, err
}
```

### 6. 分页优化
```go
// 高效分页查询（使用游标分页）
func (r *userRepository) FindUsersCursor(cursor time.Time, limit int) ([]models.User, error) {
    var users []models.User
    
    err := r.db.
        Where("created_at < ?", cursor).
        Order("created_at DESC").
        Limit(limit).
        Find(&users).Error
    
    return users, err
}

// 使用Keyset分页（性能更好）
func (r *userRepository) FindUsersKeyset(lastID uint, limit int) ([]models.User, error) {
    var users []models.User
    
    query := r.db.
        Order("id ASC").
        Limit(limit)
    
    if lastID > 0 {
        query = query.Where("id > ?", lastID)
    }
    
    err := query.Find(&users).Error
    return users, err
}
```

### 7. 批量操作
```go
// 批量插入
func (r *userRepository) BatchCreate(users []*models.User) error {
    batchSize := 100 // 每批处理100条记录
    
    return r.db.CreateInBatches(users, batchSize).Error
}

// 批量更新
func (r *userRepository) BatchUpdateStatus(ids []uint, status int) error {
    return r.db.Model(&models.User{}).
        Where("id IN ?", ids).
        Update("status", status).Error
}

// 批量删除
func (r *userRepository) BatchDelete(ids []uint) error {
    return r.db.Where("id IN ?", ids).Delete(&models.User{}).Error
}
```

## 🎯 最佳实践

1. **避免N+1查询问题**：使用Preload预加载关联数据
2. **使用索引**：为常用查询字段添加索引
3. **分页优化**：使用Keyset分页代替OFFSET分页
4. **批量操作**：使用CreateInBatches提高性能
5. **查询监控**：开启GORM的Logger监控慢查询

## ⚡ 性能提示

- 使用`Select()`指定需要的字段，避免SELECT *
- 使用`Omit()`排除不需要的字段
- 复杂查询考虑使用原生SQL
- 定期分析查询性能和使用EXPLAIN

下一步学习：**关联关系和预加载优化**
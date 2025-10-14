# 🎯 GORM 深入学习计划

## 📋 当前代码分析

基于您的项目结构，您已经掌握了GORM基础。以下是优化建议和学习路径：

## 🔍 当前代码优化点

### 1. 模型定义优化
```go
// 当前代码
type User struct {
    ID          uint        `gorm:"primarykey"`
    Username    string      `gorm:"uniqueIndex;size:50;not null"`
    // ...
}

// 优化建议
type User struct {
    ID          uint           `gorm:"primaryKey;autoIncrement"`
    Username    string         `gorm:"type:varchar(50);uniqueIndex;not null"`
    Email       string         `gorm:"type:varchar(100);uniqueIndex;not null"`
    Password    string         `gorm:"type:varchar(255);not null"`
    Mobile      string         `gorm:"type:varchar(20);index"`
    Status      int            `gorm:"type:tinyint;default:1;index"`
    CreatedAt   time.Time      `gorm:"autoCreateTime"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

### 2. 查询构建器优化
```go
// 使用Scopes和链式调用
func (r *userRepository) List(query *models.UserQuery) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    err := r.db.Transaction(func(tx *gorm.DB) error {
        // 构建查询
        db := tx.Model(&models.User{})
        if query.Username != "" {
            db = db.Where("username LIKE ?", "%"+query.Username+"%")
        }
        if query.Email != "" {
            db = db.Where("email LIKE ?", "%"+query.Email+"%")
        }
        if query.Status != nil {
            db = db.Where("status = ?", *query.Status)
        }

        // 获取总数
        if err := db.Count(&total).Error; err != nil {
            return err
        }

        // 分页查询
        return db.Offset(query.GetOffset()).
               Limit(query.GetSize()).
               Order("created_at DESC").
               Find(&users).Error
    })

    return users, total, err
}
```

## 🚀 学习路径

### 阶段一：GORM核心进阶
1. **复杂查询构建**
   - 子查询
   - 联合查询
   - 原生SQL与GORM混合使用

2. **关联关系深入**
   - 一对一、一对多、多对多
   - 预加载优化（解决N+1问题）
   - 多态关联

3. **事务管理**
   - 事务隔离级别
   - 分布式事务
   - 事务回滚策略

### 阶段二：高级特性
4. **自定义数据类型**
   - JSON字段处理
   - 枚举类型实现
   - 自定义序列化

5. **钩子函数**
   - 模型生命周期钩子
   - 自定义钩子函数
   - 业务逻辑与数据层分离

6. **数据库迁移**
   - 版本控制
   - 数据迁移脚本
   - 回滚策略

### 阶段三：性能优化
7. **查询性能**
   - 索引优化
   - 查询缓存
   - 分页优化

8. **连接池管理**
   - 连接池配置
   - 超时设置
   - 健康检查

9. **监控调试**
   - SQL日志分析
   - 性能监控
   - 慢查询优化

## 🎯 实践项目建议

### 项目1：电商系统
- 用户、商品、订单模型
- 复杂的关联关系
- 事务处理（库存扣减、订单创建）

### 项目2：博客系统
- 文章、分类、标签
- 多对多关系
- 全文搜索集成

### 项目3：任务管理系统
- 任务、用户、团队
- 状态机实现
- 时间范围查询

## 📚 学习资源

1. **官方文档**: https://gorm.io/docs/
2. **GitHub示例**: https://github.com/go-gorm/examples
3. **最佳实践文章**
4. **视频教程**

## ⚡ 下一步行动

1. 选择一个实践项目开始
2. 深入学习关联关系和预加载
3. 掌握事务管理的最佳实践
4. 学习性能监控和优化技巧

开始从**关联关系和预加载**入手，这是GORM最强大的特性之一！
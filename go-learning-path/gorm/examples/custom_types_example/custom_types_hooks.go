package custom_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// JSON类型自定义数据类型
type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("无法扫描JSON类型: %v", value)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*j = result
	return nil
}

// 枚举类型
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

// Value 实现 driver.Valuer 接口
func (s UserStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan 实现 sql.Scanner 接口
func (s *UserStatus) Scan(value interface{}) error {
	if value == nil {
		*s = StatusInactive
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("无法扫描状态类型: %v", value)
	}

	*s = UserStatus(str)
	return nil
}

// 自定义时间类型（处理时区）
type LocalTime time.Time

// Value 实现 driver.Valuer 接口
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t).Format("2006-01-02 15:04:05"), nil
}

// Scan 实现 sql.Scanner 接口
func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		*t = LocalTime(time.Time{})
		return nil
	}

	timeStr, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("无法扫描时间类型: %v", value)
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(timeStr))
	if err != nil {
		return err
	}

	*t = LocalTime(parsedTime)
	return nil
}

// MarshalJSON 自定义JSON序列化
func (t LocalTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02 15:04:05"))
}

// UnmarshalJSON 自定义JSON反序列化
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return err
	}

	*t = LocalTime(parsedTime)
	return nil
}

// 用户模型（包含自定义类型）
type User struct {
	ID          uint       `gorm:"primaryKey"`
	Username    string     `gorm:"uniqueIndex;size:50;not null"`
	Email       string     `gorm:"uniqueIndex;size:100;not null"`
	Status      UserStatus `gorm:"type:varchar(20);default:'inactive'"`
	Preferences JSONMap    `gorm:"type:json"`     // 使用自定义JSON类型
	LastLogin   LocalTime  `gorm:"type:datetime"` // 使用自定义时间类型
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`

	// 虚拟字段（不存储到数据库）
	IsOnline bool `gorm:"-"`
}

// 钩子函数示例

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("准备创建用户: %s\n", u.Username)

	// 设置默认值
	if u.Status == "" {
		u.Status = StatusInactive
	}

	if u.Preferences == nil {
		u.Preferences = JSONMap{
			"theme":         "light",
			"language":      "zh-CN",
			"notifications": true,
		}
	}

	// 验证逻辑
	if len(u.Username) < 3 {
		return fmt.Errorf("用户名长度至少3个字符")
	}

	return nil
}

// AfterCreate 创建后的钩子
func (u *User) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("用户创建成功: %s (ID: %d)\n", u.Username, u.ID)

	// 可以在这里发送通知、创建相关记录等
	log.Printf("新用户注册: %s, 邮箱: %s", u.Username, u.Email)

	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	fmt.Printf("准备更新用户: %s\n", u.Username)

	// 验证逻辑
	if u.Status == StatusBanned {
		// 检查是否可以封禁用户
		// 这里可以添加额外的业务逻辑
	}

	return nil
}

// AfterUpdate 更新后的钩子
func (u *User) AfterUpdate(tx *gorm.DB) error {
	fmt.Printf("用户更新成功: %s\n", u.Username)

	// 记录更新日志等
	if u.Status == StatusBanned {
		log.Printf("用户被封禁: %s", u.Username)
	}

	return nil
}

// BeforeDelete 删除前的钩子
func (u *User) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("准备删除用户: %s\n", u.Username)

	// 防止误删除重要用户
	if u.Username == "admin" {
		return fmt.Errorf("不能删除管理员账户")
	}

	return nil
}

// AfterDelete 删除后的钩子
func (u *User) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("用户删除成功: %s\n", u.Username)

	// 清理相关数据、发送通知等
	log.Printf("用户数据已删除: %s", u.Username)

	return nil
}

// 自定义查询方法
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (u *User) GetPreference(key string) interface{} {
	if u.Preferences == nil {
		return nil
	}
	return u.Preferences[key]
}

func (u *User) SetPreference(key string, value interface{}) {
	if u.Preferences == nil {
		u.Preferences = make(JSONMap)
	}
	u.Preferences[key] = value
}

var db *gorm.DB

func initDB() {
	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_custom?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	fmt.Println("数据库连接成功")
}

// 演示自定义类型的使用
func demoCustomTypes() {
	// 创建用户
	user := &User{
		Username: "john_doe",
		Email:    "john@example.com",
		Status:   StatusActive,
		Preferences: JSONMap{
			"theme":         "dark",
			"language":      "en-US",
			"notifications": false,
			"settings": map[string]interface{}{
				"email_alerts": true,
				"sms_alerts":   false,
			},
		},
		LastLogin: LocalTime(time.Now()),
		IsOnline:  true, // 虚拟字段
	}

	// 创建用户（会触发钩子函数）
	if err := db.Create(user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		return
	}

	// 查询用户
	var foundUser User
	if err := db.First(&foundUser, user.ID).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
		return
	}

	fmt.Printf("用户查询结果:\n")
	fmt.Printf("ID: %d\n", foundUser.ID)
	fmt.Printf("用户名: %s\n", foundUser.Username)
	fmt.Printf("状态: %s\n", foundUser.Status)
	fmt.Printf("是否活跃: %t\n", foundUser.IsActive())
	fmt.Printf("偏好设置: %+v\n", foundUser.Preferences)
	fmt.Printf("最后登录: %v\n", time.Time(foundUser.LastLogin))

	// 使用自定义方法
	fmt.Printf("主题偏好: %v\n", foundUser.GetPreference("theme"))

	// 更新偏好设置
	foundUser.SetPreference("theme", "light")
	foundUser.SetPreference("font_size", 14)

	if err := db.Save(&foundUser).Error; err != nil {
		log.Printf("更新用户失败: %v", err)
		return
	}

	fmt.Printf("更新后的偏好设置: %+v\n", foundUser.Preferences)

	// 演示枚举类型的使用
	fmt.Printf("状态检查:\n")
	fmt.Printf("是否是活跃状态: %t\n", foundUser.Status == StatusActive)
	fmt.Printf("是否是被封禁状态: %t\n", foundUser.Status == StatusBanned)

	// 切换状态
	foundUser.Status = StatusBanned
	if err := db.Save(&foundUser).Error; err != nil {
		log.Printf("更新状态失败: %v", err)
		return
	}

	fmt.Printf("用户状态已更新为: %s\n", foundUser.Status)
}

// 演示钩子函数
func demoHooks() {
	fmt.Println("\n=== 演示钩子函数 ===")

	// 创建测试用户
	testUser := &User{
		Username: "test_user",
		Email:    "test@example.com",
		Status:   StatusActive,
	}

	if err := db.Create(testUser).Error; err != nil {
		log.Printf("创建测试用户失败: %v", err)
		return
	}

	// 更新用户
	testUser.Status = StatusInactive
	if err := db.Save(testUser).Error; err != nil {
		log.Printf("更新用户失败: %v", err)
		return
	}

	// 删除用户（会触发删除钩子）
	if err := db.Delete(testUser).Error; err != nil {
		log.Printf("删除用户失败: %v", err)
		return
	}
}

func main() {
	initDB()

	fmt.Println("=== 演示自定义数据类型 ===")
	demoCustomTypes()

	demoHooks()

	fmt.Println("\n演示完成!")
}

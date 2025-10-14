# 🏗️ Go 设计模式和最佳实践

## 📚 常用设计模式

### 1. 依赖注入 (Dependency Injection)
```go
// 接口定义
type UserRepository interface {
    GetByID(id int) (*User, error)
    Create(user *User) error
}

type UserService interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// 实现
type userService struct {
    repo UserRepository
    logger *zap.Logger
}

// 使用构造函数注入依赖
func NewUserService(repo UserRepository, logger *zap.Logger) UserService {
    return &userService{
        repo: repo,
        logger: logger,
    }
}

func (s *userService) GetUser(id int) (*User, error) {
    s.logger.Info("获取用户", zap.Int("id", id))
    return s.repo.GetByID(id)
}

// 使用Wire进行依赖注入
// wire.go
func InitializeUserService() UserService {
    wire.Build(
        NewUserService,
        NewUserRepository,
        NewLogger,
    )
    return &userService{}
}
```

### 2. 工厂模式 (Factory Pattern)
```go
// 数据库连接工厂
type DBConfig struct {
    Host     string
    Port     int
    Username string
    PASSWORD     string
    Database string
}

type DBType string

const (
    MySQL    DBType = "mysql"
    PostgreSQL DBType = "postgres"
)

// 数据库工厂接口
type DBFactory interface {
    CreateConnection(config DBConfig) (*sql.DB, error)
}

// MySQL工厂实现
type MySQLFactory struct{}

func (f *MySQLFactory) CreateConnection(config DBConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
        config.Username, config.PASSWORD, config.Host, config.Port, config.Database)
    return sql.Open("mysql", dsn)
}

// 工厂方法
func NewDBFactory(dbType DBType) DBFactory {
    switch dbType {
    case MySQL:
        return &MySQLFactory{}
    case PostgreSQL:
        return &PostgreSQLFactory{}
    default:
        panic("不支持的数据库类型")
    }
}
```

### 3. 装饰器模式 (Decorator Pattern)
```go
// 基础接口
type UserService interface {
    GetUser(id int) (*User, error)
}

// 基础实现
type basicUserService struct {
    repo UserRepository
}

// 缓存装饰器
type cachedUserService struct {
    service UserService
    cache   *redis.Client
    ttl     time.Duration
}

func NewCachedUserService(service UserService, cache *redis.Client, ttl time.Duration) UserService {
    return &cachedUserService{
        service: service,
        cache:   cache,
        ttl:     ttl,
    }
}

func (s *cachedUserService) GetUser(id int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // 尝试从缓存获取
    cached, err := s.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // 缓存未命中，从服务获取
    user, err := s.service.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果
    userJSON, _ := json.Marshal(user)
    s.cache.Set(ctx, cacheKey, userJSON, s.ttl)
    
    return user, nil
}

// 日志装饰器
type loggedUserService struct {
    service UserService
    logger  *zap.Logger
}

func NewLoggedUserService(service UserService, logger *zap.Logger) UserService {
    return &loggedUserService{
        service: service,
        logger:  logger,
    }
}

func (s *loggedUserService) GetUser(id int) (*User, error) {
    start := time.Now()
    defer func() {
        s.logger.Info("GetUser执行完成",
            zap.Int("id", id),
            zap.Duration("duration", time.Since(start)))
    }()
    
    return s.service.GetUser(id)
}
```

### 4. 观察者模式 (Observer Pattern)
```go
// 事件定义
type Event interface {
    Name() string
}

type UserCreatedEvent struct {
    User User
    Timestamp time.Time
}

func (e UserCreatedEvent) Name() string { return "user.created" }

// 观察者接口
type Observer interface {
    Notify(event Event)
}

// 主题接口
type Subject interface {
    Register(observer Observer)
    Unregister(observer Observer)
    NotifyObservers(event Event)
}

// 具体实现
type EventBus struct {
    observers map[Observer]bool
    mu        sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        observers: make(map[Observer]bool),
    }
}

func (b *EventBus) Register(observer Observer) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.observers[observer] = true
}

func (b *EventBus) NotifyObservers(event Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    
    for observer := range b.observers {
        go observer.Notify(event) // 异步通知
    }
}

// 邮件通知观察者
type EmailNotifier struct{}

func (n *EmailNotifier) Notify(event Event) {
    if userCreated, ok := event.(UserCreatedEvent); ok {
        // 发送欢迎邮件
        sendWelcomeEmail(userCreated.User)
    }
}
```

## 🎯 Go 语言最佳实践

### 1. 错误处理最佳实践
```go
// 自定义错误类型
type AppError struct {
    Code    string
    Message string
    Op      string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// 错误包装和展开
func GetUser(id int) (*User, error) {
    user, err := userRepo.GetByID(id)
    if err != nil {
        return nil, &AppError{
            Code:    "USER_NOT_FOUND",
            Message: "用户不存在",
            Op:      "GetUser",
            Err:     err,
        }
    }
    return user, nil
}

// 错误检查
func HandleError() {
    user, err := GetUser(1)
    if err != nil {
        var appErr *AppError
        if errors.As(err, &appErr) {
            switch appErr.Code {
            case "USER_NOT_FOUND":
                // 处理用户不存在
            case "DATABASE_ERROR":
                // 处理数据库错误
            }
        }
        log.Printf("错误: %v", err)
    }
}
```

### 2. 并发安全最佳实践
```go
// 使用sync.Map处理并发map访问
type SafeCache struct {
    data sync.Map
}

func (c *SafeCache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *SafeCache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}

// 使用RWMutex保护共享资源
type Counter struct {
    mu    sync.RWMutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}

// 使用context控制goroutine生命周期
func ProcessWithTimeout(ctx context.Context, data []string) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    results := make(chan result, len(data))
    errCh := make(chan error, 1)
    
    var wg sync.WaitGroup
    for _, item := range data {
        wg.Add(1)
        go func(item string) {
            defer wg.Done()
            select {
            case <-ctx.Done():
                return
            default:
                res, err := processItem(ctx, item)
                if err != nil {
                    select {
                    case errCh <- err:
                    default:
                    }
                    return
                }
                results <- res
            }
        }(item)
    }
    
    go func() {
        w极Wait()
        close(results)
    }()
    
    // 收集结果...
}
```

### 3. 性能优化最佳实践
```go
// 对象池减少内存分配
var userPool = sync.Pool{
    New: func() interface{} {
        return &User{}
    },
}

func GetUserFromPool() *User {
    return userPool.Get().(*User)
}

func PutUserToPool(user *User) {
    user.Reset() // 重置对象状态
    userPool.Put(user)
}

// 使用strings.Builder构建字符串
func BuildLargeString(items []string) string {
    var builder strings.Builder
    builder.Grow(len(items) * 100) // 预分配空间
    
    for _, item := range items {
        builder.WriteString(item)
        builder.WriteString("\n")
    }
    
    return builder.String()
}

// 避免不必要的内存分配
func ProcessUsers(users []User) {
    // 不好的方式：每次循环都创建新变量
    for i := 0; i < len(users); i++ {
        user := users[i] // 复制
        process(&user)
    }
    
    // 好的方式：使用指针或索引
    for i := range users {
        process(&users[i])
    }
}
```

### 4. 测试最佳实践
```go
// 表格驱动测试
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name     string
        items    []Item
        expected float64
        hasError bool
    }{
        {
            name:     "空列表",
            items:    []Item{},
            expected: 0,
            hasError: false,
        },
        {
            name: "正常计算",
            items: []Item{
                {Price: 10, Quantity: 2},
                {Price: 5, Quantity: 1},
            },
            expected: 25,
            hasError: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateTotal(tt.items)
            if tt.hasError {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}

// 使用mock进行单元测试
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    expectedUser := &User{ID: 1, Name: "John"}
    
    mockRepo.On("GetByID", 1).Return(expectedUser, nil)
    
    service := NewUserService(mockRepo, zap.NewNop())
    user, err := service.GetUser(1)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

## 🏆 代码组织最佳实践

### 1. 项目结构
```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── models/
│   └── middleware/
├── pkg/
│   ├── database/
│   ├── logger/
│   └── utils/
├── configs/
├── deployments/
├── scripts/
└── tests/
```

### 2. 包设计原则
- **单一职责**: 每个包只做一件事
- **低耦合**: 包之间依赖最小化
- **高内聚**: 相关功能放在同一个包
- **明确接口**: 提供清晰的API边界

### 3. 配置管理
```go
// 使用环境变量和配置文件
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 📊 性能调优 checklist

- [ ] 使用pprof分析性能瓶颈
- [ ] 优化数据库查询和索引
- [ ] 使用连接池和缓存
- [ ] 减少内存分配和GC压力
- [ ] 并发优化和goroutine管理
- [ ] 监控和告警设置

## 🔧 代码质量 checklist

- [ ] 编写单元测试和集成测试
- [ ] 使用静态代码分析工具
- [ ] 遵循代码规范和安全最佳实践
- [ ] 文档和注释完善
- [ ] 代码审查和持续集成

## 🎉 总结

通过掌握这些设计模式和最佳实践，您将能够：

1. **编写更健壮的代码**: 使用适当的错误处理和并发控制
2. **提高代码可维护性**: 遵循清晰的架构和设计原则
3. **优化性能**: 减少资源消耗和提高响应速度
4. **便于测试**: 编写可测试的代码和完整的测试套件
5. **易于扩展**: 使用灵活的设计支持未来需求变化

继续实践这些模式，您将成为更优秀的Go开发者！
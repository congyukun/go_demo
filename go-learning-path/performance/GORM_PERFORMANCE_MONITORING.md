# 📊 GORM 性能监控和调试技巧

## 🎯 性能监控基础

### 1. GORM 日志配置
```go
import (
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

// 详细日志配置
func setupDB() *gorm.DB {
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second,   // 慢查询阈值
            LogLevel:                  logger.Info,    // 日志级别
            IgnoreRecordNotFoundError: true,           // 忽略记录未找到错误
            Colorful:                  true,          // 彩色打印
        },
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    return db
}

// 不同环境的日志配置
func getLoggerForEnv(env string) logger.Interface {
    config := logger.Config{
        SlowThreshold: time.Second,
        Colorful:      env != "production",
    }

    switch env {
    case "development":
        config.LogLevel = logger.Info
    case "test":
        config.LogLevel = logger.Warn
    case "production":
        config.LogLevel = logger.Error
        config.IgnoreRecordNotFoundError = true
    }

    return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), config)
}
```

### 2. 慢查询监控
```go
// 自定义慢查询处理器
type SlowQueryHandler struct {
    threshold time.Duration
}

func NewSlowQueryHandler(threshold time.Duration) *SlowQueryHandler {
    return &SlowQueryHandler{threshold: threshold}
}

func (h *SlowQueryHandler) Handle(ctx context.Context, sql string, duration time.Duration) {
    if duration > h.threshold {
        log.Printf("🚨 慢查询警告: %s (执行时间: %v)", sql, duration)
        // 可以发送到监控系统、记录到文件等
        metrics.RecordSlowQuery(sql, duration)
    }
}

// 集成到GORM
db.Callback().Query().After("gorm:query").Register("slow_query_monitor", func(db *gorm.DB) {
    if db.Statement != nil && db.Statement.SQL != nil {
        duration := time.Since(db.Statement.StartTime)
        slowQueryHandler.Handle(db.Statement.Context, db.Statement.SQL.String(), duration)
    }
})
```

## 📈 性能分析工具

### 1. 使用 pprof 分析
```go
import (
    _ "net/http/pprof"
    "net/http"
)

// 启动pprof服务器
func startPProfServer() {
    go func() {
        log.Println("PProf服务器启动在 :6060")
        log.Println(http.ListenAndServe(":6060", nil))
    }()
}

// 在GORM中集成性能追踪
func withQueryMetrics(db *gorm.DB, queryName string) *gorm.DB {
    return db.Set("query_name", queryName).Callback().
        Query().Before("gorm:query").
        Register("query_metrics_start", func(db *gorm.DB) {
            db.Set("query_start_time", time.Now())
        }).
        After("gorm:query").
        Register("query_metrics_end", func(db *gorm.DB) {
            startTime, ok := db.Get("query_start_time")
            if !ok {
                return
            }
            
            duration := time.Since(startTime.(time.Time))
            queryName, _ := db.Get("query_name")
            
            metrics.RecordQueryDuration(queryName.(string), duration)
        })
}
```

### 2. SQL 执行计划分析
```go
// 解释查询执行计划
func ExplainQuery(db *gorm.DB, query interface{}, args ...interface{}) (string, error) {
    var result string
    err := db.Raw("EXPLAIN "+query.(string), args...).Scan(&result).Error
    return result, err
}

// 自动分析慢查询
func AutoAnalyzeSlowQueries(db *gorm.DB, threshold time.Duration) {
    db.Callback().Query().After("gorm:query").Register("auto_analyze", func(db *gorm.DB) {
        duration := time.Since(db.Statement.StartTime)
        if duration > threshold && db.Statement.SQL != nil {
            explain, err := ExplainQuery(db, db.Statement.SQL.String(), db.Statement.Vars...)
            if err == nil {
                log.Printf("慢查询分析:\nSQL: %s\n执行计划: %s\n耗时: %v", 
                    db.Statement.SQL.String(), explain, duration)
            }
        }
    })
}
```

## 🔍 调试技巧

### 1. 查询调试工具
```go
// SQL调试器
type SQLDebugger struct {
    enabled bool
}

func (d *SQLDebugger) BeforeQuery(ctx context.Context, sql string, vars []interface{}) (context.Context, error) {
    if d.enabled {
        log.Printf("🔍 执行查询: %s\n参数: %+v", sql, vars)
    }
    return ctx, nil
}

func (d *SQLDebugger) AfterQuery(ctx context.Context, sql string, vars []interface{}, err error) {
    if d.enabled && err != nil {
        log.Printf("❌ 查询错误: %s\nSQL: %s", err.Error(), sql)
    }
}

// 使用调试器
debugger := &SQLDebugger{enabled: true}
db.Debug().Where("name = ?", "john").Find(&users)
```

### 2. 连接池监控
```go
import (
    "database/sql"
    "time"
)

// 监控数据库连接池
func MonitorConnectionPool(db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        return
    }

    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            stats := sqlDB.Stats()
            log.Printf("连接池状态: "+
                "OpenConnections: %d, "+
                "InUse: %d, "+
                "Idle: %d, "+
                "WaitCount: %d, "+
                "WaitDuration: %v",
                stats.OpenConnections,
                stats.InUse,
                stats.Idle,
                stats.WaitCount,
                stats.WaitDuration)
        }
    }()
}

// 连接池配置优化
func OptimizeConnectionPool(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    // 设置连接池参数
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
    sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
    sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 连接最大空闲时间

    return nil
}
```

## 📊 性能指标收集

### 1. 自定义指标收集
```go
// 性能指标结构
type QueryMetrics struct {
    QueryName     string
    Duration      time.Duration
    RowsAffected  int64
    Error         error
    Timestamp     time.Time
}

// 指标收集器
type MetricsCollector struct {
    metrics chan QueryMetrics
}

func NewMetricsCollector(bufferSize int) *MetricsCollector {
    return &MetricsCollector{
        metrics: make(chan QueryMetrics, bufferSize),
    }
}

func (c *MetricsCollector) Record(metrics QueryMetrics) {
    select {
    case c.metrics <- metrics:
    default:
        // 缓冲区满，丢弃指标
        log.Println("指标缓冲区已满")
    }
}

func (c *MetricsCollector) Start() {
    go func() {
        for metric := range c.metrics {
            // 发送到监控系统
            sendToMonitoringSystem(metric)
        }
    }()
}

// 集成到GORM
func InstrumentDB(db *gorm.DB, collector *MetricsCollector) *gorm.DB {
    return db.Callback().
        Query().After("gorm:query").
        Register("query_metrics", func(db *gorm.DB) {
            duration := time.Since(db.Statement.StartTime)
            collector.Record(QueryMetrics{
                QueryName:    getQueryName(db),
                Duration:     duration,
                RowsAffected: db.Statement.RowsAffected,
                Error:        db.Error,
                Timestamp:    time.Now(),
            })
        })
}
```

### 2. Prometheus 集成
```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus指标
var (
    queryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "gorm_query_duration_seconds",
        Help:    "GORM查询执行时间",
        Buckets: prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms 到 16s
    }, []string{"query", "status"})

    queryCount = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "gorm_query_total",
        Help: "GORM查询总数",
    }, []string{"query", "status"})
)

// Prometheus监控中间件
func WithPrometheusMetrics(db *gorm.DB) *gorm.DB {
    return db.Callback().
        Query().After("gorm:query").
        Register("prometheus_metrics", func(db *gorm.DB) {
            duration := time.Since(db.Statement.StartTime).Seconds()
            status := "success"
            if db.Error != nil {
                status = "error"
            }

            queryName := getQueryName(db)
            
            queryDuration.WithLabelValues(queryName, status).Observe(duration)
            queryCount.WithLabelValues(queryName, status).Inc()
        })
}
```

## 🐛 常见性能问题调试

### 1. N+1 查询问题检测
```go
// N+1查询检测器
func DetectNPlusOneQueries(db *gorm.DB, threshold int) {
    var queryCount int
    var lastQueryTime time.Time

    db.Callback().Query().After("gorm:query").Register("nplus1_detector", func(db *gorm.DB) {
        now := time.Now()
        
        // 如果在短时间内有大量查询，可能是N+1问题
        if now.Sub(lastQueryTime) < 100*time.Millisecond {
            queryCount++
            if queryCount > threshold {
                log.Printf("⚠️ 可能的N+1查询问题: 在 %v 内执行了 %d 次查询", 
                    now.Sub(lastQueryTime), queryCount)
                // 打印堆栈跟踪以帮助定位问题
                debug.PrintStack()
            }
        } else {
            queryCount = 1
        }
        
        lastQueryTime = now
    })
}
```

### 2. 内存泄漏检测
```go
// 内存使用监控
func MonitorMemoryUsage() {
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()

        var m runtime.MemStats
        for range ticker.C {
            runtime.ReadMemStats(&m)
            log.Printf("内存使用: Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
                m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
        }
    }()
}
```

## 🚀 性能优化建议

### 1. 查询优化
- 使用 `Select()` 指定需要的字段
- 避免使用 `*` 查询所有字段
- 使用索引优化查询
- 批量操作代替循环操作

### 2. 连接池优化
- 合理设置最大连接数
- 监控连接池状态
- 及时关闭不再使用的连接

### 3. 缓存策略
- 使用Redis缓存频繁查询的结果
- 实现查询结果缓存
- 设置合理的缓存过期时间

## 📋 性能检查清单

- [ ] 慢查询监控配置
- [ ] 连接池参数优化
- [ ] N+1查询检测启用
- [ ] 性能指标收集
- [ ] 内存使用监控
- [ ] 定期性能分析

下一步学习：**学习容器化和部署（Docker, Kubernetes）**
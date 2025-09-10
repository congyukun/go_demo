package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
	ReqLogger     *zap.Logger // 专门用于请求日志的logger
)

// LogConfig 日志配置
type LogConfig struct {
	Level       string `yaml:"level"`
	Format      string `yaml:"format"`
	OutputPath  string `yaml:"output_path"`
	ReqLogPath  string `yaml:"req_log_path"`
	MaxSize     int    `yaml:"max_size"`
	MaxBackup   int    `yaml:"max_backup"`
	MaxAge      int    `yaml:"max_age"`
	Compress    bool   `yaml:"compress"`
}

// Init 初始化日志
func Init(cfg LogConfig) error {
	// 创建日志目录
	logDir := filepath.Dir(cfg.OutputPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建请求日志目录
	if cfg.ReqLogPath != "" {
		reqLogDir := filepath.Dir(cfg.ReqLogPath)
		if err := os.MkdirAll(reqLogDir, 0755); err != nil {
			return fmt.Errorf("创建请求日志目录失败: %w", err)
		}
	}

	// 设置日志级别
	level := getLogLevel(cfg.Level)

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色输出
	}

	// 设置时间格式
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置主日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.OutputPath,
		MaxSize:    cfg.MaxSize,   // MB
		MaxBackups: cfg.MaxBackup, // 保留旧文件的最大个数
		MaxAge:     cfg.MaxAge,    // 保留旧文件的最大天数
		Compress:   cfg.Compress,  // 是否压缩/归档旧文件
		LocalTime:  true,          // 使用本地时间
	}

	// 创建写入器
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriterSync := zapcore.AddSync(fileWriter)

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, level),  // 控制台输出
		zapcore.NewCore(encoder, fileWriterSync, level), // 文件输出
	)

	// 创建主logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	SugaredLogger = Logger.Sugar()

	// 如果配置了请求日志路径，则创建请求日志logger
	if cfg.ReqLogPath != "" {
		// 配置请求日志轮转
		reqFileWriter := &lumberjack.Logger{
			Filename:   cfg.ReqLogPath,
			MaxSize:    cfg.MaxSize,   // MB
			MaxBackups: cfg.MaxBackup, // 保留旧文件的最大个数
			MaxAge:     cfg.MaxAge,    // 保留旧文件的最大天数
			Compress:   cfg.Compress,  // 是否压缩/归档旧文件
			LocalTime:  true,          // 使用本地时间
		}

		// 创建请求日志写入器
		reqFileWriterSync := zapcore.AddSync(reqFileWriter)

		// 创建请求日志核心
		reqCore := zapcore.NewTee(
			// zapcore.NewCore(encoder, consoleWriter, level),        // 控制台输出
			zapcore.NewCore(encoder, reqFileWriterSync, level),    // 请求日志文件输出
		)

		

		// 创建请求日志logger
		ReqLogger = zap.New(reqCore, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	return nil
}

// getLogLevel 获取日志级别
func getLogLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Sync 同步日志缓冲区
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// Close 关闭日志
func Close() {
	Sync()
}

// 便捷方法
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// ReqInfo 请求日志信息级别
func ReqInfo(msg string, fields ...zap.Field) {
	if ReqLogger != nil {
		ReqLogger.Info(msg, fields...)
	} else {
		// 如果没有专门的请求日志logger，则使用主logger
		Info(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sugar方法
func Debugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	SugaredLogger.Fatalf(template, args...)
}

// 字段构造函数
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

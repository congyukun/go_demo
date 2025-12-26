package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
	sugarLogger  *zap.SugaredLogger
)

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level" yaml:"level"`               // 日志级别: debug, info, warn, error
	Format     string `mapstructure:"format" yaml:"format"`             // 日志格式: json, console
	OutputPath string `mapstructure:"output_path" yaml:"output_path"`   // 日志输出路径
	ReqLogPath string `mapstructure:"req_log_path" yaml:"req_log_path"` // 请求日志路径
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`         // 单个日志文件最大大小(MB)
	MaxBackup  int    `mapstructure:"max_backup" yaml:"max_backup"`     // 保留的旧日志文件数量
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`           // 保留的旧日志文件天数
	Compress   bool   `mapstructure:"compress" yaml:"compress"`         // 是否压缩旧日志文件
}

// Init 初始化日志系统
func Init(config LogConfig) error {
	// 设置日志级别
	level := zapcore.InfoLevel
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 设置编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置编码器
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志输出
	var cores []zapcore.Core

	// 控制台输出
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	// 文件输出
	if config.OutputPath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   config.OutputPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackup,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		fileCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = globalLogger.Sugar()

	return nil
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GetLogger 获取全局logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// 如果未初始化，使用默认配置
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// GetSugarLogger 获取全局sugar logger
func GetSugarLogger() *zap.SugaredLogger {
	if sugarLogger == nil {
		sugarLogger = GetLogger().Sugar()
	}
	return sugarLogger
}

// Field 日志字段类型
type Field = zapcore.Field

// 常用字段构造函数
func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// 日志记录函数
func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	GetLogger().Panic(msg, fields...)
}

// Sugar logger 函数
func Debugf(template string, args ...interface{}) {
	GetSugarLogger().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	GetSugarLogger().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetSugarLogger().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetSugarLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	GetSugarLogger().Fatalf(template, args...)
}

// Sync 同步日志缓冲区
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

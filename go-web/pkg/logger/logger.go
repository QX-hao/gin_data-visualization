package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志实例
var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// LogConfig 日志配置结构体
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	Format     string `yaml:"format"`      // 日志格式: json, console
	Output     string `yaml:"output"`      // 输出目标: file, console, both
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max_backups"` // 最大备份文件数量
	MaxAge     int    `yaml:"max_age"`     // 日志文件最大保存天数
	Compress   bool   `yaml:"compress"`   // 是否压缩备份文件
}

// InitLogger 初始化日志记录器
func InitLogger(config *LogConfig) error {
	// 设置默认配置
	if config == nil {
		config = &LogConfig{
			Level:      "info",
			Format:     "json",
			Output:     "both",
			FilePath:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		}
	}

	// 创建日志目录
	if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if config.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建写入器
	var writeSyncer zapcore.WriteSyncer
	switch config.Output {
	case "file":
		writeSyncer = getFileWriteSyncer(config)
	case "console":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "both":
		fileSyncer := getFileWriteSyncer(config)
		consoleSyncer := zapcore.AddSync(os.Stdout)
		writeSyncer = zapcore.NewMultiWriteSyncer(fileSyncer, consoleSyncer)
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 设置日志级别
	level := getLogLevel(config.Level)
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建日志记录器
	Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	Sugar = Logger.Sugar()

	return nil
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 获取文件写入器
func getFileWriteSyncer(config *LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 日志记录方法
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

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sugar Logger 方法
func Debugf(template string, args ...interface{}) {
	Sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Sugar.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Sugar.Fatalf(template, args...)
}

// 带字段的Sugar Logger方法
func Debugw(msg string, keysAndValues ...interface{}) {
	Sugar.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Sugar.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Sugar.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Sugar.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Sugar.Fatalw(msg, keysAndValues...)
}

// 同步日志缓冲区
func Sync() error {
	return Logger.Sync()
}
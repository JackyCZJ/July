package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*Options)

// 默认值
const (
	// DefaultLevel 默认日志级别
	DefaultLevel = zapcore.DebugLevel
	// DefaultFileName 默认日志输出路径
	DefaultFilename = "./july-log.log"
	// DefaultMaxSize 默认单日志文件大小
	DefaultMaxSize = 100
)

var zapLevel = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type Options struct {
	Logger     *zap.SugaredLogger
	Filename   string
	MaxSize    int32
	LocalTime  bool
	Compress   bool
	MaxBackups int
	Level      zapcore.Level
}

// Logger 设置日志对象
func Logger(logger *zap.SugaredLogger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// Filename 日志文件名
func Filename(filename string) Option {
	return func(o *Options) {
		o.Filename = filename
	}
}

// MaxSize 单日志文件大小上限
func MaxSize(maxSize int32) Option {
	return func(o *Options) {
		o.MaxSize = maxSize
	}
}

// LocalTime 是否使用本地时间
func LocalTime(localTime bool) Option {
	return func(o *Options) {
		o.LocalTime = localTime
	}
}

// Compress 是否压缩日志文件
func Compress(compress bool) Option {
	return func(o *Options) {
		o.Compress = compress
	}
}

// Level 日志级别
func Level(l string) Option {
	return func(o *Options) {
		if level, ok := zapLevel[l]; ok {
			o.Level = level
		} else {
			o.Level = DefaultLevel
		}
	}
}

// FilterOutFunc 设置中间件忽略函数列表
//func FilterOutFunc(filterOutFunc) Option {
//	return func(o *Options) {
//		o.FilterOutFunc = filterOutFunc
//	}
//}

//最大备份数量
func MaxBackups(a int) Option {
	return func(options *Options) {
		options.MaxBackups = a
	}
}

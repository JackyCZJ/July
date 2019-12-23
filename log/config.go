package log

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	Options       *Options
	SugaredLogger *zap.SugaredLogger
}

func configure(zap *ZapLogger, ops ...Option) {
	// 默认值
	zap.Options.LocalTime = true
	zap.Options.Compress = true
	// 处理设置参数
	for _, o := range ops {
		o(zap.Options)
	}
	// 参数为空时默认值
	if zap.Options.Filename == "" {
		zap.Options.Filename = DefaultFilename
	}
	if zap.Options.MaxSize <= 0 {
		zap.Options.MaxSize = DefaultMaxSize
	}
	if zap.Options.Level < -1 || zap.Options.Level > 5 {
		zap.Options.Level = DefaultLevel
	}
}

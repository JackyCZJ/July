package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logworker *ZapLogger

func NewZapLogger(opts ...Option) *ZapLogger {
	zapLogger := &ZapLogger{
		Options: new(Options),
	}
	// 配置
	configure(zapLogger, opts...)
	// 未设置日志对象，则创建一个
	if zapLogger.Options.Logger == nil {
		// 创建zap日志对象
		syncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   zapLogger.Options.Filename,
			MaxSize:    int(zapLogger.Options.MaxSize),
			LocalTime:  zapLogger.Options.LocalTime,
			Compress:   zapLogger.Options.Compress,
			MaxBackups: zapLogger.Options.MaxBackups,
		})
		pEncoder := zap.NewDevelopmentEncoderConfig()
		pEncoder.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式
		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder //命令行颜色
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(pEncoder), syncWriter, zap.NewAtomicLevelAt(zapLogger.Options.Level)),                //生产环境核心 ， 输出日志至文件
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapLogger.Options.Level)), //开发环境核心，输出带颜色参数的日志至命令行
		)

		logger := zap.New(core, zap.AddCaller())
		zapLogger.SugaredLogger = logger.Sugar()
	} else {
		zapLogger.SugaredLogger = zapLogger.Options.Logger
	}

	Logworker = zapLogger
	return zapLogger
}

//use ZapLogger to implement middleware
func (zap *ZapLogger) StreamClient(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
	//todo: filter func
	defer func() {
		if err != nil {
			zap.SugaredLogger.Errorw("Streamer function", "method", method, "err", err)
		} else if zap.Options.Level == zapcore.DebugLevel {
			zap.SugaredLogger.Debugw("Streamer function", "method", method, "err", err)
		}
	}()
	cs, err = streamer(ctx, desc, cc, method, opts...)
	return
}

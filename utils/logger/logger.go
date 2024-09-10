package logger

import (
	"configService/conf"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"tc.nari.com/debug"
)

var logger *zap.SugaredLogger

// 自定义日志输出，部分输出到debug

// Init 初始化日志
func Init(args ...string) {
	cfg := conf.Cfg
	// 初始化debug
	if len(args) >= 2 {
		err := debug.Init(cfg.Zookeeper.Addr, cfg.Service.Debug, "service", fmt.Sprintf("%s:%s", args[0], args[1]))
		if err != nil {
			fmt.Println("调试服务初始化失败:", err)
		}
	}

	// 初始化日志
	encodeConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "caller",
		StacktraceKey: "trace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		/*EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},*/
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	level := zap.NewAtomicLevelAt(zapcore.Level(cfg.Runtime.LogLevel))

	logr := lumberjack.Logger{
		Filename:   cfg.Logger.Filename,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxDays,
		Compress:   cfg.Logger.Compress,
	}

	hook := zapcore.AddSync(&logr)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encodeConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), hook),
		&level,
	)
	// 设置调用CallerSkip(2)跳过两层调用
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

}

// Close 关闭日志
func Close() {
	_ = debug.Logout()
}

// Debug 调试类型日志, 输出到debug及日志文件
func Debug(args ...interface{}) {
	if len(args) <= 0 {
		return
	} else if len(args) == 1 {
		logger.Debug(args...)
		return
	} else {
		content := fmt.Sprint(args[1:]...)
		logger.Debugf("[%v]%v", args[0], content)
		if title, ok := args[0].(string); ok {
			debug.Debug(title, content)
		}
	}
}

// Debugf 调试类型日志, 输出到日志文件
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Debugw 调试类型日志, 输出到日志文件
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

// Info 提示类型日志, 输出到debug和日志文件
func Info(args ...interface{}) {
	if len(args) <= 0 {
		return
	} else if len(args) == 1 {
		logger.Info(args...)
		return
	} else {
		content := fmt.Sprint(args[1:]...)
		logger.Infof("[%v]%v", args[0], content)
		if title, ok := args[0].(string); ok {
			debug.Debug(title, content)
		}
	}
}

// Infof 提示类型日志, 输出到日志文件
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Infow 提示类型日志, 输出到日志文件
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// Warn 告警类型日志, 输出到debug及日志文件
func Warn(args ...interface{}) error {
	if len(args) <= 0 {
		return errors.New("")
	} else if len(args) == 1 {
		logger.Warn(args...)
	} else {
		content := fmt.Sprint(args[1:]...)
		logger.Warnf("[%v]%v", args[0], content)
		if title, ok := args[0].(string); ok {
			debug.Debug(title, content)
		}
	}
	return errors.New(fmt.Sprint(args...))
}

// Warnf 告警类型日志, 输出到日志文件
func Warnf(template string, args ...interface{}) error {
	logger.Warnf(template, args...)
	if len(args) <= 0 && template == "" {
		return errors.New("")
	}
	if template == "" {
		return errors.New(fmt.Sprint(args...))
	}
	return errors.New(fmt.Sprintf(template, args...))
}

// Warnw 告警类型日志, 输出到日志文件
func Warnw(msg string, keysAndValues ...interface{}) error {
	logger.Warnw(msg, keysAndValues...)
	if len(keysAndValues) <= 0 && msg == "" {
		return errors.New("")
	}
	if msg == "" {
		return errors.New(fmt.Sprint(keysAndValues...))
	}
	return errors.New(msg + "," + fmt.Sprint(keysAndValues...))
}

// Error 错误类型日志, 输出到debug和日志文件, 返回error
func Error(args ...interface{}) error {
	if len(args) <= 0 {
		return errors.New("")
	} else if len(args) == 1 {
		logger.Error(args...)
	} else {
		content := fmt.Sprint(args[1:]...)
		logger.Errorf("[%v]%v", args[0], content)
		if title, ok := args[0].(string); ok {
			debug.Debug(title, content)
		}
	}
	return errors.New(fmt.Sprint(args...))
}

// Errorf 错误类型日志, 输出到日志文件, 返回error
func Errorf(template string, args ...interface{}) error {
	logger.Errorf(template, args...)
	if len(args) <= 0 && template == "" {
		return errors.New("")
	}
	if template == "" {
		return errors.New(fmt.Sprint(args...))
	}
	return errors.New(fmt.Sprintf(template, args...))
}

// Errorw 错误类型日志, 输出到日志文件, 返回error
func Errorw(msg string, keysAndValues ...interface{}) error {
	logger.Errorw(msg, keysAndValues...)
	if len(keysAndValues) <= 0 && msg == "" {
		return errors.New("")
	}
	if msg == "" {
		return errors.New(fmt.Sprint(keysAndValues...))
	}
	return errors.New(msg + "," + fmt.Sprint(keysAndValues...))
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	logger.DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

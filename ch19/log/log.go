package log

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field
type Logger interface {
	Debug(msg string)
	DebugC(context context.Context, msg string)
	Debugf(format string, args ...interface{})
	DebugfC(context context.Context, args ...interface{})
	DebugW(msg string, keysAndValues ...interface{})
	DebugWC(context context.Context, keysAndValues ...interface{})
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	zapLogger *zap.Logger
}

// Debug implements Logger.
func (z *zapLogger) Debug(msg string) {
	z.zapLogger.Debug(msg)
}

// DebugC implements Logger.
func (z *zapLogger) DebugC(context context.Context, msg string) {
	panic("unimplemented")
}

// DebugW implements Logger.
func (z *zapLogger) DebugW(msg string, keysAndValues ...interface{}) {
	panic("unimplemented")
}

// DebugWC implements Logger.
func (z *zapLogger) DebugWC(context context.Context, keysAndValues ...interface{}) {
	panic("unimplemented")
}

// Debugf implements Logger.
func (z *zapLogger) Debugf(format string, args ...interface{}) {
	panic("unimplemented")
}

// DebugfC implements Logger.
func (z *zapLogger) DebugfC(context context.Context, args ...interface{}) {
	panic("unimplemented")
}

var (
	defaultLogger = New(NewOptions())
	mu            sync.Mutex
)

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	loggerConfig := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
	}

	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
	}

	return logger
}

func Init(opt *Options) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = New(opt)
}

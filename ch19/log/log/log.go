package log

import (
	"context"
	stdlog "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogHelper interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
}

// Init 用指定的选项初始化 logger
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// New 通过opts创建 logger, 可以通过命令参数自定义
func New(opts *Options) *Logger {
	if opts == nil {
		opts = NewOptions() // 如果不传递 opts 则使用默认创建
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil { // 获取 opts.Level 等级 需要转化一下
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder //将Level序列化为全大写字符串。例如, InfoLevel被序列化为INFO
	// 当输出到本地路径时，禁止带颜色
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder // 将Level序列化为全大写字符串并添加颜色。例如，InfoLevel被序列化为INFO，并被标记为蓝色。
	}

	encoderConfig := zapcore.EncoderConfig{ // 设置 编码器
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,   // 每行结束 换行
		EncodeLevel:    encodeLevel,                 // log 等级
		EncodeTime:     timeEncoder,                 // TODO 注释看看
		EncodeDuration: milliSecondsDurationEncoder, // TODO 注释看看
		EncodeCaller:   zapcore.ShortCallerEncoder,  // TODO 注释看看
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // logger 启动的最小级别
		Development:       opts.Development,               // 是否是开发模式
		DisableCaller:     opts.DisableCaller,             // 停止使用调用函数的文件名和行号注释日志, 在默认情况下 所有日志都有注释
		DisableStacktrace: opts.DisableStacktrace,         // 禁用自动捕获堆栈跟踪。默认情况下，在 dev 环境中捕获WarnLevel及以上级别的日志，在 pro 环境中捕获ErrorLevel及以上级别的日志
		Sampling: &zap.SamplingConfig{ // 设置采样策略。nil 禁用采样。
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format, // 设置记录器的编码 有效值 console 和 json
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,      // 输出 默认输出到标准输出
		ErrorOutputPaths: opts.ErrorOutputPaths, // 错误输出 默认输出到标准错误
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1)) // zap.AddCallerSkip 因为日志记录器本身也会被记录，所以一般会跳过 1 层，以便获取真正的调用者信息
	if err != nil {
		panic(err)
	}
	logger := &Logger{
		Logger: l,

		//有时我们稍微封装了一下记录日志的方法，但是我们希望输出的文件名和行号是调用封装函数的位置。这时可以使用zap.AddCallerSkip(skip int)向上跳 1 层：
		skipCaller: l.WithOptions(zap.AddCallerSkip(1)),

		minLevel:         zapLevel,
		errorStatusLevel: zap.ErrorLevel,
		caller:           true, // 填写 true 在 span 中传递 调用的函数名 调用的文件位置 调用的行位置
		withTraceID:      true, // 打印log时 额外打印 trace_id 信息
		//stackTrace:       true, // 填写 true 在 span 中传递调用栈信息
	}

	/*
		当应用程序中同时使用了 log 模块和 zap 日志库时，如果不对两者进行整合，就可能会出现一些问题：
			1. 日志信息不统一：log 模块的日志输出格式和 zap 日志库的输出格式可能不一致，导致在分析和处理日志时比较麻烦。
			2. 日志信息遗漏：如果不将 log 模块的日志信息也输入到 zap 日志库中，就可能会出现部分日志信息遗漏的情况。
			3. 日志信息错乱：如果将 log 模块的日志信息和 zap 日志库的信息分别输出到不同的文件或控制台中，就可能会出现信息错乱或混杂的情况。
		因此，我们可以使用 zap.RedirectStdLog 方法将 log 模块的日志信息重定向到 zap 日志库中，从而避免上述问题。
		这个方法的作用是将 log 模块的输出重定向到 zap 日志库中，从而使得两者的日志信息都能够被 zap 日志库捕获和处理，同时也能够保证日志信息的统一、完整和正确输出。
	*/
	zap.RedirectStdLog(l)

	return logger
}

// ZapLogger used for other log wrapper such as klog.
func ZapLogger() *zap.Logger {
	return std.Logger
}

// CheckIntLevel used for other log wrapper such as klog which return if logging a
// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.Logger.Check(lvl, "")

	return checkEntry != nil
}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	std.Logger.Debug(msg, fields...)
}

// DebugC method output debug level log.
func DebugC(ctx context.Context, msg string, fields ...Field) {
	std.DebugContext(ctx, msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	std.Logger.Sugar().Debugf(format, v...)
}

// DebugfC method output debug level log.
func DebugfC(ctx context.Context, format string, v ...interface{}) {
	std.DebugfContext(ctx, format, v...)
}

// Debugw method output debug level log.
func Debugw(msg string, keysAndValues ...interface{}) {
	std.Logger.Sugar().Debugw(msg, keysAndValues...)
}

func DebugwC(ctx context.Context, msg string, keysAndValues ...interface{}) {
	std.DebugfContext(ctx, msg, keysAndValues...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	std.Logger.Info(msg, fields...)
}

func InfoC(ctx context.Context, msg string, fields ...Field) {
	std.InfoContext(ctx, msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	std.Logger.Sugar().Infof(format, v...)
}

func InfofC(ctx context.Context, format string, v ...interface{}) {
	std.InfofContext(ctx, format, v...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	std.Logger.Warn(msg, fields...)
}

func WarnC(ctx context.Context, msg string, fields ...Field) {
	std.WarnContext(ctx, msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	std.Logger.Sugar().Warnf(format, v...)
}

func WarnfC(ctx context.Context, format string, v ...interface{}) {
	std.WarnfContext(ctx, format, v...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	std.Logger.Error(msg, fields...)
}

func ErrorC(ctx context.Context, msg string, fields ...Field) {
	std.ErrorContext(ctx, msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	std.Logger.Sugar().Errorf(format, v...)
}

func ErrorfC(ctx context.Context, format string, v ...interface{}) {
	std.ErrorfContext(ctx, format, v...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	std.Logger.Panic(msg, fields...)
}

func PanicC(ctx context.Context, msg string, fields ...Field) {
	std.PanicContext(ctx, msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	std.Logger.Sugar().Panicf(format, v...)
}

func PanicfC(ctx context.Context, format string, v ...interface{}) {
	std.PanicfContext(ctx, format, v...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	std.Logger.Fatal(msg, fields...)
}

func FatalC(ctx context.Context, msg string, fields ...Field) {
	std.PanicContext(ctx, msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	std.Logger.Sugar().Fatalf(format, v...)
}

func FatalfC(ctx context.Context, format string, v ...interface{}) {
	std.FatalfContext(ctx, format, v...)
}

func StdInfoLogger() *stdlog.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.Logger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

func Flush() { std.Flush() }

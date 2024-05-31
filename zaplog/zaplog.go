package zaplog

import (
	"context"
	"fmt"
	"github.com/archine/logger-zap/zaplog/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// ZapLogger zap logger
type ZapLogger struct {
	*zap.Logger
	fieldApplyFunc func(ctx context.Context) []zap.Field
}

var defaultLogger *ZapLogger

// Init creates a new zap logger
func Init(conf *config.Config) error {
	if conf == nil {
		return fmt.Errorf("log configuration is nil")
	}
	if conf.Level == "" {
		conf.Level = "debug"
	}
	if conf.Formatter == "" {
		conf.Formatter = "console"
	}
	if conf.Syncer == nil {
		conf.Syncer = zapcore.AddSync(os.Stdout)
	}
	level, err := zapcore.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	conf.Options = append(conf.Options, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	ec := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		FunctionKey:    zapcore.OmitKey,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		StacktraceKey:  "stacktrace",
	}
	var encoder zapcore.Encoder
	if conf.Formatter == "json" {
		encoder = zapcore.NewJSONEncoder(ec)
	} else {
		if conf.ConsoleSeparator == "" {
			conf.ConsoleSeparator = " | "
		}
		ec.ConsoleSeparator = conf.ConsoleSeparator
		encoder = zapcore.NewConsoleEncoder(ec)
	}

	core := zapcore.NewCore(encoder, conf.Syncer, level)
	logger := zap.New(core, conf.Options...)

	defaultLogger = &ZapLogger{logger, conf.ApplyFields}
	return nil
}

// clone
func (z *ZapLogger) clone() *ZapLogger {
	nz := *z
	return &nz
}

func (z *ZapLogger) With(fields ...zap.Field) *ZapLogger {
	if len(fields) == 0 {
		return z
	}
	nz := z.clone()
	nz.Logger = nz.Logger.With(fields...)
	return nz
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

// WithContext returns a new logger with the given context
func WithContext(ctx context.Context) *ZapLogger {
	if defaultLogger.fieldApplyFunc == nil {
		return defaultLogger
	}
	return defaultLogger.With(defaultLogger.fieldApplyFunc(ctx)...)
}

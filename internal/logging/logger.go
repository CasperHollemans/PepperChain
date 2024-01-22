package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type ZapLogger struct {
	zap *zap.Logger
}

func NewZapLogger(debug bool) *ZapLogger {
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{zap: logger}
}

func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.zap.Sugar().Debugf(msg, args...)
}

func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.zap.Sugar().Infof(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.zap.Sugar().Warnf(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.zap.Sugar().Errorf(msg, args...)
}

func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
	l.zap.Sugar().Fatalf(msg, args...)
}

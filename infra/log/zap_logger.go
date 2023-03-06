package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func (l zapLogger) Print(v ...interface{}) {
	l.sugaredLogger.Info(v)
}

func (l zapLogger) Warn(v ...interface{}) {
	l.sugaredLogger.Warn(v)
}

func (l zapLogger) Error(v ...interface{}) {
	l.sugaredLogger.Error(v)
}

func (l zapLogger) Printf(format string, v ...interface{}) {
	l.sugaredLogger.Infof(format, v)
}

func (l zapLogger) Errorf(format string, v ...interface{}) {
	l.sugaredLogger.Errorf(format, v)
}

func newZapLogger() *zapLogger {
	var logger *zap.Logger
	var err error
	if config.IsProd {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.FunctionKey = "func"
		cfg.DisableStacktrace = true
		logger, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	} else {
		cfg := zap.NewDevelopmentConfig()
		cfg.DisableStacktrace = true
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		logger, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	}

	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			sentry.CaptureException(err)
		}
	}(logger)
	return &zapLogger{
		sugaredLogger: logger.Sugar(),
	}
}

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

type zapLogger struct {
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
}

func (l zapLogger) Print(v ...interface{}) {
	l.SugaredLogger.Info(v)
}

func (l zapLogger) Warn(v ...interface{}) {
	l.SugaredLogger.Warn(v)
}

func (l zapLogger) Error(v ...interface{}) {
	l.SugaredLogger.Error(v)
}

func (l zapLogger) Printf(format string, v ...interface{}) {
	l.SugaredLogger.Infof(format, v)
}

func (l zapLogger) Errorf(format string, v ...interface{}) {
	l.SugaredLogger.Errorf(format, v)
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
		logger, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	}

	defer logger.Sync()
	return &zapLogger{
		Logger:        logger,
		SugaredLogger: logger.Sugar(),
	}
}

package lgr

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"

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

func (l zapLogger) Sync() {
	_ = l.sugaredLogger.Sync()
}

func newZapLogger() *zapLogger {
	var logger *zap.Logger
	var err error
	if config.IsProd {
		core, err := adapter.New(
			adapter.SetClientOptions(
				axiom.SetOrganizationID(config.App.AxiomOrganizationID),
				axiom.SetAPITokenConfig(config.App.AxiomToken),
			),
			adapter.SetDataset(config.App.AxiomDataset),
		)
		if err != nil {
			panic(err)
		}
		logger = zap.New(core)
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
	return &zapLogger{
		sugaredLogger: logger.Sugar(),
	}
}

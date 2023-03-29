package lgr

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"
	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

var zapper *zap.SugaredLogger

func newZapLogger() *zap.SugaredLogger {
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
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
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
	return logger.Sugar()
}

func SetupLogger() {
	zapper = newZapLogger()
}

func Print(v ...interface{}) {
	zapper.Info(v)
}

func Warn(v ...interface{}) {
	zapper.Warn(v)
}

func Error(v ...interface{}) {
	zapper.Error(v)
}

func Printf(format string, v ...interface{}) {
	zapper.Infof(format, v)
}

func Errorf(format string, v ...interface{}) {
	zapper.Errorf(format, v)
}

func Sync() {
	if err := zapper.Sync(); err != nil {
		sentry.CaptureException(err)
	}
}

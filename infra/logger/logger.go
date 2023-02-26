package logger

import "go.uber.org/zap"

// Logger is the logging interface of the project
type Logger interface {
	Print(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Printf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

var Log Logger
var RequestLogger *zap.Logger

func SetupLogger() {
	zapLogger := newZapLogger()
	Log = zapLogger
	RequestLogger = zapLogger.Logger
}

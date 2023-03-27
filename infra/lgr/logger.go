package lgr

// Logger is the logging interface of the project
type Logger interface {
	Print(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Printf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Sync()
}

var Default Logger

func SetupLogger() {
	Default = newZapLogger()
}

package gobacktest

type BackLog interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})

	Debugf(format string, v ...interface{})
	Debug(v ...interface{})

	Infof(format string, v ...interface{})
	Info(v ...interface{})

	Warnf(format string, v ...interface{})
	Warn(v ...interface{})

	Errorf(format string, v ...interface{})
	Error(v ...interface{})

	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})

	Panicf(format string, v ...interface{})
	Panic(v ...interface{})
}

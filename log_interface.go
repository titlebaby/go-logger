package logger

type LogInterface interface {
	init()
	SetLevel(level int)
	Debug(format string, arg ...interface{})
	Trace(format string, arg ...interface{})
	Info(format string, arg ...interface{})
	Warn(format string, arg ...interface{})
	Error(format string, arg ...interface{})
	Fatal(format string, arg ...interface{})
	Close()
}

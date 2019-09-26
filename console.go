package logger

import "C"
import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func (c *ConsoleLogger) init() {

}

func NewConsoleLogger(config map[string]string) (log LogInterface, err error) {

	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not  found log_level")
		return
	}
	level := getLevel(logLevel)
	log = &ConsoleLogger{
		level: level,
	}
	return
}
func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		c.level = LogLevelDebug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	//str := fmt.Sprintf(format,arg...)
	if c.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelDebug, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	logData := writeLog( LogLevelFatal, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *ConsoleLogger) Close() {

}

package logger

import (
	"fmt"
	"path"
	"time"
)

const (
	LogLevelDebug  = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

const(
	logSplitTypeHour = iota
	logSplitTypeSize
)

func getLevel(level string) int  {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	}
	return LogLevelDebug
}

func getLevelText(level int) string  {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	}
	return "DEBUG"
}
/**
1。 当业务调试用打日志的方法时，我们把日志相关的数据写入到chan（队列）
2。 然后我们有一个后台的线程不断的从chan里面获取这些日志，最终写入到文件
 */
/*
同步写法
func writeLog(file *os.File, level int, format string, args ...interface{}) {
	//if f.level > LogLevelDebug {
	//	return
	//}
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05")
	levelStr := getLevelText(level)
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(file, "%s %s (%s %s : %d )%s\n", nowStr, levelStr, fileName, funcName, lineNo, msg)
	fmt.Fprintln(file)
}*/
func writeLog(level int, format string, args ...interface{}) *logData  {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05")
	levelStr := getLevelText(level)
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	logData := &logData{
		Message:msg,
		TimeStr:nowStr,
		LevelStr:levelStr,
		Filename:fileName,
		FuncName:funcName,
		LineNo:lineNo,
		warnAndFatal:false,
	}
	if level == LogLevelWarn || level == LogLevelError || level == LogLevelFatal {
		logData.warnAndFatal = true
	}
	return logData
}

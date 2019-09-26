package logger

import (
	"runtime"
)
//	存储log数据
type logData struct {
	Message  string
	TimeStr string
	LevelStr string
	Filename string
	FuncName string
	LineNo int
	warnAndFatal bool
}


func GetLineInfo() (fileName string, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}

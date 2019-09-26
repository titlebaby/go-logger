package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	logDataChan   chan *logData //可以存任何类型 int string object ，放指针性能更好 避免数据的拷贝
	logSplitType  int
	logSplitSize  int64
	lastSplitHour int
}

func (f *FileLogger) init() {
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed ,err:%v\n", filename, err))
	}
	f.file = file
	//写错误日志和fatal日志文件
	filename = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v\n", f.logName, err))
	}
	f.warnFile = file

	//后台写log的线程
	go f.writeLogBackground()
}

func (f *FileLogger) splitFileHour(warnFile bool)  {
	now := time.Now()
	hour := now.Hour()

	if hour == f.lastSplitHour {
		return
	}
	f.lastSplitHour = hour
	var backupFilename string
	var filename string
	if warnFile {
		backupFilename = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		filename = fmt.Sprintf("%s/%s.log.wf",
			f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s/%s.log%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		filename = fmt.Sprintf("%s/%s.log",
			f.logPath, f.logName)
	}
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	file.Close()
	//备份
	os.Rename(filename, backupFilename)
	//重新打开日志文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		//打开失败 直接程序挂掉
		return
	}
	//更新
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}
// 按大小分日志
func (f *FileLogger) splitFileSize(warnFile bool)  {
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	//获取当前文件的基本信息 大小 最后修改时间
	statInfo , err := file.Stat()

	if err != nil {
		return
	}

	fileSize := statInfo.Size()
	fmt.Println(fileSize, f.logSplitSize)
	if fileSize <= f.logSplitSize {
		return
	}

	var backupFilename string
	var filename string
	now := time.Now()
	//备份 到秒数 保证没有冲突 不会一秒钟产生10兆的日志
	if warnFile {
		backupFilename = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),now.Hour(), now.Minute(),now.Second())
		filename = fmt.Sprintf("%s/%s.log.wf",
			f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s/%s.log%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),now.Hour() ,now.Minute(),now.Second())
		filename = fmt.Sprintf("%s/%s.log",
			f.logPath, f.logName)
	}
	file.Close()
	//实际是重命名文件
	os.Rename(filename, backupFilename)
	fmt.Println(backupFilename)

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		return
	}
	if warnFile {
		f.warnFile =file
	} else {
		f.file = file
	}
}


func (f *FileLogger) checkoutSplitFile(warnFile bool) {
	if f.logSplitType == logSplitTypeHour {
		f.splitFileHour(warnFile)
		return
	} else {
		f.splitFileSize(warnFile)
		return
	}
}

//后台写log的线程
func (f *FileLogger) writeLogBackground() {
	for logData := range f.logDataChan {
		var file *os.File = f.file
		if logData.warnAndFatal {
			file = f.warnFile
		}
		f.checkoutSplitFile(logData.warnAndFatal)
		fmt.Fprintf(file, "%s %s (%s %s : %d )%s\n", logData.TimeStr,
			logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)

	}
}

//func NewFileLogger(level int, logPath, logName string) LogInterface {
func NewFileLogger(config map[string]string) (log LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name1")
		return
	}
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}
	level := getLevel(logLevel)
	//查看是否设置通道大小
	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}
	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}
	//查看是否设置日志分割类型、以及日志分割大小
	var logSplitType int = logSplitTypeHour
	var logSplitSize int64
	logSplitStr, ok := config["log_split_type"]
	if !ok {
		logSplitStr = "hour"
	} else {
		if logSplitStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857"
			}
			logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64) // 十进制的64位
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = logSplitTypeSize
		} else {
			logSplitType = logSplitTypeHour
		}
	}

	log = &FileLogger{
		level:         level,
		logPath:       logPath,
		logName:       logName,
		logDataChan:   make(chan *logData, chanSize),
		logSplitSize:  logSplitSize,
		logSplitType:  logSplitType,
		lastSplitHour: time.Now().Hour(),
	}
	log.init()
	return
}

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		f.level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	//str := fmt.Sprintf(format,arg...)
	if f.level > LogLevelDebug {
		return
	}
	//writeLog(f.file, LogLevelDebug, format, args...)
	logData := writeLog(LogLevelDebug, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}
func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}
func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}
func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	logData := writeLog(LogLevelFatal, format, args...)
	//判断管道是否满  满了就直接丢弃*****
	select {
	case f.logDataChan <- logData:
	default:
		//满了 直接default丢弃，防止管道堵塞
	}
	f.logDataChan <- logData
}
func (f *FileLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}

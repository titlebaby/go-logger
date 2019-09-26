##go 日志库学习
    主要实现了对文件和控制台日志的打印，充分运用了go的渠道和线程，同时对文件日志进行了按小时或者大小的切割。
1. 使用
```
package main
import (
	"fmt"
	"github.com/pgxy/logger"
	"time"
)
var log logger.LogInterface
func initLogger(logPath, logName string, level string )  (err error) {
	//2次封装
	m := make(map[string]string)
	m["log_path"] = logPath
	m["log_name"] = logName
	m["log_level"] = level
	m["log_split_type"] = "size"
	err = logger.InitLogger("file", m)
	//err = logger.InitLogger("console", m)

	if err!=nil {
		fmt.Println("yesy",err)
		return err
	}
	logger.Debug("init logger success!!!")
	return
}

func Run()  {
	for {
		logger.Debug("user server is running，panic: runtime error: invalid memory address or nil pointer dereference")
	}

}

func main()  {
	err := initLogger("/Users/linger/logs/","user_server", "debug")
	if err!=nil {
	    fmt.Println("init failed :",err)
        return 
	}
	Run()
	return

}


```
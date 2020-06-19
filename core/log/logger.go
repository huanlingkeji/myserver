package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

//支持日志级别
//文件滚动存储
//日志格式
//
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func Info(v ...interface{}) {
	now := time.Now()
	prefix := fmt.Sprintf(`%v/%2d/%2d-%2d:%2d:%2d`, now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	file, line := getFileLine()
	logger.Println(prefix, " ", file, " ", line, " ", v)
}

func getFileLine() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file, line
}

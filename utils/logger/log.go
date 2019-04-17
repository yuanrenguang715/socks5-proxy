package logger

import (
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var logDir = "./log"
var logPath = logDir + "/%Y%m%d.log"
var rotationTime = 24 * 3600 //按天数生成日志
//var maxAge = 30 * 24 * 3600  //自动删除30天前的日志

// Init 初始化日志模块
func Init() {
	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(logDir, os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	rl, err := rotatelogs.New(
		logPath,
		//rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Second),
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetFlags(0)
	log.SetOutput(rl)
}

// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Info(args ...interface{}) {
	write("[I]  ", "", args...)
}

// Infof calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	write("[I]  ", format, args...)
}

// Error calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Error(args ...interface{}) {
	write("[E]  ", "", args...)
}

// Errorf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	write("[E]  ", format, args...)
}

func write(level, format string, args ...interface{}) {
	log.SetPrefix(time.Now().Format("[2006-01-02 15:04:05]  ") + level)
	if len(format) == 0 {
		log.Print(args...)
	} else {
		log.Printf(format, args...)
	}
}

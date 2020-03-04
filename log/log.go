package log

import (
	"fmt"
	"log"
)

func Info(args ...interface{}) {
	log.Println("INFO: " + fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	log.Println("INFO: " + fmt.Sprintf(format, args...))
}

func Warning(args ...interface{}) {
	log.Println("WARN: " + fmt.Sprint(args...))
}

func Warningf(format string, args ...interface{}) {
	log.Println("WARN: " + fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	log.Println("ERROR: " + fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	log.Println("ERROR: " + fmt.Sprintf(format, args...))
}
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

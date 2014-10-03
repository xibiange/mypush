package log

import (
	"fmt"
)

func Error(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

func Debug(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

func Info(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

func Warn(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

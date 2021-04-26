package logger

import "fmt"

func Info(msg string) {
	fmt.Printf(fmt.Sprintf("%s\n", msg))
}

func Error(msg string, err error) {
	fmt.Printf(fmt.Sprintf("%s :: %+v\n", msg, err))
}
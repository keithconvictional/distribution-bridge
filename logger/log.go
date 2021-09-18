package logger

import "fmt"

func Info(requestID string, domain string, msg string) {
	fmt.Printf(fmt.Sprintf("%s\n", msg))
}

func Error(requestID string, domain string, msg string, err error) {
	fmt.Printf(fmt.Sprintf("%s :: %+v\n", msg, err))
}
package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(message string)
	Error(message string)
}

type StdoutLogger struct{}

func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{}
}

func (log *StdoutLogger) Info(message string) {
	fmt.Printf("INFO (%v) : %v \n", time.Now().Format("2006-01-02 15:04:05.0000"), message)
}

func (log *StdoutLogger) Error(message string) {
	fmt.Printf("ERROR (%v) : %v \n", time.Now().Format("2006-01-02 15:04:05.0000"), message)
}

package logger

import (
	"log"
	"os"
)

const (
	LSilent = 0
	LInfo   = 1
	LDebug  = 2
)

func NewDebugLogger() *log.Logger {
	return log.New(os.Stderr, "\033[34m[DEBUG] \033[0m", log.Lmsgprefix)
}

func NewErrorLogger() *log.Logger {
	return log.New(os.Stderr, "\033[31m[ERROR] \033[0m", log.Lmsgprefix)
}

package logger

import (
	"log"
	"runtime"
)

func ServerLog(err error, message string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	if err != nil {
		log.Printf("[ERROR] %s:%d | err=%v | msg=%s", file, line, err, message)
	} else {
		log.Printf("[ERROR] %s:%d | msg=%s", file, line, message)
	}
}

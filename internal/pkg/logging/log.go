package logging

import (
	"log"
	"os"
	"sync"
)

const logFile string = "gotify.log"

var once sync.Once
var logger *log.Logger

func GetLoggerInstance() *log.Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *log.Logger {
	file, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return log.New(file, "", log.Ltime|log.Lshortfile|log.Ldate)
}

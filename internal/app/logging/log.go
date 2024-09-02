package logging

import (
	"log"
	"os"
)

const logFile string = "gotify_log.txt"

var (
	GotifyLogger *log.Logger
)

func init() {
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	GotifyLogger = log.New(file, "", log.Lshortfile|log.Ldate|log.Ltime)
}

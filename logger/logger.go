package logger

import (
	"log"
	"os"
)

func Init(logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LogInfo(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

func LogError(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

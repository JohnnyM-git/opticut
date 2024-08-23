package logger

import (
	"log"
	"os"
)

func LogError(msg string) {
	file, err := os.OpenFile("appErrors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer file.Close()
	logmsg := log.New(file, "Error: ", log.Ldate|log.Ltime)
	logmsg.Println(msg)
}
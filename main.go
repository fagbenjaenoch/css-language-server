package main

import (
	"log"
	"os"
)

func main() {
	log := getLogger("/home/enoch/dev/css-language-server/log.txt")
	log.Println("server started")
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	return log.New(logFile, "[css-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}

package main

import (
	"bufio"
	"log"
	"os"

	"github.com/fagbenjaenoch/css-language-server/rpc"
)

func main() {
	log := getLogger("/home/enoch/dev/css-language-server/log.txt")
	log.Println("server started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		rpcRequest := scanner.Text()
		log.Printf("received a message from the client of %d bytes", len(rpcRequest))
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	return log.New(logFile, "[css-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}

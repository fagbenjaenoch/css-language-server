package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/fagbenjaenoch/css-language-server/lsp"
	"github.com/fagbenjaenoch/css-language-server/rpc"
)

func main() {
	log := getLogger("/home/enoch/dev/css-language-server/log.txt")
	log.Println("server started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		rpcRequest := scanner.Bytes()
		log.Printf("received a message from the client of %d bytes", len(rpcRequest))

		method, body, err := rpc.DecodeMessage(rpcRequest)
		if err != nil {
			log.Printf("could not parse request body: %s", err)
			continue
		}
		log.Printf("received a request of %s method", method)

		switch method {
		case rpc.MethodInitialize:
			{
				log.Println("initializing server")
				var request lsp.InitializeParams
				if err := json.Unmarshal(body, &request); err != nil {
					log.Printf("could not unmarshal request body: %s", err)
					continue
				}

				log.Printf("connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

				initializeResponse := lsp.NewInitializeResponse(request.ID)
				response := rpc.EncodeMessage(initializeResponse)
				os.Stdout.Write([]byte(response))
			}
		}
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	return log.New(logFile, "[css-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}

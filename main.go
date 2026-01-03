package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/fagbenjaenoch/css-language-server/lsp"
	"github.com/fagbenjaenoch/css-language-server/rpc"
)

func main() {
	logger := getLogger("/home/enoch/dev/css-language-server/log.txt")
	logger.Println("server started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	writer := os.Stdout

	for scanner.Scan() {
		rpcRequest := scanner.Bytes()
		logger.Printf("received a message from the client of %d bytes", len(rpcRequest))

		method, body, err := rpc.DecodeMessage(rpcRequest)
		if err != nil {
			logger.Printf("could not parse request body: %s", err)
			continue
		}

		handleRequest(logger, writer, method, body)
	}
}

func handleRequest(logger *log.Logger, writer io.Writer, method string, body []byte) {
	logger.Printf("recieved method %s from client", method)

	switch method {
	case rpc.MethodInitialize:
		{
			log.Println("initializing server")
			var request lsp.InitializeParams
			if err := json.Unmarshal(body, &request); err != nil {
				log.Printf("could not unmarshal request body: %s", err)
			}

			log.Printf("connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

			initializeResponse := lsp.NewInitializeResponse(request.ID)
			response := rpc.EncodeMessage(initializeResponse)
			writer.Write([]byte(response))
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

package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/fagbenjaenoch/css-language-server/analysis"
	"github.com/fagbenjaenoch/css-language-server/lsp"
	"github.com/fagbenjaenoch/css-language-server/rpc"
)

func main() {
	logger := getLogger("/home/enoch/dev/css-language-server/log.txt")
	logger.Println("server started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	writer := os.Stdout

	state := analysis.NewState()

	for scanner.Scan() {
		rpcRequest := scanner.Bytes()
		logger.Printf("received a message from the client of %d bytes", len(rpcRequest))

		method, body, err := rpc.DecodeMessage(rpcRequest)
		if err != nil {
			logger.Printf("could not parse request body: %s", err)
			continue
		}

		handleRequest(logger, writer, state, method, body)
	}
}

func handleRequest(logger *log.Logger, writer io.Writer, state *analysis.State, method string, body []byte) {
	logger.Printf("recieved method '%s' from client", method)

	switch method {
	case rpc.MethodInitialize:
		logger.Println("initializing server")
		var request lsp.InitializeRequest
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Printf("could not unmarshal request body: %s", err)
		}

		logger.Printf("connected to %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		initializeResponse := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, initializeResponse)

	case rpc.MethodTextDocumentDidOpen:
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Printf("could not parse notification from textDocument/didOpen")
		}

		diagnostics := state.OpenDocument(request.Params.TextDocument.Uri, request.Params.TextDocument.Text)
		if len(diagnostics) == 0 {
			return
		}
		response := lsp.PublishDiagnosticsParams{
			Uri:         request.Params.TextDocument.Uri,
			Diagnostics: diagnostics,
		}
		writeResponse(writer, response)

	case rpc.MethodTextDocumentDidChange:
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Printf("could not parse notification from textDocument/didChange")
			return
		}

		diagnostics := state.UpdateDocument(request.Params.TextDocument.Uri, request.Params.ContentChanges[0].Text)
		if len(diagnostics) == 0 {
			return
		}
		response := lsp.PublishDiagnosticsParams{
			Uri:         request.Params.TextDocument.Uri,
			Diagnostics: diagnostics,
		}
		writeResponse(writer, response)
	case rpc.MethodTextDocumentCompletion:
		var request lsp.CompletionRequest
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Printf("could not parse %s request body", rpc.MethodTextDocumentCompletion)
			return
		}

		// response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.Uri)
		// writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	response := rpc.EncodeMessage(msg)
	writer.Write([]byte(response))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	return log.New(logFile, "[css-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}

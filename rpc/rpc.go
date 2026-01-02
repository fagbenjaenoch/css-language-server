package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	jsonPayload, err := json.Marshal(msg)
	if err != nil {
		panic("could not serialize json payload")
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(jsonPayload), jsonPayload)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(data []byte) (string, []byte, error) {
	header, body, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("incomplete request body")
	}

	bodyLengthBytes := header[len("Content-Length: "):]
	bodyLength, err := strconv.Atoi(string(bodyLengthBytes))
	if err != nil {
		return "", nil, errors.New("could not parse content length from header")
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(body, &baseMessage); err != nil {
		return "", nil, errors.New("could not parse message body")
	}

	return baseMessage.Method, body[:bodyLength], nil
}

func DecodeMessage(data []byte) (any, error) {
	return nil, nil
}

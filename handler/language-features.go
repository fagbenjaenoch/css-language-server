package handler

import (
	"github.com/fagbenjaenoch/css-language-server/mappers"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	var completionItems []protocol.CompletionItem

	for word, mapper := range mappers.CSSMappers {
		mapperCopy := mapper

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &mapperCopy,
			InsertText: &mapperCopy,
		})
	}

	return completionItems, nil
}

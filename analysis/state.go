package analysis

import (
	"strings"

	"github.com/fagbenjaenoch/css-language-server/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() *State {
	return &State{
		Documents: map[string]string{},
	}
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos. Don't forget to like & subscribe to streamers using it :)",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	return getDiagnosticsForFile(text)
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "hello") {
			idx := strings.Index(line, "hello")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    lsp.LineRange(row, idx, idx+len("hello")),
				Severity: 2,
				Source:   "Wisdom from a thousand years",
				Message:  "You should use 'hello world' instead",
			})
		}
	}

	return diagnostics

}

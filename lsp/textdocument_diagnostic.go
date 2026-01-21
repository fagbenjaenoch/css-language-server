package lsp

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}

type PublishDiagnosticsParams struct {
	Uri         string       `json:"uri"`
	Version     *int         `json:"version"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

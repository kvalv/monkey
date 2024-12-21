package msg

type Body interface {
	Bytes() []byte
	MethodName() string
}

// we've embedded the JsonRPC field into each struct to avoid
// too much struct embeddings
type (
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	// everything is optional so we omit for now
	TextDocument struct {
		URI DocumentURI `json:"uri"`
	}

	DocumentURI      = string
	TextDocumentItem struct {
		URI DocumentURI `json:"uri"`
		// e.g. "c", "go", "python", "dockerfile", ...
		LanguageID string `json:"languageId"`
		Version    int    `json:"version"`
		Text       string `json:"text"`
	}

	TextDocumentIdentifier struct {
		Uri DocumentURI `json:"uri"`
	}
	VersionedTextDocumentIdentifier struct {
		TextDocumentIdentifier
		Version int `json:"version"`
	}
	Position struct {
		Line      int `json:"line"`
		Character int `json:"character"`
	}
	Range struct {
		Start Position `json:"start"`
		End   Position `json:"end"`
	}
	TextDocumentPositionParams struct {
		TextDocument TextDocumentIdentifier `json:"textDocument"`
		Position     Position               `json:"position"`
	}
	TextDocumentContentChangeEvent struct {
		// Range 	 Range `json:"range"`
		// RangeLength int `json:"rangeLength"`

		// we only operate on full text changes so the above variant
		// does not apply
		Text string `json:"text"`
	}
	// "markdown", "plaintext"
	MarkupKind    = string
	MarkupContent struct {
		// "markdown", "plaintext"
		Kind MarkupKind `json:"kind"`

		// The content itself
		Value string `json:"value"`
	}
)

type (
	CompletionRequest = Request[CompletionParams]
)

// requests
type (
	CompletionParams struct {
		TextDocument `json:"textDocument"`
		Position     `json:"position"`
	}
)

const (
	TEXT_DOCUMENT_SYNC_KIND_NONE        = 0
	TEXT_DOCUMENT_SYNC_KIND_FULL        = 1
	TEXT_DOCUMENT_SYNC_KIND_INCREMENTAL = 2
)

// Structs related to the initialize request

func NewCompletionRequest(id int, uri string, pos Position) *Request[CompletionParams] {
	return &Request[CompletionParams]{
		Id:     id,
		Method: METHOD_REQUEST_COMPLETION,
		Params: CompletionParams{
			TextDocument: TextDocument{URI: uri},
			Position:     pos,
		},
	}
}

package msg

import "encoding/json"

type (
	Response[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  T      `json:"result"`
		Error   *Error `json:"error,omitempty"`
	}
)

func (r *Response[T]) MarshalJSON() ([]byte, error) {
	// when marshalling we'll auto-populate the JsonRPC field
	type alias Response[T]
	return json.Marshal(&struct {
		JsonRPC string `json:"jsonrpc"`
		*alias
	}{
		JsonRPC: "2.0",
		alias:   (*alias)(r),
	})
}
func (r *Response[T]) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}
func (r *Response[T]) MethodName() string { return "" }

// no methods in response, so don't need that

type (
	ServerCapabilities struct {
		// for convenience, should be TEXT_DOCUMENT_SYNC_KIND_FULL
		TextDocumentSync int  `json:"textDocumentSync"`
		HoverProvider    bool `json:"hoverProvider"`
	}
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	InitializeResult struct {
		// required
		Capabilities ServerCapabilities `json:"capabilities"`
		// optional
		ServerInfo `json:"serverInfo"`
	}

	ResponseHover = Response[Hover]
	Hover         struct {
		Contents MarkupContent `json:"contents"`
		Range    *Range        `json:"range,omitempty"`
	}
)

func NewInitializeResult(id int) *Response[InitializeResult] {
	return &Response[InitializeResult]{
		Id: id,
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TEXT_DOCUMENT_SYNC_KIND_FULL,
				HoverProvider:    true,
			},
			ServerInfo: ServerInfo{
				Name:    "monkey-lsp",
				Version: "1.0.0",
			},
		},
		Error: nil,
	}
}

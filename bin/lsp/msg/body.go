package msg

import "encoding/json"

type Body interface {
	Bytes() []byte
}

// we've embedded the JsonRPC field into each struct to avoid
// too much struct embeddings
type (
	Request[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int    `json:"id"` // actually string | number, but we just use string
		Method  string `json:"method"`
		Params  T      `json:"params"`
	}
	Response[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  T      `json:"result"`
		Error   *Error `json:"error,omitempty"`
	}
	Notification[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		// the method to be invoked
		Method string `json:"method"`
		// The notification's params
		Params T `json:"Params"`
	}
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	// everything is optional so we omit for now
	TextDocument struct {
		URI string `json:"uri"`
	}
	Position struct {
		Line      int `json:"line"`
		Character int `json:"character"`
	}
)

type (
	InitializeRequest       = Request[InitializeParams]
	CompletionRequest       = Request[CompletionParams]
	DidSaveTextNotification = Notification[DidSaveTextDocumentParams]
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

func (r *Request[T]) MarshalJSON() ([]byte, error) {
	// when marshalling we'll auto-populate the JsonRPC field
	type alias Request[T]
	return json.Marshal(&struct {
		JsonRPC string `json:"jsonrpc"`
		*alias
	}{
		JsonRPC: "2.0",
		alias:   (*alias)(r),
	})
}

func (n *Notification[T]) MarshalJSON() ([]byte, error) {
	// when marshalling we'll auto-populate the JsonRPC field
	type alias Notification[T]
	return json.Marshal(&struct {
		JsonRPC string `json:"jsonrpc"`
		*alias
	}{
		JsonRPC: "2.0",
		alias:   (*alias)(n),
	})
}

func (r *Request[T]) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}
func (r *Response[T]) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}
func (n *Notification[T]) Bytes() []byte {
	b, _ := json.Marshal(n)
	return b
}

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

type (
	InitializeParams struct {
		ProcessId  int `json:"processId"`
		ClientInfo `json:"clientInfo"`
	}
	ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	InitializeResult struct {
		// required
		Capabilities ServerCapabilities `json:"capabilities"`
		// optional
		ServerInfo `json:"serverInfo"`
	}
	ServerCapabilities struct {
		// for convenience, should be TEXT_DOCUMENT_SYNC_KIND_FULL
		TextDocumentSync int `json:"textDocumentSync"`
	}
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)

type (
	DidSaveTextDocumentParams struct {
		TextDocument `json:"textDocument"`
		Text         string `json:"text,omitempty"`
	}
)

func NewCompletionRequest(id int, uri string, pos Position) *Request[CompletionParams] {
	return &Request[CompletionParams]{
		Id:     id,
		Method: METHOD_COMPLETION,
		Params: CompletionParams{
			TextDocument: TextDocument{URI: uri},
			Position:     pos,
		},
	}
}

func NewInitializeResult(id int) *Response[InitializeResult] {
	return &Response[InitializeResult]{
		Id: id,
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TEXT_DOCUMENT_SYNC_KIND_FULL,
			},
			ServerInfo: ServerInfo{
				Name:    "monkey-lsp",
				Version: "1.0.0",
			},
		},
		Error: nil,
	}
}

package msg

import "encoding/json"

type (
	Request[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		Id      int    `json:"id"` // actually string | number, but we just use string
		Method  string `json:"method"`
		Params  T      `json:"params"`
	}
)

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
func (r *Request[T]) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}
func (r *Request[T]) MethodName() string { return r.Method }

const (
	METHOD_REQUEST_INITIALIZE = "initialize"
	METHOD_REQUEST_HOVER      = "textDocument/hover"
	METHOD_REQUEST_COMPLETION = "textDocument/completion"
)

type (
	HoverParams struct {
		TextDocumentPositionParams `json:"textDocument"`
		Position                   Position `json:"position"`
	}
	RequestHover = Request[HoverParams]

	ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	InitializeParams struct {
		ProcessId  int `json:"processId"`
		ClientInfo `json:"clientInfo"`
	}
	RequestInitialize = Request[InitializeParams]
)

func getRequestBody(method string) Body {
	switch method {
	case METHOD_REQUEST_INITIALIZE:
		return &RequestInitialize{}
	case METHOD_REQUEST_HOVER:
		return &RequestHover{}
	case METHOD_REQUEST_COMPLETION:
		return &CompletionRequest{}
	default:
		return nil
	}
}

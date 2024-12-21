package msg

import "encoding/json"

type (
	Notification[T any] struct {
		JsonRPC string `json:"jsonrpc"`
		// the method to be invoked
		Method string `json:"method"`
		// The notification's params
		Params T `json:"Params"`
	}
)

func (n *Notification[T]) Bytes() []byte {
	b, _ := json.Marshal(n)
	return b
}
func (r *Notification[T]) MethodName() string { return r.Method }
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

const (
	METHOD_NOTIFICATION_DID_SAVE    = "textDocument/didSave"
	METHOD_NOTIFICATION_DID_OPEN    = "textDocument/didOpen"
	METHOD_NOTIFICATION_DID_CLOSE   = "textDocument/didClose"
	METHOD_NOTIFICATION_DID_CHANGE  = "textDocument/didChange"
	METHOD_NOTIFICATION_HOVER       = "textDocument/hover"
	METHOD_NOTIFICATION_INITIALIZED = "initialized"
)

type (
	DidSaveNotification       = Notification[DidSaveTextDocumentParams]
	DidSaveTextDocumentParams struct {
		TextDocument TextDocumentItem `json:"textDocument"`

		// Depends on includeText on server capabilities
		Text string `json:"text,omitempty"`
	}

	DidOpenNotification       = Notification[DidOpenTextDocumentParams]
	DidOpenTextDocumentParams struct {
		TextDocument TextDocumentItem `json:"textDocument"`
	}

	DidCloseNotification       = Notification[DidCloseTextDocumentParams]
	DidCloseTextDocumentParams struct {
		TextDocument TextDocumentItem `json:"textDocument"`
	}

	DidChangeNotification       = Notification[DidChangeTextDocumentParams]
	DidChangeTextDocumentParams struct {
		TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
		ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
	}

	InitializedNotification = Notification[InitializeParams]
	InitializedParams       struct{}
)

func getNotificationBody(method string) Body {
	switch method {
	case METHOD_NOTIFICATION_DID_SAVE:
		return &DidSaveNotification{}
	case METHOD_NOTIFICATION_DID_OPEN:
		return &DidOpenNotification{}
	case METHOD_NOTIFICATION_DID_CLOSE:
		return &DidCloseNotification{}
	case METHOD_NOTIFICATION_DID_CHANGE:
		return &DidChangeNotification{}
	case METHOD_NOTIFICATION_INITIALIZED:
		return &InitializedNotification{}
	default:
		return nil
	}
}

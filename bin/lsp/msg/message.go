package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

const (
	SEP       = "\r\n"
	DOUBLESEP = "\r\n\r\n"
)

// A Message contains a header and a body and represents a LSP message
type Message struct {
	*Header
	Body
}

func New(msg Body) *Message {
	h := Header{ContentLength: len(msg.Bytes())}
	return &Message{&h, msg}
}

// Returns an appropriate method based on the bytes received
func FromBytes(b []byte) (*Message, error) {
	parts := bytes.SplitN(b, []byte(DOUBLESEP), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("FromBytes: expected a header and content part")
	}
	headerBytes := parts[0]
	contentBytes := parts[1]

	header, err := HeaderFromBytes(headerBytes)
	if err != nil {
		return nil, err
	}

	var data struct {
		Method string `json:"method"`
		Result map[string]any
	}
	log.Printf("content=%s", string(contentBytes))

	if err := json.Unmarshal(contentBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}

	var body Body

	// if method is provided, we'll set the body to one of the supported methods
	if data.Method != "" {
		switch data.Method {
		case METHOD_INITIALIZE:
			body = &InitializeRequest{}
		case METHOD_COMPLETION:
			body = &CompletionRequest{}
		case METHOD_DID_SAVE:
			body = &DidSaveTextNotification{}
		}
	}
	if body == nil {
		return nil, fmt.Errorf("not implemented for method %q cl=%d", data.Method, header.ContentLength)
	}
	if err := json.Unmarshal(contentBytes, body); err != nil {
		return nil, err
	}
	return &Message{header, body}, nil
}

func (r *Message) String() string {
	return fmt.Sprintf("%s\r\n%s", r.Header.Bytes(), r.Body.Bytes())
}

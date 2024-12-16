package msg_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kvalv/monkey/bin/lsp/msg"
)

func TestMessageFromBytes(t *testing.T) {
	t.Run("requests", func(t *testing.T) {
		t.Run("request", func(t *testing.T) {
			body := `{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"processId":600660,"clientInfo":{"name":"Neovim","version":"0.11.0-dev+gf9dd682621"}}}`
			want := &msg.InitializeRequest{
				JsonRPC: "2.0",
				Method:  "initialize",
				Id:      1,
				Params: msg.InitializeParams{
					ProcessId: 600660,
					ClientInfo: msg.ClientInfo{
						Name:    "Neovim",
						Version: "0.11.0-dev+gf9dd682621",
					},
				},
			}
			parsed := mustMessageFromString(t, body)
			got, ok := parsed.Body.(*msg.InitializeRequest)
			if !ok {
				t.Fatalf("expected *msg.InitializeRequest, got %T", parsed.Body)
			}
			// we'll just match by string comparison for now
			if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", want) {
				t.Fatalf("mismatch\n got = %+v\nwant = %+v", got, want)
			}
		})
	})
	t.Run("notifications", func(t *testing.T) {
		t.Run("didSave", func(t *testing.T) {
			body := `{ "jsonrpc": "2.0", "method": "textDocument/didSave", "params": { "textDocument": { "uri": "file:///path/to/file.go" }, "text": "package main\n\nfunc main() {\n\tprintln(\"Hello, LSP\")\n}" } } `
			want := &msg.DidSaveTextNotification{
				JsonRPC: "2.0",
				Method:  "textDocument/didSave",
				Params: msg.DidSaveTextDocumentParams{
					TextDocument: msg.TextDocument{
						URI: "file:///path/to/file.go",
					},
					Text: "package main\n\nfunc main() {\n\tprintln(\"Hello, LSP\")\n}",
				},
			}
			parsed := mustMessageFromString(t, body)
			got, ok := parsed.Body.(*msg.DidSaveTextNotification)
			if !ok {
				t.Fatalf("expected *msg.DidSaveTextNotification, got %T", parsed.Body)
			}
			if uri, want := got.Params.URI, "file:///path/to/file.go"; uri != want {
				t.Fatalf("uri mismatch; wanted %q but got %q", want, uri)
			}
			// we'll just match by string comparison for now
			if fmt.Sprintf("%+v", got) != fmt.Sprintf("%+v", want) {
				t.Fatalf("mismatch\n got = %+v\nwant = %+v", got, want)
			}
		})
	})

}

func TestMessageToBytes(t *testing.T) {
	cases := []struct {
		message *msg.Message
		want    string
	}{
		{
			message: msg.New(msg.NewCompletionRequest(1, "foo.txt", msg.Position{Line: 123, Character: 234})),
			want:    `content-length:143\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"textDocument/completion","params":{"textDocument":{"uri":"foo.txt"},"position":{"line":123,"character":234}}}`,
		},
		{
			message: msg.New(msg.NewInitializeResult(1)),
			want:    `content-length:126\r\n\r\n{"jsonrpc":"2.0","id":1,"result":{"capabilities":{"textDocumentSync":1},"serverInfo":{"name":"monkey-lsp","version":"1.0.0"}}}`,
		},
		{
			message: msg.New(msg.NewInitializeResult(1)),
			want:    `content-length:126\r\n\r\n{"jsonrpc":"2.0","id":1,"result":{"capabilities":{"textDocumentSync":1},"serverInfo":{"name":"monkey-lsp","version":"1.0.0"}}}`,
		},
	}
	for _, tc := range cases {
		got := tc.message.String()
		want := strings.ReplaceAll(tc.want, `\r\n`, "\r\n")
		if got != want {
			t.Fatalf("string mismatch\ngot =%q\nwant=%q", got, want)
		}
	}
}

func mustMessageFromString(t *testing.T, s string) *msg.Message {
	t.Helper()
	contentLength := len(s)
	rpc := fmt.Sprintf("content-length:%d\r\n\r\n%s", contentLength, s)
	parsed, err := msg.FromBytes([]byte(rpc))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return parsed
}

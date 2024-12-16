package scanner_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/kvalv/monkey/bin/lsp/msg"
	"github.com/kvalv/monkey/bin/lsp/scanner"
)

func TestScan(t *testing.T) {
	t.Run("two completion requests", func(t *testing.T) {
		msg1 := msg.New(msg.NewCompletionRequest(1, "foo.txt", msg.Position{Line: 123, Character: 234})).String()
		msg2 := msg.New(msg.NewCompletionRequest(2, "foo.txt", msg.Position{Line: 345, Character: 567})).String()
		r, w := io.Pipe()
		go func() {
			b := &bytes.Buffer{}
			fmt.Fprintf(b, msg1)
			fmt.Fprintf(b, msg2)
			fmt.Fprintf(b, "hello world") // just some arbitrary data, not yet large enough for a message
			w.Write(b.Bytes())
		}()

		sc := scanner.New(r)
		parsed1 := expectNextMessage[*msg.CompletionRequest](t, sc)
		if got := parsed1.Params.Line; got != 123 {
			t.Fatalf("expected line 123, got %d", got)
		}
		parsed2 := expectNextMessage[*msg.CompletionRequest](t, sc)
		if got := parsed2.Params.Line; got != 345 {
			t.Fatalf("expected line 123, got %d", got)
		}
		expectNoScan(t, sc)
	})

	t.Run("initialize", func(t *testing.T) {
		rpc := mustMessageFromString(t, `{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"processId":600660,"clientInfo":{"name":"Neovim","version":"0.11.0-dev+gf9dd682621"}}}`)
		r, w := io.Pipe()
		go func() {
			fmt.Fprintf(w, rpc.String())
		}()
		sc := scanner.New(r)
		req := expectNextMessage[*msg.InitializeRequest](t, sc)
		if req.Params.ProcessId != 600660 {
			t.Fatalf("mismatch: got=%d", req.Params.ProcessId)
		}
	})
}

func expectNextMessage[T msg.Body](t *testing.T, sc *scanner.Scanner) T {
	t.Helper()
	sc.Scan()
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %s", err)
	}
	got := sc.Next()
	if got == nil {
		t.Fatalf("next is nil")
	}
	parsed, ok := got.Body.(T)
	if !ok {
		var emp T
		t.Fatalf("expected message to be of type %T got=%T", emp, got)
	}
	return parsed
}

func expectNoScan(t *testing.T, sc *scanner.Scanner) {
	t.Helper()
	if sc.Scan() {
		t.Fatalf("expected no next message, but scan succeeded")
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

package msg

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Header represents the header part of an RPC message
type Header struct {
	ContentLength int
	rest          map[string]any
}

var (
	// separator between two headers
	headerSep = []byte("\r\n")

	// separator between key and value; same as HTTP headers
	kvSep = []byte(":")
)

const (
	HEADER_CONTENT_LENGTH = "content-length"
)

func normalizeKey(s string) string { return strings.Trim(strings.ToLower(s), " ") }

func parseHeader(b []byte) (string, int, error) {
	parts := bytes.SplitN(b, kvSep, 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("failed to split header into key and value (%q)", b)
	}

	key := normalizeKey(string(parts[0]))

	// for now we're assuming all headers are integers
	value, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse header value to int")
	}

	return key, value, nil
}

func HeaderFromBytes(b []byte) (*Header, error) {
	header := &Header{rest: make(map[string]any)}
	for _, pair := range bytes.Split(b, headerSep) {
		if len(pair) == 0 {
			// trailing -- we're done
			return header, nil
		}
		key, value, err := parseHeader(pair)
		if err != nil {
			return nil, err
		}
		switch key {
		case HEADER_CONTENT_LENGTH:
			header.ContentLength = value
		default:
			return nil, fmt.Errorf("unknown header: %q", key)
		}
	}
	return header, nil
}

func (h *Header) Bytes() []byte {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "content-length:%d\r\n", h.ContentLength)
	return b.Bytes()
}

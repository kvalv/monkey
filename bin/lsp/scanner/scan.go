package scanner

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/kvalv/monkey/bin/lsp/msg"
)

// A Scanner reads until it finds a message that matches the lsp protocol format.
type Scanner struct {
	sc   *bufio.Scanner
	err  error
	next *msg.Message
}

var (
	// separator between two header values
	sep = []byte("\r\n")

	// separator between header and body of our message
	headerBodySep = append(sep, sep...)
)

func New(r io.Reader) *Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		// we want to find the header, because it contains the content-length of the body.
		i := bytes.Index(data, headerBodySep)
		if i == 0 {
			// Request more data.
			return 0, nil, nil
		}
		if i == -1 {
			// We don't have a full header yet.
			log.Printf("data: %d, %s", i, string(data))
			return 0, nil, nil
		}

		// We include 2 additional bytes (\r\n) for the last header key-value pair
		headerLength := i + len(sep)

		log.Printf("data: %d, %s", i, string(data[:headerLength]))
		header, err := msg.HeaderFromBytes(data[:headerLength])
		if err != nil {
			return 0, nil, err
		}

		// total message size = size(header) + size(\r\n) + size(body)
		if msgSize := headerLength + len(sep) + header.ContentLength; msgSize <= len(data) {
			return msgSize, data[:msgSize], nil
		} else {
			// Request more data; full message hasn't arrived yet.
			return 0, nil, nil
		}
	})
	return &Scanner{sc, nil, nil}

}
func (s *Scanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.sc.Err()
}

func (s *Scanner) Scan() bool {
	ok := s.sc.Scan()
	if !ok {
		s.err = s.sc.Err()
		return false
	}
	s.next, s.err = msg.FromBytes(s.sc.Bytes())
	return s.err == nil
}

func (s *Scanner) Next() *msg.Message {
	if s.err != nil {
		return nil
	}
	return s.next
}

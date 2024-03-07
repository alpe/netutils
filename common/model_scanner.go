package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const endOfStream = '\000'

// json scanner. does not support concurrent calls.
type scanner struct {
	buf        []byte
	err        error
	cursor     int
	inQuotes   bool
	bufferSize int
}

func (s *scanner) scan(in io.Reader) (string, error) {
	s.buf = make([]byte, s.bufferSize)
	n, err := in.Read(s.buf)
	if err != nil {
		if err != io.EOF {
			return "", fmt.Errorf("read input stream: %w", err)
		}
	}
	s.buf = append(s.buf[0:n], endOfStream)

	if c := s.skipWhiteSpaces(); c != '{' {
		return "", fmt.Errorf("invalid character: %c", c)
	}
	model := captureAttributeValue(s, "model")
	return model, s.err
}

func captureAttributeValue(s *scanner, name string) string {
	switch c := s.skipWhiteSpaces(); c {
	case '"':
		s.inQuotes = true
		if key := captureStringValue(s); key != name {
			// not model as first attribute
			return ""
		}
		if c := s.skipWhiteSpaces(); c != ':' {
			if c != endOfStream {
				s.err = fmt.Errorf("expected ':' but got character: %c", c)
			}
			return ""
		}

		if c := s.skipWhiteSpaces(); c != '"' {
			if c != endOfStream {
				s.err = fmt.Errorf("not string type: %c", c)
			}
			return ""
		}
		s.inQuotes = true
		return captureStringValue(s)
	case endOfStream:
		return ""
	default:
		s.err = fmt.Errorf("invalid character: %c", c)
		return ""
	}
}

func captureStringValue(s *scanner) string {
	var captured string
	for {
		switch c := s.readNext(); c {
		case '"':
			if s.inQuotes {
				s.inQuotes = false
				return captured
			}
		case '\\':
			return captured + captureEscapedString(s)
		case endOfStream:
			return ""
		default:
			captured += string(c)
		}
	}
}

func captureEscapedString(s *scanner) string {
	for {
		switch c := s.readNext(); c {
		case '"', '\\':
			return string(c) + captureStringValue(s)
		case 'b', 'f', 'n', 'r', 'v', 't': // ignore backspace, page break, new line,carriage return, vertical tab, horizontal tab
			return captureStringValue(s)
		case 'u':

			s.err = errors.New("can not handle unicode")
			return ""
		case endOfStream:
			return ""
		}
	}
}

func (s *scanner) readNext() byte {
	defer func() { s.cursor++ }()
	if len(s.buf) == s.cursor {
		return endOfStream
	}
	return s.buf[s.cursor]
}

func (s *scanner) skipWhiteSpaces() byte {
	c := s.readNext()
	switch c {
	case ' ', '\n', '\t', '\r', '\v':
		return s.skipWhiteSpaces()
	default:
		return c
	}
}

func PeekModel(in io.Reader, bufferSize int) (string, io.ReadCloser, error) {
	buff := bytes.NewBuffer([]byte{})
	s := &scanner{bufferSize: bufferSize}

	model, err := s.scan(io.TeeReader(in, buff))
	if err != nil {
		return "", nil, err
	}
	return model, MultiReadCloser(buff, in), nil
}

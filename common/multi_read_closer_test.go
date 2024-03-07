package common

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiReadCloser(t *testing.T) {
	captureCloseMock := captureClose("bar")
	specs := map[string]struct {
		readers       []io.Reader
		expContent    string
		expMockClosed bool
	}{
		"single": {
			readers:    []io.Reader{strings.NewReader("foo")},
			expContent: "foo",
		},
		"multiple": {
			readers:    []io.Reader{strings.NewReader("foo"), strings.NewReader(" "), strings.NewReader("bar")},
			expContent: "foo bar",
		},
		"multiple with closer": {
			readers:       []io.Reader{strings.NewReader("foo"), strings.NewReader(" "), captureCloseMock},
			expContent:    "foo bar",
			expMockClosed: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			captureCloseMock.Reset()
			m := MultiReadCloser(spec.readers...)

			b, err := io.ReadAll(m)
			require.NoError(t, err)

			require.Equal(t, spec.expContent, string(b))

			err = m.Close()
			require.NoError(t, err)
			assert.Equal(t, spec.expMockClosed, captureCloseMock.closed)
		})
	}
}

var _ io.Closer = (*captureReadCloser)(nil)

type captureReadCloser struct {
	io.Reader
	closed bool
}

func captureClose(c string) *captureReadCloser {
	return &captureReadCloser{Reader: strings.NewReader(c)}
}

func (c *captureReadCloser) Close() error {
	c.closed = true
	return nil
}

func (c *captureReadCloser) Reset() {
	c.closed = false
}

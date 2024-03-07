package common

import "io"

var _ io.WriterTo = (*multiReadCloser)(nil)

type multiReadCloser struct {
	closer      []io.Closer
	multiReader io.Reader
}

// MultiReadCloser returns a ReadCloser that's the logical concatenation of
// the provided input readers. It uses the stdlib io.MultiReader under the hood but adds the Close function
// for all supported readers.
func MultiReadCloser(readers ...io.Reader) io.ReadCloser {
	closers := make([]io.Closer, 0, len(readers))
	for _, r := range readers {
		if o, ok := r.(io.Closer); ok {
			closers = append(closers, o)
		}
	}
	return &multiReadCloser{
		multiReader: io.MultiReader(readers...),
		closer:      closers,
	}
}

func (m *multiReadCloser) Read(p []byte) (n int, err error) {
	return m.multiReader.Read(p)
}

func (m *multiReadCloser) WriteTo(w io.Writer) (n int64, err error) {
	return m.multiReader.(io.WriterTo).WriteTo(w) //nolint: forcetypeassert
}

func (m *multiReadCloser) Close() error {
	for _, c := range m.closer {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

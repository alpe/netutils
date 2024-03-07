package common

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelParser(t *testing.T) {
	specs := map[string]struct {
		source   []byte
		expModel any
		expErr   bool
	}{
		"is first attribute": {
			source:   []byte(`{"model":"my-model", "post-data":"bar"}`),
			expModel: "my-model",
		},
		"not first attribute": {
			source:   []byte(`{"pre-data":"foo","model":"my-model", "post-data":"bar"}`),
			expModel: "",
		},
		"no model attribute json": {
			source:   []byte(`{"foo": "model"`),
			expModel: "",
		},
		"escaped text value": {
			source:   []byte(`{"model":"\"my-model"`),
			expModel: `"my-model`,
		},
		"trim invisible escaped text value": {
			source:   []byte(`{ "model": "\nmy-model\t"`),
			expModel: "my-model",
		},
		"escaped text value: backslash": {
			source:   []byte(`{ "model": "\\my-model"`),
			expModel: `\my-model`,
		},
		"escaped text value: unicode": {
			source: []byte(`{ "model": "\u00C4my-model"`),
			expErr: true,
		},
		"invalid json": {
			source: []byte(`{"model" "foo"`),
			expErr: true,
		},
		"unsupported type: slice": {
			source: []byte(`{"model": []}`),
			expErr: true,
		},
		"unsupported type: number": {
			source: []byte(`{"model": 1}`),
			expErr: true,
		},
		"unsupported type: null": {
			source: []byte(`{"model": null}`),
			expErr: true,
		},
		"unsupported type: object": {
			source: []byte(`{"model": {}}`),
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			gotModel, ioStream, gotErr := PeekModel(bytes.NewReader(spec.source), 128)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, spec.expModel, gotModel)
			// and common intact
			allData, err := io.ReadAll(ioStream)
			require.NoError(t, err)
			assert.Equal(t, spec.source, allData)
		})
	}
}

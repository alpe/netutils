package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type jsonType struct {
	Model string `json:"model"`
}

func BenchmarkDecoding(b *testing.B) {
	specs := map[string]struct {
		src string
	}{
		"minimal": {
			src: `{"model":"my-model" }`,
		},
		"small": {
			src: fmt.Sprintf(`{"model":"my-model", "other": %q }`, strings.Repeat("a", 100)),
		},
		"medium": {
			src: fmt.Sprintf(`{"model":"my-model", "other": %q }`, strings.Repeat("a", 50_000)),
		},
		"big": {
			src: fmt.Sprintf(`{"model":"my-model", "other": %q }`, strings.Repeat("a", 1_000_000)),
		},
	}
	for name, spec := range specs {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			b.Run("PeekModel", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _, err := PeekModel(strings.NewReader(spec.src), 128)
					require.NoError(b, err)
				}
			})

			b.Run("JsonDecoder", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					var target jsonType
					require.NoError(b, json.NewDecoder(strings.NewReader(spec.src)).Decode(&target))
				}
			})
		})
	}
}

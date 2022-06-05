package L_test

import (
	"io"
	"testing"

	"github.com/scott-cotton/L"
)

func BenchmarkBasic(b *testing.B) {
	L := L.New(&L.Config{
		Labels: map[string]int{},
		W:      io.Discard,
		//F:      &L.CLI{Fields: []string{"key3", "key0"}},
		F: L.JSONFmter(),
		E: L.EPanic,
	})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		L.Dict().
			Field("key0", 22).
			Field("key2", false).
			Field("key3", "hello susan").
			Log()
	}
}

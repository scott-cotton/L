package L_test

import (
	"io"
	"testing"

	"github.com/scott-cotton/L"
)

func BenchmarkConstruct1(b *testing.B) {
	b.StopTimer()
	L := L.New(&L.Config{
		Labels: map[string]int{},
		W:      io.Discard,
		F:      L.JSONFmter(),
		E:      L.EPanic,
	})
	o := L.Dict().
		Field("key0", 22).
		Field("key2", false).
		Field("key3", "hello susan").
		Field("key4", "bjez").
		Field("key5", 22.0).
		Field("key6", 77e10).
		Field("key7", "seven").
		Field("key8", "nine")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		o.Clone().Field("nine", "nine")
	}
}

func BenchmarkBasic(b *testing.B) {
	b.StopTimer()
	L := L.New(&L.Config{
		Labels: map[string]int{},
		W:      io.Discard,
		F:      L.JSONFmter(),
		E:      L.EPanic,
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

func BenchmarkConstructDict(b *testing.B) {
	b.StopTimer()
	L := L.New(&L.Config{
		Labels: map[string]int{},
		W:      io.Discard,
		F:      L.JSONFmter(),
		E:      L.EPanic,
	})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		L.Dict().
			Field("key0", 22).
			Field("key2", false).
			Field("key3", "hello susan")
	}
}

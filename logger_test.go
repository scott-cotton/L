package L_test

import (
	"bytes"
	"testing"

	"github.com/scott-cotton/L"
)

var testW = bytes.NewBuffer(nil)

var testL = L.New(&L.Config{
	W: testW,
	F: &L.TableFmter{
		Fields: []string{"key0", "key2"},
		Sep:    " ",
	},
	E: L.EPanic,
})

func TestBasic(t *testing.T) {
	testL.Dict().
		Field("key0", "hello kitty").
		Field("key1", 11).
		Field("key2", false).Log()
	t.Logf("%s", testW.String())
}

package L_test

import (
	"bytes"
	"fmt"

	"github.com/scott-cotton/L"
)

var eW = bytes.NewBuffer(nil)
var cfg = &L.Config{
	W: eW,
	F: L.JSONFmter(),
	E: func(l L.Logger, _ *L.Config, e error) {
		fmt.Printf("logger couldn't log obj: %s", e.Error())
	},
}

var eL = L.New(cfg)

func Example_configE() {
	// make a mistake and don't create real json.
	eL.Str("hello").Str("again").Log()
	fmt.Printf("%s\n", eW.String())

	// Output:
	// logger couldn't log obj: invalid character ',' after top-level value
}

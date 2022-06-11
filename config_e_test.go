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
	E: L.ELog,
}

var eL = L.New(cfg)

func Example_configE() {
	// make a mistake and don't create real json.
	eL.Str("hello").Str("again").Log()
	fmt.Printf("%s\n", eW.String())

	// Output:
	// {"LE":"invalid character ',' after top-level value"}
}

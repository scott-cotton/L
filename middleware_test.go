package L_test

import (
	"bytes"
	"fmt"

	"github.com/scott-cotton/L"
)

var out = bytes.NewBuffer(nil)

var config = &L.Config{
	W:    out,
	F:    L.JSONFmter(),
	E:    L.EPanic,
	Post: []L.Middleware{L.Pkg()},
}

var mwL = L.New(config)

func Example_pkgMiddleware() {
	d := mwL.Dict()
	d.Field("i", 3)
	mwL.Log(d)

	fmt.Println(out.String())

	// Output:
	// {"i":3,"Lpkg":"github.com/scott-cotton/L_test"}
}

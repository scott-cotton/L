package L_test

import (
	"bytes"
	"fmt"

	"github.com/scott-cotton/L"
)

type LogLevel int

const (
	Panic LogLevel = 1 + iota
	Fatal
	Warn
	Notice
	Info
	Debug
	Trace
)

func (lvl LogLevel) String() string {
	return map[LogLevel]string{
		Panic:  "panic",
		Fatal:  "fatal",
		Warn:   "warn",
		Notice: "notice",
		Info:   "info",
		Debug:  "debug",
		Trace:  "trace",
	}[lvl]
}

// set output to a global buffer for testing.
var lvlW = bytes.NewBuffer(nil)

func tagLevel(l L.Logger, cfg *L.Config, o *L.Obj) *L.Obj {
	v := cfg.Labels["Lop"]
	if v == 0 {
		return o
	}
	return o.Field("Lop", LogLevel(v).String())
}

func LConfig(lvl LogLevel) *L.Config {
	res := &L.Config{
		Labels: map[string]int{"Lop": int(lvl)},
		W:      lvlW,
		F:      L.JSONFmter(),
		E:      L.EPanic,
		Pre: []L.Middleware{
			L.Geq("Lop", int(lvl)),
		},
		Post: []L.Middleware{L.Pkg(), tagLevel},
	}
	return res
}

var (
	Ltrace  = L.New(LConfig(Trace))
	Ldebug  = Ltrace.With("Lop", int(Debug))
	Lnotice = Ldebug.With("Lop", int(Notice))
	Linfo   = Ltrace.With("Lop", int(Info))
	Lwarn   = Ltrace.With("Lop", int(Warn))
	Lfatal  = Ltrace.With("Lop", int(Fatal))
	Lpanic  = Ltrace.With("Lop", int(Panic))
)

func Example_levels() {
	Ltrace.Dict().Field("hello-trace", 22).Log()
	Ldebug.Dict().Field("hello-trace", 33).Log()
	fmt.Printf("%s\n", lvlW.String())

	// Output:
	// {"hello-trace":22,"Lpkg":"github.com/scott-cotton/L_test","Lop":"trace"}

}

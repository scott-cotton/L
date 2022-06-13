package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/scott-cotton/L"
	"github.com/scott-cotton/L/rpc"
)

const usage = `Lctl <name|url> <cmd>
<cmd> can be one of
- loggers
	retrieve all label information about loggers listening on <url>.
- apply <input>
	<input> can be a file or '-' for standard input
`

var logger = L.New(&L.Config{
	Labels: map[string]int{},
	Post:   []L.Middleware{L.Pkg(), L.TimeFormat("time", time.RFC1123)},
	W:      os.Stderr,
	E:      L.EPanic,
	F: &L.TableFmter{
		Fields: []string{"msg", "time", "method", "Lerr"},
		Sep:    " ",
	},
})

var Lerr = logger.With("Lcmd.err", 1)

var url = flag.String("addr", "http://localhost:4321/L", "url of L service")
var key = flag.String("key", "", "key for communicating with addr")

func main() {
	wo := Lerr.Dict()
	flag.Parse()
	client, err := rpc.NewClient(*key, *url)
	if err != nil {
		wo.Err(err).Fatal()
		os.Exit(1)
	}
	args := flag.Args()
	if len(args) == 0 {
		wo.Fmt("no args specified, usage:\n%s", usage).Fatal()
	}
	switch args[0] {
	case "loggers":
		res, err := client.Loggers()
		if err != nil {
			wo.Err(err).Fatal()
		}
		jenc := json.NewEncoder(os.Stdout)
		jenc.SetIndent("", "  ")
		if err := jenc.Encode(res); err != nil {
			wo.Err(err).Fatal()
		}
	case "apply":
		if len(args) == 1 {
			wo.Err(fmt.Errorf("no args specified, usage:\n%s", usage)).Fatal()
		}
		fname := args[1]
		r := os.Stdin
		var err error
		if fname != "-" {
			r, err = os.Open(fname)
			if err != nil {
				wo.Err(err).Fatal()
			}
			defer r.Close()
		}
		var params rpc.ApplyParams
		if err := json.NewDecoder(r).Decode(&params); err != nil {
			wo.Err(err).Fatal()
		}
		res, err := client.Apply(&params)
		if err != nil {
			wo.Err(err).Fatal()
		}
		jenc := json.NewEncoder(os.Stdout)
		jenc.SetIndent("", "  ")
		if err := jenc.Encode(res); err != nil {
			wo.Err(err).Fatal()
		}

	default:
		wo.Fmt("unknown method %q", args[0]).Fatal()
	}
}

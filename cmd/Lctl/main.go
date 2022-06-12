package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

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
	Post:   []L.Middleware{L.Pkg()},
	W:      os.Stderr,
	E:      L.EPanic,
	F: &L.TableFmter{
		Fields: []string{"msg", "time", "method", "Lerr"},
		Sep:    " ",
	},
})

var warn = logger.With("warn", 1)

var url = flag.String("addr", "http://localhost:4321/L", "url of L service")
var key = flag.String("key", "", "key for communicating with addr")

func main() {
	wo := warn.Dict()
	flag.Parse()
	client, err := rpc.NewClient(*key, *url)
	if err != nil {
		wo.Err(err).Log()
		os.Exit(1)
	}
	args := flag.Args()
	if len(args) == 0 {
		wo.Err(fmt.Errorf("no args specified, usage:\n%s", usage)).Log()
		os.Exit(2)
	}
	switch args[0] {
	case "loggers":
		res, err := client.Loggers()
		if err != nil {
			wo.Err(err).Log()
			os.Exit(3)
		}
		jenc := json.NewEncoder(os.Stdout)
		jenc.SetIndent("", "  ")
		if err := jenc.Encode(res); err != nil {
			wo.Err(err).Log()
			os.Exit(4)
		}
	case "apply":
		if len(args) == 1 {
			wo.Err(fmt.Errorf("no args specified, usage:\n%s", usage)).Log()
			os.Exit(2)
		}
		fname := args[1]
		r := os.Stdin
		var err error
		if fname != "-" {
			r, err = os.Open(fname)
			if err != nil {
				warn.Dict().Err(err).Log()
				os.Exit(4)
			}
			defer r.Close()
		}
		var params rpc.ApplyParams
		if err := json.NewDecoder(r).Decode(&params); err != nil {
			warn.Dict().Err(err).Log()
			os.Exit(4)
		}
		res, err := client.Apply(&params)
		if err != nil {
			wo.Err(err).Log()
			os.Exit(4)
		}
		jenc := json.NewEncoder(os.Stdout)
		jenc.SetIndent("", "  ")
		if err := jenc.Encode(res); err != nil {
			wo.Err(err).Log()
		}

	default:
		wo.Err(fmt.Errorf("unknown method %q", args[0])).Log()
	}
}

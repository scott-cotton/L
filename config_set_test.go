package L_test

import (
	"fmt"
	"sort"

	"github.com/scott-cotton/L"
)

var setCfg = L.NewConfig("Lset")

func Example_config_set() {
	cfg := L.NewConfig("a", "b", "c", ".d", "e")
	setCfg.Labels["b"] = 11
	setCfg.Labels[".d"] = 73
	cfg.Apply(setCfg)
	keys := make([]string, 0, len(cfg.Labels))
	for k := range cfg.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s: %d\n", key, cfg.Labels[key])
	}

	// Output:
	// a: 0
	// b: 11
	// c: 0
	// e: 0
	// github.com/scott-cotton/L_test.d: 73

}

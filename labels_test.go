package L_test

import (
	"fmt"
	"sort"

	"github.com/scott-cotton/L"
)

var lblConfig = L.NewConfig("a", ".b", "_c:")

func Example_labels() {
	lbls := lblConfig.Labels
	// sort keys for replicable test
	keys := []string{}
	for k := range lbls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, lbl := range keys {
		fmt.Printf("%s\n", lbl)
	}

	// Output:
	// _c:
	// a
	// github.com/scott-cotton/L_test.b
}

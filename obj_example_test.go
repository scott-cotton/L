package L_test

import (
	"encoding/json"
	"fmt"
)

func Example_obj() {
	s := "hello"
	obj := testL.Str(s)
	var t string
	json.Unmarshal(obj.D(), &t)
	if t != s {
		fmt.Printf("oops")
		return
	}
	fmt.Printf("obj matched\n")

	// Output:
	// obj matched
}

func Example_twoChildrenPanic() {
	defer func() {
		e := recover()
		if e == nil {
			fmt.Printf("no panic\n")
			return
		}
		fmt.Printf("panic %s\n", e)
	}()
	d := testL.Dict()
	c1 := d.Array()
	c2 := d.Array()
	_, _ = c1, c2

	// Output:
	// panic nonlinear
}

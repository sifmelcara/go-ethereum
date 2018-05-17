package ret_parser

import "fmt"

func ExampleNegInt() {
	var buf bytes.Buffer

	f := [1]byte{INT}

	d := Parse(f, "-123")

	fmt.Printf("%v\n", d)

	// Output: [-1]
}
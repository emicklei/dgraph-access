package dga

import "fmt"

func ExampleBlankUID() {
	fmt.Println(BlankUID("canada").RDF())
	// Output: _:canada
}

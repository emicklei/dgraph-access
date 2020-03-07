package dga

import "fmt"

func ExampleBlankUID() {
	fmt.Println(BlankUID("canada").RDF())
	// Output: _:canada
}

func ExampleIntegerUID() {
	fmt.Println(IntegerUID(42).RDF())
	// Output: <0x2a>
}

func ExampleStringUID() {
	fmt.Println(StringUID("name").RDF())
	// Output: <name>
}

func ExampleBlankNQuad() {
	fmt.Println(BlankNQuad("subject", "predicate", 42).RDF())
	// Output: _:subject <predicate> "42" .
}

func ExampleBlankNQuad_blankobject() {
	fmt.Println(BlankNQuad("parent-Sylvia", "children", BlankUID("child-Max")).RDF())
	// Output: _:parent-Sylvia <children> _:child-Max .
}

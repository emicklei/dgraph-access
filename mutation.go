package dga

import (
	"bytes"
	"fmt"
)

// Mutation represents an action with multiple RDF Triples.
type Mutation struct {
	// set, delete
	Command string
	NQuads  []NQuad
}

// RDF returns the string representation of the Mutation.
func (m Mutation) RDF() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{\n\t%s {", m.Command)
	for _, each := range m.NQuads {
		fmt.Fprintf(b, "\n\t\t")
		b.Write(each.Bytes())
	}
	fmt.Fprintf(b, "\n\t}\n}")
	return b.String()
}

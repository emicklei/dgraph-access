package dga

import (
	"bytes"
	"fmt"
)

const (
	// Star is used to model any predicate or any object in an NQuad.
	Star = "*"
)

// NQuad represents an RDF S P O pair.
type NQuad struct {
	// Subject is the node for which the predicate must be created/modified.
	Subject UID
	// Predicate is a known schema predicate or a Star
	Predicate string
	// Object can be a primitive value or a UID or a Star (constant)
	Object interface{}
}

// Bytes returns the mutation line.
func (n NQuad) Bytes() []byte {
	b := new(bytes.Buffer)
	b.WriteString(n.Subject.NQuadString())
	if n.Predicate == Star {
		fmt.Fprint(b, " * ")
	} else {
		fmt.Fprintf(b, " <%s> ", n.Predicate)
	}
	if s, ok := n.Object.(string); ok {
		if s == Star {
			fmt.Fprint(b, "* ")
		} else {
			fmt.Fprintf(b, "%q ", s)
		}
	} else if uid, ok := n.Object.(UID); ok {
		fmt.Fprintf(b, "%s ", uid.NQuadString())
	} else {
		fmt.Fprintf(b, "%v ", n.Object)
	}
	fmt.Fprintf(b, ".")
	return b.Bytes()
}

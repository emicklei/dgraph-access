package dga

import (
	"bytes"
	"fmt"
	"time"
)

const (
	// Star is used to model any predicate or any object in an NQuad.
	Star = "*"

	// DateTimeFormat is the format used by Dgraph for facet values of type dateTime.
	DateTimeFormat = "2006-01-02T15:04:05"

	// DgraphType is a reserved predicate name to refer to a type definition.
	DgraphType = "dgraph.type"
)

// NQuad represents an RDF S P O pair.
type NQuad struct {
	// Subject is the node for which the predicate must be created/modified.
	Subject UID

	// Predicate is a known schema predicate or a Star
	Predicate string

	// Object can be a primitive value or a UID or a Star (constant)
	Object interface{}

	// StorageType is used to optionally specify the type when storing the object
	// see https://docs.dgraph.io/mutations/#language-and-rdf-types
	// Example: dga.RDFString
	StorageType RDFDatatype

	// Maps to string, bool, int, float and dateTime.
	// For int and float, only 32-bit signed integers and 64-bit floats are accepted.
	Facets map[string]interface{}
}

// BlankNQuad returns an NQuad value with a Blank UID subject.
// Use BlankUID if you want the object also to be a Blank UID from a name.
func BlankNQuad(subjectName string, predicate string, object interface{}) NQuad {
	return NQuad{
		Subject:   BlankUID(subjectName),
		Predicate: predicate,
		Object:    object,
	}
}

// RDFDatatype is to set the StorageType of an NQuad.
type RDFDatatype string

// see https://docs.dgraph.io/mutations/#language-and-rdf-types
const (
	// RDFString is a RDF type
	RDFString   = RDFDatatype("<xs:string>")
	RDFDateTime = RDFDatatype("<xs:dateTime>")
	RDFDate     = RDFDatatype("<xs:date>")
	RDFInteger  = RDFDatatype("<xs:int>")
	RDFBoolean  = RDFDatatype("<xs:boolean>")
	RDFDouble   = RDFDatatype("<xs:double>")
	RDFFloat    = RDFDatatype("<xs:float>")
)

// WithStorageType returns a copy with its StorageType set.
// Use DetectStorageType(any interface{})
func (n NQuad) WithStorageType(t RDFDatatype) NQuad {
	return NQuad{
		Subject:     n.Subject,
		Predicate:   n.Predicate,
		Object:      n.Object,
		StorageType: t,
		Facets:      n.Facets,
	}
}

// WithFacet returns a copy with an additional facet (key=value).
func (n NQuad) WithFacet(key string, value interface{}) NQuad {
	f := n.Facets
	if f == nil {
		f = map[string]interface{}{}
	}
	f[key] = value
	return NQuad{
		Subject:     n.Subject,
		Predicate:   n.Predicate,
		Object:      n.Object,
		StorageType: n.StorageType,
		Facets:      f,
	}
}

// Bytes returns the mutation line.
func (n NQuad) Bytes() []byte {
	b := new(bytes.Buffer)
	b.WriteString(n.Subject.RDF())
	if n.Predicate == Star {
		fmt.Fprint(b, " * ")
	} else {
		fmt.Fprintf(b, " <%s> ", n.Predicate)
	}
	if s, ok := n.Object.(string); ok {
		if s == Star {
			fmt.Fprint(b, "*")
		} else {
			fmt.Fprintf(b, "%q", s)
		}
	} else if uid, ok := n.Object.(UID); ok {
		fmt.Fprintf(b, "%s", uid.RDF())
	} else if s, ok := n.Object.(string); ok {
		fmt.Fprintf(b, "%q", s)
	} else if i, ok := n.Object.(int); ok {
		fmt.Fprintf(b, "\"%d\"", i)
	} else {
		fmt.Fprintf(b, "%v", n.Object)
	}
	if len(n.StorageType) > 0 {
		fmt.Fprintf(b, "^^%s ", n.StorageType)
	} else {
		fmt.Fprintf(b, " ")
	}
	if n.Facets != nil && len(n.Facets) > 0 {
		fmt.Fprintf(b, "(")
		first := true
		for k, v := range n.Facets {
			if !first {
				fmt.Fprintf(b, ", ")
			}
			var s string
			if t, ok := v.(time.Time); ok {
				s = t.Format(DateTimeFormat)
			} else if q, ok := v.(string); ok {
				s = fmt.Sprintf("%q", q)
			} else {
				s = fmt.Sprintf("%v", v)
			}
			fmt.Fprintf(b, "%s=%s", k, s)
			first = false
		}
		fmt.Fprintf(b, ") ")
	}
	fmt.Fprintf(b, ".")
	return b.Bytes()
}

// RDF returns the string version of its Bytes representation.
func (n NQuad) RDF() string { return string(n.Bytes()) }

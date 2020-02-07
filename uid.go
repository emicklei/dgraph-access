package dga

import (
	"bytes"
	"fmt"
)

// HasUID is used in CreateNode to set the assigned UID to a typed value.
type HasUID interface {
	SetUID(uid UID)
	GetUID() UID
	SetType(typeName string)
	GetTypes() []string
}

// unknownUID is the zero UID, uninitialized
var unknownUID = UID{}

// UID represents a DGraph uid which can be expressed using an integer,string or undefined value.
type UID struct {
	// Str is exposed for JSON marshalling. Do not use it to read/write it directly.
	Str string
	raw string
}

// BlankUID returns an UID with an undefined uid and a local name only valid for one write transaction.
// .RDF() => _:name
func BlankUID(name string) UID {
	return UID{raw: "_:" + name}
}

// StringUID returns an UID using a string value for uid.
// .RDF() => <id>
func StringUID(id string) UID {
	return UID{Str: id}
}

// FunctionUID returns an UID calling the uid function on the argument.
// .RDF() => uid(s)
func FunctionUID(s string) UID {
	return UID{raw: fmt.Sprintf("uid(%s)", s)}
}

// NewUID returns an UID that is printed as is.
// .RDF() => s
func NewUID(s string) UID {
	return UID{raw: s}
}

// IntegerUID returns an UID using the integer value.
// .RDF() => <0x...>
func IntegerUID(i int) UID {
	return UID{raw: fmt.Sprintf("<0x%x>", uint64(i))}
}

// IsZero returns whether this UID is a zero value
func (u UID) IsZero() bool {
	return u == unknownUID || len(u.Str) == 0 && len(u.raw) == 0
}

// String is for debugging only
func (u UID) String() string {
	return fmt.Sprintf("UID(%s)", u.RDF())
}

// RDF returns a string presentation for use in an NQuad.
func (u UID) RDF() string {
	if len(u.raw) > 0 {
		return u.raw
	}
	return fmt.Sprintf("<%s>", u.Str)
}

// Assigned TODO
func (u UID) Assigned() string {
	// TODO handle raw
	//i, _ := strconv.ParseInt(u.Str, 10, 64)
	return u.Str
}

// MarshalJSON is part of JSON
func (u UID) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%q", u.RDF())
	return b.Bytes(), nil
}

// UnmarshalJSON is part of JSON
func (u *UID) UnmarshalJSON(data []byte) error {
	u.Str = string(data[1 : len(data)-1])
	return nil
}

package dga

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// HasUID is used in CreateNode to set the assigned UID to a typed value.
type HasUID interface {
	SetUID(uid UID)
	GetUID() UID
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
func BlankUID(name string) UID {
	return UID{raw: "_:" + name}
}

// StringUID returns an UID using a string value for uid, will be printed as "<some_id>".
func StringUID(s string) UID {
	return UID{Str: s}
}

// FunctionUID returns an UID that is printed as "uid(s)".
func FunctionUID(s string) UID {
	return UID{raw: fmt.Sprintf("uid(%s)", s)}
}

// NewUID returns an UID that is printed in RDF as is.
func NewUID(s string) UID {
	return UID{raw: s}
}

// IntegerUID returns an UID using the integer value.
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

// MarshalJSON is part of JSON
func (u UID) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%q", u.RDF())
	return b.Bytes(), nil
}

// UnmarshalJSON is part of JSON
func (u *UID) UnmarshalJSON(data []byte) error {
	type uid struct {
		UID string `json:"uid"`
	}
	var r uid
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	u.Str = r.UID
	return nil
}

package dga

import "fmt"

// HasUID is used in CreateNode to set the assigned UID to a typed value.
type HasUID interface {
	SetUID(uid UID)
	GetUID() UID
}

// unknownUID is the zero UID, uninitialized
var unknownUID = UID{}

// UID represents a DGraph uid which can be expressed using an integer,string or undefined value.
type UID struct {
	int uint64
	// Str is exposed for JSON marshalling. Do not use it to read/write it directly.
	Str   string `json:"uid"`
	blank string
}

// NewUID returns an UID with an undefined uid and a local name only valid for one write transaction.
func NewUID(name string) UID {
	return UID{blank: name}
}

// StringUID returns an UID using a string value for uid, will be printed as "<0x23>".
func StringUID(s string) UID {
	return UID{Str: s}
}

func IntegerUID(i int) UID {
	return UID{int: uint64(i)}
}

func (u UID) String() string {
	return fmt.Sprintf("uid(%s)", u.NQuadString())
}

func (u UID) NQuadString() string {
	if len(u.Str) > 0 {
		return fmt.Sprintf("<%s>", u.Str)
	}
	if u.int > 0 {
		return fmt.Sprintf("<0x%x>", u.int)
	}
	return "_:" + u.blank
}

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
	int uint64
	// Str is exposed for JSON marshalling. Do not use it to read/write it directly.
	Str   string
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

func (u UID) IsZero() bool {
	return u == unknownUID || u.int == 0 && len(u.blank) == 0 && len(u.Str) == 0
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

func (u UID) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%q", u.NQuadString())
	return b.Bytes(), nil
}

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

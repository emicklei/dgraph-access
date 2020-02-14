package dga

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStringUID(t *testing.T) {
	if got, want := StringUID("test").RDF(), "<test>"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBlankUID(t *testing.T) {
	if got, want := BlankUID("test").RDF(), "_:test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIntegerUID(t *testing.T) {
	if got, want := IntegerUID(42).RDF(), "<0x2a>"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFunctionUID(t *testing.T) {
	if got, want := FunctionUID("v").RDF(), "uid(v)"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNewUID(t *testing.T) {
	if got, want := NewUID("raw").RDF(), "raw"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestUIDJSON(t *testing.T) {
	type pair struct {
		uid      UID
		assigned string
	}
	pairs := append([]pair{},
		pair{NewUID("raw"), "raw"},
		pair{IntegerUID(42), "<0x2a>"},
		pair{StringUID("1234"), "<1234>"},
	)
	for _, u := range pairs {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.Encode(u.uid)
		w := UID{}
		json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&w)
		if got, want := w.Assigned(), u.assigned; got != want {
			t.Errorf("got [%v] want [%v]", got, want)
		}
	}
}

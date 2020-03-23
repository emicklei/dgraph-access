package dga

import (
	"testing"
)

func TestReflectNQuads(t *testing.T) {
	type Virus struct {
		Node
		Name  string `json:"alias,omitempty"`
		Empty string
	}
	v := new(Virus)
	v.Name = "corona"
	subject := NewUID("uid<node>")
	q := ReflectNQuads(subject, v)
	if got, want := len(q), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	first := q[0]
	if got, want := first.Subject, subject; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := first.Predicate, "alias"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := first.Object, v.Name; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

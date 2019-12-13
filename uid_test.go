package dga

import "testing"

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

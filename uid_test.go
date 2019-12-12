package dga

import "testing"

func TestStringUID(t *testing.T) {
	if got, want := StringUID("test").String(), "uid(<test>)"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBlankUID(t *testing.T) {
	if got, want := BlankUID("test").NQuadString(), "_:test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIntegerUID(t *testing.T) {
	if got, want := IntegerUID(42).NQuadString(), "<0x2a>"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

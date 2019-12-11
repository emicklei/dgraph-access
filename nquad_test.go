package dga

import "testing"

func TestNQuadString(t *testing.T) {
	q := NQuad{Subject: StringUID("0x34"), Predicate: "name", Object: "hello"}
	if got, want := string(q.Bytes()), `<0x34> <name> "hello" .`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNQuadString2(t *testing.T) {
	q := NQuad{Subject: StringUID("0x34"), Predicate: "name", Object: IntegerUID(0x57432143214)}
	if got, want := string(q.Bytes()), `<0x34> <name> <0x57432143214> .`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNQuadStarStar(t *testing.T) {
	q := NQuad{Subject: StringUID("0x34"), Predicate: Star, Object: Star}
	if got, want := string(q.Bytes()), `<0x34> * * .`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// https://docs.dgraph.io/mutations/#blank-nodes-and-uid
func TestNQuadTutorial1(t *testing.T) {
	q := NQuad{
		Subject:   BlankUID("class"),
		Predicate: "student",
		Object:    BlankUID("x"),
	}
	if got, want := string(q.Bytes()), "_:class <student> _:x ."; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

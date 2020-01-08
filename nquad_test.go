package dga

import "testing"

import "time"

func TestNQuadInteger(t *testing.T) {
	q := NQuad{Subject: StringUID("0x11"), Predicate: "name", Object: 42}
	if got, want := string(q.Bytes()), `<0x11> <name> "42" .`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

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
func TestNQuadSingleFacet(t *testing.T) {
	q := NQuad{
		Subject:   BlankUID("foo"),
		Predicate: "equals",
		Object:    BlankUID("bar"),
	}
	q = q.WithFacet("i", 24)
	if got, want := string(q.Bytes()), "_:foo <equals> _:bar (i=24) ."; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestNQuadStringAndDateTimeFacet(t *testing.T) {
	q := NQuad{
		Subject:   BlankUID("luke"),
		Predicate: "loves",
		Object:    BlankUID("force"),
	}
	when, _ := time.Parse(DateTimeFormat, "2006-01-02T15:04:05")
	q = q.WithFacet("when", when)
	q = q.WithFacet("weapon", "light saber")
	got := string(q.Bytes())
	one := "_:luke <loves> _:force (when=2006-01-02T15:04:05, weapon=\"light saber\") ."
	two := "_:luke <loves> _:force (weapon=\"light saber\", when=2006-01-02T15:04:05) ."
	if got != one && got != two {
		t.Error("missing or incorrect facet notation")
	}
}

func TestNQuadWithStringStorageType(t *testing.T) {
	q := NQuad{
		Subject:     BlankUID("foo"),
		Predicate:   "equals",
		Object:      "bar",
		StorageType: RDFString,
	}
	if got, want := string(q.Bytes()), "_:foo <equals> \"bar\"^^<xs:string> ."; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

package dga

import "testing"

func TestUID(t *testing.T) {
	if got, want := StringUID("test").String(), "uid(<test>)"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

package dga

import (
	"testing"
)

func TestEmbeddedNode(t *testing.T) {
	type Task struct {
		Node
		Name string
	}
	k := new(Task)
	k.SetUID(BlankUID("test"))
	k.SetType("Task")
	if got, want := k.GetTypes(), []string{"Task"}; len(got) != len(want) {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := k.GetTypes(), []string{"Task"}; got[0] != want[0] {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

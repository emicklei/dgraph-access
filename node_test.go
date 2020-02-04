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
	hasCheck(t, k)
	k.SetUID(BlankUID("test"))
	k.SetType("Task")
}

func hasCheck(t *testing.T, h HasUID) {
	if _, ok := h.(HasUID); !ok {
		t.Error("must implement HasUID")
	}
}

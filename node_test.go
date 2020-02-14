package dga

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestEmbeddedNode(t *testing.T) {
	type Task struct {
		Node
		Name string
	}
	k := new(Task)
	k.SetUID(StringUID("0x2a"))
	k.SetType("Task")
	if got, want := k.GetTypes(), []string{"Task"}; len(got) != len(want) {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := k.GetTypes(), []string{"Task"}; got[0] != want[0] {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(k); err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf.Bytes()))
	m := Task{}
	if err := json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&m); err != nil {
		t.Fatal(err)
	}
	if got, want := m.GetUID().Assigned(), "<0x2a>"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

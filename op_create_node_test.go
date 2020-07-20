package dga

import "testing"

type Bird struct {
	Node `json:",inline"`
	// scalar
	Name string `json:"name,omitempty"`
}

func TestCreateNode(t *testing.T) {
	t.Skip()
	b := &Bird{Name: "hawk"}
	c := &CreateNode{Node: b}
	dac := NewDGraphAccess(nil)
	_, err := c.Do(dac)
	if err != nil {
		t.Error(err)
	}
}

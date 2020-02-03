package dga

// Node is an abstract type to encapsulate a Dgraph identity (UID) and type (DType)
// Node can be used to embed in your own entity type, e.g.:
//
// type Person struct {
//      *dga.Node `json:",inline"`
//      Name string `json:"name"`
// }
type Node struct {
	UID   UID      `json:"uid,omitempty"`
	DType []string `json:"graph.dtype,omitempty"`
}

func (u *Node) SetUID(uid UID) { u.UID = uid }

func (u Node) GetUID() UID { return u.UID }

func NewNode(typeNames ...string) *Node {
	return &Node{
		UID:   unknownUID,
		DType: typeNames,
	}
}

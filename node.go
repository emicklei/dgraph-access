package dga

// Node is an abstract type to encapsulate a Dgraph identity (UID) and type (DType)
// Node can be used to embed in your own entity type, e.g.:
//
// type Person struct {
//      dga.Node    `json:",inline"`
//      Name string `json:"name"`
// }
type Node struct {
	UID   UID      `json:"uid,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

// SetUID sets the dgraph uid
func (n *Node) SetUID(uid UID) { n.UID = uid }

// GetUID gets the dgraph uid
func (n Node) GetUID() UID { return n.UID }

// SetType set or adds a graph.type for value that embeds the node.
func (n *Node) SetType(typeName string) {
	n.DType = append(n.DType, typeName)
}

// GetTypes returns the graph.type value(s).
func (n Node) GetTypes() []string { return n.DType }

// NewNode creates a new (embedded) Node with dgraph.type information.
func NewNode(typeNames ...string) *Node {
	return &Node{
		UID:   unknownUID,
		DType: typeNames,
	}
}

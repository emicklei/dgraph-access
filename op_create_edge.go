package dga

import (
	"github.com/dgraph-io/dgo/v2/protos/api"
)

// CreateEdge represents a Dgraph operation.
type CreateEdge struct {
	Subject   HasUID
	Predicate string
	Object    interface{}
	Facets    map[string]interface{}
}

// Do creates a new Edge (using an NQuad).
// If Subject is a non-created Node than create it first ; abort if error
// If Object is a non-created Node than create it first ; abort if error
// Returns an error if the mutation fails.
// Returns whether the edge was created when the absent check was requested.
// Requires a DGraphAccess with a Write transaction.
func (c CreateEdge) Do(d *DGraphAccess) (created bool, fail error) {
	if err := d.checkState(); err != nil {
		return false, err
	}
	// create subject if new Node
	if c.Subject.GetUID().IsZero() {
		if err := d.Fluent().CreateNode(c.Subject); err != nil {
			return false, err
		}
	}
	// create object if new Node
	object := c.Object
	if huid, ok := c.Object.(HasUID); ok {
		if huid.GetUID().IsZero() {
			if err := d.Fluent().CreateNode(huid); err != nil {
				return false, err
			}
		}
		object = huid.GetUID()
	}
	nq := NQuad{
		Subject:   c.Subject.GetUID(),
		Predicate: c.Predicate,
		Object:    object,
		Facets:    c.Facets,
	}
	nQuads := nq.Bytes()
	if d.traceEnabled {
		trace("CreateEdge", "nquad", string(nQuads))
	}
	_, err := d.txn.Mutate(d.ctx, &api.Mutation{SetNquads: nQuads})
	return true, err
}

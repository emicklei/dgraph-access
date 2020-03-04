package dga

import (
	"fmt"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

type CreateEdge struct {
	Subject   HasUID
	Predicate string
	Object    interface{}
	Facets    map[string]interface{}
	IfAbsent  bool
}

// Do creates a new Edge (using an NQuad).
// Returns an error if the mutation fails.
// Returns whether the edge was created when the absent check was requested.
// Requires a DGraphAccess with a Write transaction.
// If Subject is a non-created Node than create it first ; abort if error
// If Object is a non-created Node than create it first ; abort if error
func (c CreateEdge) Do(d *DGraphAccess) (created bool, fail error) {
	if err := d.checkState(); err != nil {
		return false, err
	}
	if c.Subject.GetUID().IsZero() {
		if err := d.CreateNode(c.Subject); err != nil {
			return false, err
		}
	}
	object := c.Object
	if uid, ok := c.Object.(HasUID); ok {
		if uid.GetUID().IsZero() {
			if err := d.CreateNode(uid); err != nil {
				return false, err
			}
		}
		object = uid.GetUID()
	}
	nq := NQuad{
		Subject:   c.Subject.GetUID(),
		Predicate: c.Predicate,
		Object:    object,
		Facets:    c.Facets,
	}
	nQuads := nq.Bytes()
	if d.traceEnabled {
		trace(fmt.Sprintf("RDF mutation (nquad): [%s]", string(nQuads)))
	}
	_, err := d.txn.Mutate(d.ctx, &api.Mutation{SetNquads: nQuads})
	return true, err
}

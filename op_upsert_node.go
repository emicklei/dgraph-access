package dga

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

// UpsertNode models the operation to insert (create) or update a Dgraph node.
type UpsertNode struct {
	Node      HasUID
	condition predicateCondition
}

// InsertUnless set the condition to update versus insert the node.
func (u *UpsertNode) InsertUnless(predicate string, object interface{}) {
	u.condition = predicateCondition{
		Predicate: predicate,
		Object:    object,
	}
}

// Do creates or updates a new Node.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (u UpsertNode) Do(d *DGraphAccess) (created bool, fail error) {
	if len(u.Node.GetTypes()) == 0 {
		u.Node.SetType(simpleType(u.Node))
	}
	dtype := u.Node.GetTypes()[0]
	subject := NewUID("uid(node)")
	nquads := ReflectNQuads(subject, u.Node)
	nquads = append(nquads, NQuad{
		Subject:   subject,
		Predicate: "dgraph.type",
		Object:    dtype,
	})
	data := bytesFromNQuads(nquads)
	mu := &api.Mutation{
		SetNquads: data,
	}
	query := fmt.Sprintf("query {node as q(func: type(%s)) @filter(%s) { uid }}", dtype, findFilterContent(u.condition.Predicate, u.condition.Object))
	if d.traceEnabled {
		trace("UpsertNode", query)
		trace("UpsertNode", "NQuads\n", string(data))
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
	}
	resp, err := d.txn.Do(d.ctx, req)
	if err != nil {
		return false, err
	}
	if d.traceEnabled {
		trace("UpsertNode", "resp", resp)
	}
	if len(resp.GetUids()) == 0 {
		// this means the node already exists
		// then we need to inspect the return of the query called "q" for getting the UID
		type queryResult struct {
			Q []struct {
				UID string `json:"uid"`
			} `json:"q"`
		}
		var qr queryResult
		if err := json.Unmarshal(resp.Json, &qr); err != nil {
			return false, err
		}
		u.Node.SetUID(StringUID(qr.Q[0].UID))
		return false, nil
	}
	var first string
	for _, uid := range resp.GetUids() {
		first = uid
		break
	}
	u.Node.SetUID(StringUID(first))
	return true, nil
}

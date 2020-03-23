package dga

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

// CreateNode models the operation to (conditionally) create a Dgraph node.
type CreateNode struct {
	Node      HasUID
	condition predicateCondition
}

// CreateUnless is a means to conditionally create the node. Create unless [predicate=object] for uids of the same type is true.
func (c *CreateNode) CreateUnless(predicate string, object interface{}) {
	c.condition = predicateCondition{
		Predicate: predicate,
		Object:    object,
	}
}

// Do creates a new Node.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (c CreateNode) Do(d *DGraphAccess) (created bool, fail error) {
	if c.condition.Object != nil {
		return c.conditional(d)
	}
	err := c.unconditional(d)
	return true, err
}

func (c CreateNode) unconditional(d *DGraphAccess) error {
	if err := d.CheckState(); err != nil {
		return err
	}
	if c.Node.GetUID().IsZero() {
		c.Node.SetUID(BlankUID("temp"))
	}
	if len(c.Node.GetTypes()) == 0 {
		c.Node.SetType(simpleType(c.Node))
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(c.Node)
	if err != nil {
		return fmt.Errorf("CreateNode|%v", err)
	}
	if d.traceEnabled {
		trace("CreateNode", "mutation", buf.String())
	}
	mu := &api.Mutation{SetJson: buf.Bytes()}
	resp, err := d.txn.Mutate(d.ctx, mu)
	if err != nil {
		return err
	}
	var first string
	for _, uid := range resp.GetUids() {
		first = uid
		break
	}
	c.Node.SetUID(StringUID(first))
	return nil
}

func (c *CreateNode) conditional(d *DGraphAccess) (created bool, fail error) {
	// TODO check uid of node
	c.Node.SetUID(BlankUID("temp"))
	if len(c.Node.GetTypes()) == 0 {
		c.Node.SetType(simpleType(c.Node))
	}
	dtype := c.Node.GetTypes()[0]
	data, err := json.Marshal(c.Node)
	if err != nil {
		return false, err
	}
	mu := &api.Mutation{
		Cond:    `@if(eq(len(node), 0))`,
		SetJson: data,
	}
	query := fmt.Sprintf("query {node as var(func: type(%s)) @filter(%s)}", dtype, findFilterContent(c.condition.Predicate, c.condition.Object))
	if d.traceEnabled {
		trace("CreateNode", query)
		trace("CreateNode", "cond", mu.Cond)
		trace("CreateNode", "JSON", string(data))
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
		trace("CreateNode", "resp", resp)
	}
	if len(resp.GetUids()) == 0 {
		// not absent, no node created
		return false, nil
	}
	var first string
	for _, uid := range resp.GetUids() {
		first = uid
		break
	}
	c.Node.SetUID(StringUID(first))
	return true, nil
}

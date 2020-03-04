package dga

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

type CreateNode struct {
	Node      HasUID
	condition PredicateCondition
}

func (c *CreateNode) Condition(predictate string, object interface{}) {
	c.condition = PredicateCondition{
		Predicate: predictate,
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
	if err := d.checkState(); err != nil {
		return err
	}
	if c.Node.GetUID().IsZero() {
		c.Node.SetUID(BlankUID("temp"))
	}
	if len(c.Node.GetTypes()) == 0 {
		c.Node.SetType(simpleType(c.Node))
	}
	data, err := json.Marshal(c.Node)
	if err != nil {
		return err
	}
	if d.traceEnabled {
		trace("JSON mutation:", string(data))
	}
	mu := &api.Mutation{SetJson: data}
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
		trace("CreateNodeIfAbsent query:", query)
		trace("CreateNodeIfAbsent cond:", mu.Cond)
		trace("CreateNodeIfAbsent JSON:", string(data))
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
	}
	resp, err := d.txn.Do(d.ctx, req)
	if err != nil {
		return false, err
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

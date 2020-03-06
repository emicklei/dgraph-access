package dga

import "github.com/dgraph-io/dgo/v2/protos/api"

// AlterSchema represents a Dgraph operation.
type AlterSchema struct {
	Source string
}

// Do performs the alert operation on the client.
func (a AlterSchema) Do(d *DGraphAccess) error {
	if err := d.CheckState(); err != nil {
		return err
	}
	op := &api.Operation{Schema: a.Source}
	if d.traceEnabled {
		trace("AlterSchema", "src", a.Source)
	}
	return d.client.Alter(d.ctx, op)
}

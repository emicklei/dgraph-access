package dga

import "github.com/dgraph-io/dgo/v2/protos/api"

type AlterSchema struct {
	Source string
}

func (a AlterSchema) Do(d *DGraphAccess) error {
	if err := d.checkState(); err != nil {
		return err
	}
	op := &api.Operation{Schema: a.Source}
	if d.traceEnabled {
		trace("AlterSchema", "src", a.Source)
	}
	return d.client.Alter(d.ctx, op)
}

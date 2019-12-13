package dga

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo/v2"
)

func TestUpsertTwoNQuads(t *testing.T) {
	dac := NewDGraphAccess(new(dgo.Dgraph))
	dac.ctx = context.Background()
	dac.txn = &mockTransaction{T: t, ExpectedQuery: "?", ExpectedNquads: `_:luke <learns> "the force" .`}
	nQuads := []NQuad{BlankNQuad("luke", "learns", "the force")}
	if err := dac.Upsert("?", nQuads); err != nil {
		t.Error(err)
	}
}

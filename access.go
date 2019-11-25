package dga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"golang.org/x/net/context"
)

var (
	// ErrNoClient is a DGraphAccess state error
	ErrNoClient = errors.New("dgo client not initialized")

	// ErrNoTransaction is a DGraphAccess state error
	ErrNoTransaction = errors.New("dgo transaction not created")

	// ErrNoContext is a DGraphAccess state error
	ErrNoContext = errors.New("dgo transaction context not created")
)

// DGraphAccess is a decorator for a dgo.Client that holds a Context and Transaction to perform queries and mutations.
type DGraphAccess struct {
	client       *dgo.Dgraph
	ctx          context.Context
	txn          *dgo.Txn
	traceEnabled bool
}

// checkState verifies that the Access can be used for a transaction (write | read only)
func (d *DGraphAccess) checkState() error {
	if d.client == nil {
		return ErrNoClient
	}
	if d.txn == nil {
		return ErrNoTransaction
	}
	if d.ctx == nil {
		return ErrNoContext
	}
	return nil
}

// NewDGraphAccess returns a new DGraphAccess using a client.
func NewDGraphAccess(client *dgo.Dgraph) *DGraphAccess {
	return &DGraphAccess{
		client:       client,
		traceEnabled: false,
	}
}

// WithTraceLogging returns a copy of DGraphAccess that will trace parts of its internals.
func (d *DGraphAccess) WithTraceLogging() *DGraphAccess {
	return &DGraphAccess{
		client:       d.client,
		txn:          d.txn,
		ctx:          d.ctx,
		traceEnabled: true,
	}
}

// ForReadWrite returns a copy of DGraphAccess ready to perform mutations.
func (d *DGraphAccess) ForReadWrite(ctx context.Context) *DGraphAccess {
	return &DGraphAccess{
		client:       d.client,
		txn:          d.client.NewTxn(),
		ctx:          ctx,
		traceEnabled: d.traceEnabled,
	}
}

// ForReadOnly returns a copy of DGraphAccess ready to perform read operations only.
func (d *DGraphAccess) ForReadOnly(ctx context.Context) *DGraphAccess {
	return &DGraphAccess{
		client:       d.client,
		txn:          d.client.NewReadOnlyTxn(),
		ctx:          ctx,
		traceEnabled: d.traceEnabled,
	}
}

// AlterSchema uses a schema definition to change the current DGraph schema.
// This operation is idempotent.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) AlterSchema(source string) error {
	if err := d.checkState(); err != nil {
		return err
	}
	op := &api.Operation{Schema: source}
	return d.client.Alter(d.ctx, op)
}

// CommitTransaction completes the current transaction.
// Return an error if the DGraphAccess is in the wrong state or if the Commit fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CommitTransaction() error {
	if err := d.checkState(); err != nil {
		return err
	}
	t, c := d.txn, d.ctx
	d.ctx = nil
	d.txn = nil
	return t.Commit(c)
}

// DiscardTransaction aborts the current transaction (unless absent).
func (d *DGraphAccess) DiscardTransaction() error {
	if d.txn != nil && d.ctx != nil {
		err := d.txn.Discard(d.ctx)
		d.txn = nil
		d.ctx = nil
		return err
	}
	return nil
}

// InTransactionDo calls a function with a prepared DGraphAccess with a Write transaction.
// Return an error if the Commit fails.
func (d *DGraphAccess) InTransactionDo(ctx context.Context, do func(da *DGraphAccess) error) error {
	wtx := d.ForReadWrite(ctx)
	defer wtx.DiscardTransaction()
	if err := do(wtx); err != nil {
		return err
	}
	return wtx.CommitTransaction()
}

// CreateEdge creates a new Edge (using an NQuad).
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateEdge(subject HasUID, predicate string, object interface{}) error {
	if err := d.checkState(); err != nil {
		return err
	}
	if uid, ok := object.(HasUID); ok {
		object = uid.GetUID()
	}
	nq := NQuad{
		Subject:   subject.GetUID(),
		Predicate: predicate,
		Object:    object,
	}
	nQuads := nq.Bytes()
	if d.traceEnabled {
		trace(fmt.Sprintf("mutate nquad: [%s]", string(nQuads)))
	}
	_, err := d.txn.Mutate(d.ctx, &api.Mutation{SetNquads: nQuads})
	return err
}

// CreateNode creates a new Node        .
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateNode(node HasUID) error {
	if err := d.checkState(); err != nil {
		return err
	}
	if node.GetUID().IsZero() {
		node.SetUID(NewUID("temp"))
	}
	data, err := json.Marshal(node)
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
	node.SetUID(StringUID(first))
	return nil
}

// FindNodeWithTypeAndAttribute finds a node of with dgraph.type = <typeName> and with predicateName = <value>
func (d *DGraphAccess) FindNodeWithTypeAndPredicate(typeName, predicateName, value string) (UID, bool, error) {
	q := fmt.Sprintf(`query FindNodeWithTypeAndPredicate {
		q(func: type(%s)) @filter(eq(%s,%q)) {
		  uid		  
		}
	  }`, typeName, predicateName, value)
	if d.traceEnabled {
		trace(q)
	}
	resp, err := d.txn.Query(d.ctx, q)
	if err != nil {
		return unknownUID, false, err
	}
	if d.traceEnabled {
		trace(string(resp.Json))
	}
	result := map[string][]UID{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return unknownUID, false, err
	}
	findOne := result["q"]
	if len(findOne) == 0 {
		return unknownUID, false, nil
	}
	return findOne[0], true, nil
}

func trace(msg ...interface{}) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "[dgraph-access-trace]")
	for _, each := range msg {
		fmt.Fprintf(b, " %v", each)
	}
	log.Println(b.String())
}

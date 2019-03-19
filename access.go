package dga

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
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
	client *dgo.Dgraph
	ctx    context.Context
	txn    *dgo.Txn
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
		client: client,
	}
}

// ForReadWrite returns a copy of DGraphAccess ready to perform mutations.
func (d *DGraphAccess) ForReadWrite() *DGraphAccess {
	return &DGraphAccess{
		client: d.client,
		txn:    d.client.NewTxn(),
		ctx:    context.Background(),
	}
}

// ForReadOnly returns a copy of DGraphAccess ready to perform read operations only.
func (d *DGraphAccess) ForReadOnly() *DGraphAccess {
	return &DGraphAccess{
		client: d.client,
		txn:    d.client.NewReadOnlyTxn(),
		ctx:    context.Background(),
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
func (d *DGraphAccess) InTransactionDo(do func(da *DGraphAccess) error) error {
	wtx := d.ForReadWrite()
	defer wtx.DiscardTransaction()
	if err := do(wtx); err != nil {
		return err
	}
	return wtx.CommitTransaction()
}

// CreateEdge creates a new Edge (using an NQuad).
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateEdge(subject UID, predicate string, object interface{}) error {
	if err := d.checkState(); err != nil {
		return err
	}
	nq := NQuad{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
	_, err := d.txn.Mutate(d.ctx, &api.Mutation{SetNquads: nq.Bytes()})
	return err
}

// CreateNode creates a new Node        .
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateNode(node HasUID) error {
	if err := d.checkState(); err != nil {
		return err
	}
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	mu := &api.Mutation{SetJson: data}
	assigned, err := d.txn.Mutate(d.ctx, mu)
	if err != nil {
		return err
	}
	node.SetUID(StringUID(assigned.Uids["blank-0"]))
	return nil
}

// FindNodeHasEqualsString find the first Node that has a given predicate (hasPredicate)
// and equals to a given string value for another predicate (filterPredicate).
func (d *DGraphAccess) FindNodeHasEqualsString(hasPredicate, filterPredicate, value string) (UID, bool) {
	q := fmt.Sprintf(`query findOne($param: string) {
		findOne(func: has(%s)) @filter(eq(%s,$param)) {
		  uid		  
		}
	  }`, hasPredicate, filterPredicate)
	resp, err := d.txn.QueryWithVars(d.ctx, q, map[string]string{"$param": value})
	if err != nil {
		log.Println(err)
		return unknownUID, false
	}
	//{"findOne":[{"uid":"0x2718"}]}
	result := map[string][]UID{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Println(err)
		return unknownUID, false
	}
	findOne := result["findOne"]
	if len(findOne) == 0 {
		return unknownUID, false
	}
	return findOne[0], true
}

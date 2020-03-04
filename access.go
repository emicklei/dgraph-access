package dga

import (
	"errors"

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

	// ErrNoResultsFound is returned from FindEquals when no node matches.
	ErrNoResultsFound = errors.New("no results found")

	// ErrUnmarshalQueryResult is returned when the result of a query cannot be unmarshalled from JSON
	ErrUnmarshalQueryResult = errors.New("failed to unmarshal query result")
)

// DGraphAccess is a decorator for a dgo.Client that holds a Context and Transaction to perform queries and mutations.
type DGraphAccess struct {
	client       *dgo.Dgraph
	ctx          context.Context
	txn          DGraphTransaction
	traceEnabled bool
}

// DGraphTransaction exists for testing. It has only the methods this package needs from a *dgo.Txn
type DGraphTransaction interface {
	Mutate(ctx context.Context, mu *api.Mutation) (*api.Response, error)
	Commit(ctx context.Context) error
	Discard(ctx context.Context) error
	Do(ctx context.Context, req *api.Request) (*api.Response, error)
	Query(ctx context.Context, q string) (*api.Response, error)
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

// Transaction returns the encapsulated transaction (if present).
func (d *DGraphAccess) Transaction() *dgo.Txn {
	if d.txn == nil {
		return nil
	}
	if nonMock, ok := d.txn.(*dgo.Txn); ok {
		return nonMock
	}
	return nil
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

// Commit completes the current transaction.
// Return an error if the DGraphAccess is in the wrong state or if the Commit fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) Commit() error {
	if err := d.checkState(); err != nil {
		return err
	}
	t, c := d.txn, d.ctx
	d.ctx = nil
	d.txn = nil
	return t.Commit(c)
}

// Discard aborts the current transaction (unless absent).
func (d *DGraphAccess) Discard() error {
	if d.txn != nil && d.ctx != nil {
		err := d.txn.Discard(d.ctx)
		d.txn = nil
		d.ctx = nil
		return err
	}
	return nil
}

// InTransactionDo calls a function with a prepared DGraphAccess with a Write transaction.
// The encapsulated transaction is available from the receiver using Transaction().
// Return an error if the Commit fails.
func (d *DGraphAccess) InTransactionDo(ctx context.Context, do func(da *DGraphAccess) error) error {
	wtx := d.ForReadWrite(ctx)
	defer wtx.Discard()
	if err := do(wtx); err != nil {
		return err
	}
	return wtx.Commit()
}

// Operation is for dispatching commands using a DGraphAccess.
type Operation interface {
	Do(d *DGraphAccess) (hadEffect bool, err error)
}

// Do executes the operation. Return whether the operation had effect or an error.
func (d *DGraphAccess) Do(o Operation) (hadEffect bool, err error) {
	return o.Do(d)
}

// Fluent gives access to the fluent interface of DGraphAccess.
func (d *DGraphAccess) Fluent() Fluent {
	return Fluent{access: d}
}

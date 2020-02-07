package dga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

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

// NoFacets can be used in CreateEdge for passing no facets.
var NoFacets map[string]interface{} = nil

// CreateEdge creates a new Edge (using an NQuad).
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateEdge(subject HasUID, predicate string, object interface{}) error {
	return d.CreateEdgeWithFacets(subject, predicate, object, NoFacets)
}

// CreateEdgeWithFacets creates a new Edge (using an NQuad) that has facets (can be nil or empty)
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateEdgeWithFacets(subject HasUID, predicate string, object interface{}, facetsOrNil map[string]interface{}) error {
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
		Facets:    facetsOrNil,
	}
	nQuads := nq.Bytes()
	if d.traceEnabled {
		trace(fmt.Sprintf("RDF mutation (nquad): [%s]", string(nQuads)))
	}
	_, err := d.txn.Mutate(d.ctx, &api.Mutation{SetNquads: nQuads})
	return err
}

// CreateNode creates a new Node.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) CreateNode(node HasUID) error {
	if err := d.checkState(); err != nil {
		return err
	}
	if node.GetUID().IsZero() {
		node.SetUID(BlankUID("temp"))
	}
	if len(node.GetTypes()) == 0 {
		node.SetType(simpleType(node))
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

// Upsert runs a mutation if the query yields no results.
// Requires a DGraphAccess with a Write transaction.
func (d *DGraphAccess) Upsert(query string, nQuads []NQuad) error {
	if err := d.checkState(); err != nil {
		return err
	}
	req := d.UpsertRequest(query, nQuads)
	r, err := d.txn.Do(d.ctx, req)
	if d.traceEnabled {
		trace(fmt.Sprintf("%#v", r))
	}
	return err
}

func (d *DGraphAccess) UpsertRequest(query string, nQuads []NQuad) *api.Request {
	b := new(bytes.Buffer)
	for _, each := range nQuads {
		b.Write(each.Bytes())
		b.WriteString("\n")
	}
	if d.traceEnabled {
		trace(fmt.Sprintf(`
upsert {
	%s
	mutation {
		set {
%s
		}
	}
}`, query, b.String()))
	}
	mu := &api.Mutation{SetNquads: b.Bytes()}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
	}
	return req
}

func simpleType(result interface{}) string {
	tokens := strings.Split(fmt.Sprintf("%T", result), ".")
	return tokens[len(tokens)-1]
}

var ErrNotFound = errors.New("node not found")
var ErrUnmarshalQueryResult = errors.New("failed to unmarshal query result")

// FindEquals populates the result with the result of matching a predicate with a value.
func (d *DGraphAccess) FindEquals(result interface{}, predicateName, value interface{}) error {
	st := simpleType(result)
	var valueString string
	if s, ok := value.(string); ok {
		valueString = fmt.Sprintf("\"%s\"", s)
	}
	if n, ok := value.(HasUID); ok {
		valueString = n.GetUID().QueryFunction()
	}
	q := fmt.Sprintf(`
query FindWithTypeAndPredicate {
	q(func: type(%s)) @filter(eq(%s,%s)) {
		uid	
		dgraph.type
		expand(%s)
	}
}`, st, predicateName, valueString, st)
	if d.traceEnabled {
		trace(q)
	}
	resp, err := d.txn.Query(d.ctx, q)
	if err != nil {
		return ErrNotFound
	}
	if d.traceEnabled {
		trace(string(resp.Json))
	}
	qresult := map[string][]interface{}{}
	err = json.Unmarshal(resp.Json, &qresult)
	if err != nil {
		return ErrUnmarshalQueryResult
	}
	findOne := qresult["q"]
	if len(findOne) == 0 {
		return ErrNotFound
	}
	// mapstructure pkg did not work for this case
	resultData := new(bytes.Buffer)
	json.NewEncoder(resultData).Encode(findOne[0])
	resultBytes := resultData.Bytes()
	return json.NewDecoder(bytes.NewReader(resultBytes)).Decode(result)
}

func trace(msg ...interface{}) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "[dgraph-access-trace]")
	for _, each := range msg {
		fmt.Fprintf(b, " %v", each)
	}
	log.Println(b.String())
}

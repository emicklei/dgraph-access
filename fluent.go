package dga

// Fluent gives access to the fluent interface of DGraphAccess.
type Fluent struct {
	access *DGraphAccess
}

// AlterSchema uses a schema definition to change the current DGraph schema.
// This operation is idempotent.
// Requires a DGraphAccess with a Write transaction.
func (f Fluent) AlterSchema(source string) error {
	a := AlterSchema{
		Source: source,
	}
	return a.Do(f.access)
}

// CreateEdge creates a new Edge (using an NQuad).
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
// If subject is a non-created Node than create it first ; abort if error
// If object is a non-created Node than create it first ; abort if error
func (f Fluent) CreateEdge(subject HasUID, predicate string, object interface{}) error {
	c := CreateEdge{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
	_, err := c.Do(f.access)
	return err
}

// CreateNode creates a new Node.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (f Fluent) CreateNode(node HasUID) error {
	c := CreateNode{
		Node: node,
	}
	_, err := c.Do(f.access)
	return err
}

// RunQuery executes the raw query and populates the result with the data found using a given key.
func (f Fluent) RunQuery(result interface{}, query string, dataKey string) (bool, error) {
	r := RunQuery{
		Result:  result,
		Query:   query,
		DataKey: dataKey,
	}
	return r.Do(f.access)
}

// FindEquals populates the result with the result of matching a predicate with a value.
func (f Fluent) FindEquals(result interface{}, predicateName string, value interface{}) (bool, error) {
	e := FindEquals{
		Result:    result,
		Predicate: predicateName,
		Object:    value,
	}
	return e.Do(f.access)
}

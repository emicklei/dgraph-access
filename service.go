package dga

// Service is the API to use the operation types with DGraphAccess.
type Service struct {
	access *DGraphAccess
}

// AlterSchema uses a schema definition to change the current DGraph schema.
// This operation is idempotent.
// Requires a DGraphAccess with a Write transaction.
func (s Service) AlterSchema(source string) error {
	a := AlterSchema{
		Source: source,
	}
	return a.Do(s.access)
}

// CreateEdge creates a new Edge (using an NQuad).
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
// If subject is a non-created Node than create it first ; abort if error
// If object is a non-created Node than create it first ; abort if error
func (s Service) CreateEdge(subject HasUID, predicate string, object interface{}) error {
	c := CreateEdge{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
	_, err := c.Do(s.access)
	return err
}

// CreateNode creates a new Node.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (s Service) CreateNode(node HasUID) error {
	c := CreateNode{
		Node: node,
	}
	_, err := c.Do(s.access)
	return err
}

// UpsertNode creates (insert) or updates a Node.
// The operation will update iff the predicate -> object.
// Return an error if the mutation fails.
// Requires a DGraphAccess with a Write transaction.
func (s Service) UpsertNode(node HasUID, predicate string, object interface{}) (wasCreated bool, err error) {
	c := UpsertNode{
		Node: node,
	}
	c.InsertUnless(predicate, object)
	return c.Do(s.access)
}

// RunQuery executes the raw query and populates the result with the data found using a given key.
func (s Service) RunQuery(result interface{}, query string, dataKey string) (bool, error) {
	r := RunQuery{
		Result:  result,
		Query:   query,
		DataKey: dataKey,
	}
	return r.Do(s.access)
}

// FindEquals populates the result with the result of matching a predicate with a value.
func (s Service) FindEquals(result interface{}, predicateName string, value interface{}) (bool, error) {
	e := FindEquals{
		Result:    result,
		Predicate: predicateName,
		Object:    value,
	}
	return e.Do(s.access)
}

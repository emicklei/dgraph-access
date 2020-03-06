package dga

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FindEquals populates the result with the result of matching a predicate with a value.
type FindEquals struct {
	Predicate string
	Object    interface{}
	Result    interface{}
}

// Do populates the result with the result of matching a predicate with a value.
func (f FindEquals) Do(d *DGraphAccess) (hadEffect bool, err error) {
	st := simpleType(f.Result)
	filterContent := findFilterContent(f.Predicate, f.Object)
	q := fmt.Sprintf(`
query FindWithTypeAndPredicate {
	q(func: type(%s)) @filter(%s) {
		uid	
		dgraph.type
		expand(%s)
	}
}`, st, filterContent, st)
	if d.traceEnabled {
		trace("FindEquals", "query", q)
	}
	resp, err := d.txn.Query(d.ctx, q)
	if err != nil {
		return false, err
	}
	if d.traceEnabled {
		trace("FindEquals", "resp", string(resp.Json))
	}
	qresult := map[string][]interface{}{}
	err = json.Unmarshal(resp.Json, &qresult)
	if err != nil {
		return false, ErrUnmarshalQueryResult
	}
	findOne := qresult["q"]
	if len(findOne) == 0 {
		return false, nil
	}
	// mapstructure pkg did not work for this case
	// TODO optimize this
	resultData := new(bytes.Buffer)
	json.NewEncoder(resultData).Encode(findOne[0])
	resultBytes := resultData.Bytes()
	return true, json.NewDecoder(bytes.NewReader(resultBytes)).Decode(f.Result)
}

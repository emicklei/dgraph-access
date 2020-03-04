package dga

import (
	"bytes"
	"encoding/json"
	"log"
)

// RunQuery executes the raw query and populates the result with the data found using a given key.
type RunQuery struct {
	Result  interface{}
	Query   string
	DataKey string
}

// Do executes the raw query and populates the result with the data found using a given key.
func (r RunQuery) Do(d *DGraphAccess) (hadEffect bool, err error) {
	if d.traceEnabled {
		trace("RunQuery", "query", r.Query)
	}
	resp, err := d.txn.Query(d.ctx, r.Query)
	if err != nil {
		// TODO check error
		log.Println(err)
		return false, ErrNoResultsFound
	}
	if d.traceEnabled {
		trace("RunQuery", "resp", string(resp.Json))
	}
	qresult := map[string][]interface{}{}
	err = json.Unmarshal(resp.Json, &qresult)
	if err != nil {
		return false, ErrUnmarshalQueryResult
	}
	findOne := qresult[r.DataKey]
	if len(findOne) == 0 {
		return false, ErrNoResultsFound
	}
	// mapstructure pkg did not work for this case
	// TODO optimize this
	resultData := new(bytes.Buffer)
	json.NewEncoder(resultData).Encode(findOne[0])
	resultBytes := resultData.Bytes()
	return true, json.NewDecoder(bytes.NewReader(resultBytes)).Decode(r.Result)
}

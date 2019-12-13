package dga

import (
	"strings"
	"testing"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"golang.org/x/net/context"
)

type mockTransaction struct {
	T              *testing.T
	ExpectedQuery  string
	ExpectedNquads string
}

func (m *mockTransaction) Mutate(ctx context.Context, mu *api.Mutation) (*api.Response, error) {
	return nil, nil
}
func (m *mockTransaction) Commit(ctx context.Context) error {
	return nil
}
func (m *mockTransaction) Discard(ctx context.Context) error {
	return nil
}
func (m *mockTransaction) Do(ctx context.Context, req *api.Request) (*api.Response, error) {
	if len(m.ExpectedQuery) > 0 {
		if got, want := req.Query, m.ExpectedQuery; got != want {
			m.T.Errorf("got [%v] want [%v]", got, want)
		}
	}
	if len(m.ExpectedNquads) > 0 {
		if got, want := clean(string(req.Mutations[0].GetSetNquads())), m.ExpectedNquads; got != want {
			m.T.Errorf("got [%v] want [%v]", got, want)
		}
	}
	return nil, nil
}
func (m *mockTransaction) Query(ctx context.Context, q string) (*api.Response, error) {
	return nil, nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", "")
}

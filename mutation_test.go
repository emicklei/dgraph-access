package dga

import (
	"testing"
)

func TestDeleteMutation(t *testing.T) {
	m := Mutation{
		Operation: "delete",
		NQuads: []NQuad{
			{
				Subject:   StringUID("0xf11168064b01135b"),
				Predicate: "died",
				Object:    1998},
		},
	}
	if got, want := m.RDF(), `{
	delete {
		<0xf11168064b01135b> <died> "1998" .
	}
}`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

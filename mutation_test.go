package dga

import (
	"testing"
)

func TestDeleteMutation(t *testing.T) {
	m := Mutation{
		Operation: "delete",
		NQuads: []NQuad{
			NQuad{
				Subject:   StringUID("0xf11168064b01135b"),
				Predicate: "died",
				Object:    1998},
		},
	}
	t.Log(m.RDF())
}

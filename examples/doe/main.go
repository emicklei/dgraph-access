package main

import (
	"context"
	"log"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	dga "github.com/emicklei/dgraph-access"
	"google.golang.org/grpc"
)

type Person struct {
	*dga.Node `json:",inline"`
	//
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

func main() {
	ctx := context.Background()
	client := newClient()

	// Warn: Cleaning up the database
	if err := client.Alter(ctx, &api.Operation{DropAll: true}); err != nil {
		log.Fatal(err)
	}

	// create an accessor
	dac := dga.NewDGraphAccess(client)

	// for debugging only
	dac = dac.WithTraceLogging()

	// set schema
	if err := dac.InTransactionDo(ctx, alterSchema); err != nil {
		log.Fatal(err)
	}

	// insert data
	if err := dac.InTransactionDo(ctx, insertData); err != nil {
		log.Fatal(err)
	}

	// query data
	dac = dac.ForReadOnly(ctx)

	// find using type and name
	p := Person{}
	err := dac.FindEquals(&p, "name", "John")
	if err != nil {
		log.Println(err)
	}
	log.Println("uid:", p.UID, "name:", p.Name, "surname:", p.Surname)
}

func insertData(da *dga.DGraphAccess) error {
	john := &Person{Node: dga.NewNode("Person"), Name: "John", Surname: "Doe"}
	jane := &Person{Node: dga.NewNode("Person"), Name: "Jane", Surname: "Doe"}
	if err := da.CreateNode(john); err != nil {
		return err
	}
	if err := da.CreateNode(jane); err != nil {
		return err
	}
	if err := da.CreateEdge(john, "isMarriedTo", jane); err != nil {
		return err
	}
	if err := da.CreateEdgeWithFacets(jane, "isMarriedTo", john, dga.NoFacets); err != nil {
		return err
	}
	props := map[string]interface{}{
		"style": "spanish",
	}
	if err := da.CreateEdgeWithFacets(jane, "likesToDanceWith", john, props); err != nil {
		return err
	}
	return nil
}

func alterSchema(da *dga.DGraphAccess) error {
	return da.AlterSchema(`
	name: string @index(exact) .
	surname: string .

	type Person {
		name
		surname
	}
`)
}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

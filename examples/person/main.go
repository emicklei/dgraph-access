package main

import (
	"context"
	"flag"
	"log"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	dga "github.com/emicklei/dgraph-access"
	"google.golang.org/grpc"
)

type Person struct {
	dga.Node `json:",inline"`
	// scalar
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

var drop = flag.Bool("drop", false, "cleanup the database at startup")

func main() {
	flag.Parse()
	ctx := context.Background()
	client := newClient()

	if *drop {
		log.Println("Cleaning up the database")
		if err := client.Alter(ctx, &api.Operation{DropAll: true}); err != nil {
			log.Fatal("drop all failed ", err)
		}
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
	ok, err := dac.Fluent().FindEquals(&p, "name", "John")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("uid:", p.UID, "name:", p.Name, "surname:", p.Surname)
	}

	// create Jack if missing
	jack := &Person{Name: "Jack", Surname: "Doe"}
	op := dga.CreateNode{Node: jack}
	op.Unless("name", jack.Name)

	dac.InTransactionDo(ctx, func(d *dga.DGraphAccess) error {
		_, err := d.Do(op)
		return err
	})
	log.Println("uid:", jack.UID, "name:", jack.Name, "surname:", jack.Surname)
}

func insertData(d *dga.DGraphAccess) error {
	john := &Person{Name: "John", Surname: "Doe"}
	jane := &Person{Name: "Jane", Surname: "Doe"}

	// use the operation
	op := dga.CreateEdge{
		Subject:   john,
		Predicate: "isMarriedTo",
		Object:    jane,
	}
	if _, err := d.Do(op); err != nil {
		return err
	}
	// use the fluent interface
	f := d.Fluent()
	if err := f.CreateEdge(jane, "isMarriedTo", john); err != nil {
		return err
	}
	if err := f.CreateEdge(john, "parent", &Person{Name: "Jesse", Surname: "Doe"}); err != nil {
		return err
	}

	// create with a facet requires to use the operation
	op = dga.CreateEdge{
		Subject:   jane,
		Predicate: "likesToDanceWith",
		Object:    john,
		Facets: map[string]interface{}{
			"style": "spanish",
		},
	}
	if _, err := d.Do(op); err != nil {
		return err
	}

	// update the first name of John


	return nil
}

func alterSchema(d *dga.DGraphAccess) error {
	return d.Fluent().AlterSchema(`
	name: string @index(exact) .
	surname: string .
	parent: uid .

	type Person {
		name
		surname
		parent
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

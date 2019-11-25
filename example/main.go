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
	// dgraph
	Uid   dga.UID `json:"uid,omitempty"`
	DType string  `json:"dgraph.type,omitempty"`
	//
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

func (p *Person) SetUID(uid dga.UID) {
	p.Uid = uid
}

func (p Person) GetUID() dga.UID {
	return p.Uid
}

func main() {
	ctx := context.Background()
	client := newClient()

	// Warn: Cleaning up the database
	// if err := client.Alter(ctx, &api.Operation{DropAll: true}); err != nil {
	// 	log.Fatal(err)
	// }

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

	uid, ok, err := dac.FindNodeWithTypeAndPredicate("Person", "name", "John")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Println("not found")
		return
	}
	log.Println("uid:", uid)
}

func insertData(da *dga.DGraphAccess) error {
	john := &Person{Name: "John", Surname: "Doe", DType: "Person"}
	jane := &Person{Name: "Jane", Surname: "Doe", DType: "Person"}
	if err := da.CreateNode(john); err != nil {
		return err
	}
	if err := da.CreateNode(jane); err != nil {
		return err
	}
	if err := da.CreateEdge(john, "isMarriedTo", jane); err != nil {
		return err
	}
	if err := da.CreateEdge(jane, "isMarriedTo", john); err != nil {
		return err
	}
	return nil
}

func alterSchema(da *dga.DGraphAccess) error {
	return da.AlterSchema(`
	name: string @index(exact) .

	type Person {
		name: string
		surname: string
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

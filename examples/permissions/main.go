package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	dga "github.com/emicklei/dgraph-access"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	client := newClient()

	// Warn: Cleaning up the database
	if err := client.Alter(ctx, &api.Operation{DropAll: true}); err != nil {
		log.Fatal("drop all failed ", err)
	}

	// create an accessor
	dac := dga.NewDGraphAccess(client)

	// for debugging only
	dac = dac.WithTraceLogging()

	// set schema
	if err := dac.InTransactionDo(ctx, alterSchema); err != nil {
		log.Fatal("alter schema failed ", err)
	}

	// insert data
	if err := dac.InTransactionDo(ctx, insertData); err != nil {
		log.Fatal(err)
	}

	// query data
	// dac = dac.ForReadOnly(ctx)
}

func insertData(xs *dga.DGraphAccess) error {
	sa := ServiceAccount{
		Node: dga.NewNode("ServiceAccount"),
		Name: "compute-default",
	}
	if err := xs.CreateNode(sa); err != nil {
		return err
	}
	fmt.Println("serviceAccount:", sa.UID)
	pr := Project{
		Node: dga.NewNode("Project"),
		Name: "my-project",
	}
	if err := xs.CreateNode(pr); err != nil {
		return err
	}
	fmt.Println("project:", pr.UID)
	pip := PermissionsInProject{
		Node:        dga.NewNode("PermissionsInProject"),
		Permissions: []string{"role/editor"},
	}
	if err := xs.CreateNode(pip); err != nil {
		return err
	}
	fmt.Println("permissions-in-project:", pip.UID)
	if err := xs.CreateEdge(pip, "serviceAccount", sa); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip.UID, "->", "service-account", sa.UID)
	return nil
}

func alterSchema(da *dga.DGraphAccess) error {
	content, err := ioutil.ReadFile("schema.txt")
	if err != nil {
		return err
	}
	return da.AlterSchema(string(content))
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

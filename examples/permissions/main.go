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

	// which permissions does user(john.doe) have?
	// Method 1: find john then find its permissions
	dac = dac.ForReadOnly(ctx)

	john := new(CloudIdentity)
	if err := dac.FindEquals(john, "user", "john.doe"); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", john)

	pip := new(PermissionsInProject)
	if err := dac.FindEquals(pip, "groupOrUser", john); err != nil {
		log.Fatal(err)
	}
	log.Printf("(with node) %#v", pip)

	{ // if you only have the uid of John
		pip := new(PermissionsInProject)
		if err := dac.FindEquals(pip, "groupOrUser", john.UID); err != nil {
			log.Fatal(err)
		}
		log.Printf("(uid only) %#v", pip)
	}

	// which permissions does user(john.doe) have?
	// Method 2: find permissions filtering groupOrUser predicate
	query := `{
		q(func: type(PermissionsInProject)) @cascade {
				groupOrUser @filter(eq(user,"john.doe"))
				permissions
		}
	  }`
	data := map[string][]string{}
	if err := dac.RunQuery(&data, query, "q"); err != nil {
		log.Fatal(err)
	}
	log.Printf("(filter predicate) %#v", data["permissions"])

	// for which projects has service account (compute-default) permissions (role/editor) ?
	query = `{
		q(func: type(PermissionsInProject)) @cascade {
				serviceAccount @filter(eq(serviceaccount_name,"compute-default"))
				permissions
		}
	  }`
	data = map[string][]string{}
	if err := dac.RunQuery(&data, query, "q"); err != nil {
		log.Fatal(err)
	}
	log.Printf("(filter account and permissions) %#v", data["permissions"])
}

func insertData(xs *dga.DGraphAccess) error {
	// serviceAcount(compute-default) has permission(role/editor) in project(my-project)
	sa := &ServiceAccount{
		Name: "compute-default",
	}
	if err := xs.CreateNode(sa); err != nil {
		return err
	}
	fmt.Println("serviceAccount:", sa.UID)
	pr := &Project{
		Name: "my-project",
	}
	if err := xs.CreateNode(pr); err != nil {
		return err
	}
	fmt.Println("project:", pr.UID)
	pip := &PermissionsInProject{
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

	// user(john.doe) has permission(role/viewer) in project(my-project)
	pip2 := &PermissionsInProject{
		Permissions: []string{"role/viewer"},
	}
	if err := xs.CreateNode(pip2); err != nil {
		return err
	}
	ci := &CloudIdentity{
		User: "john.doe",
	}
	if err := xs.CreateNode(ci); err != nil {
		return err
	}
	if err := xs.CreateEdge(pip2, "groupOrUser", ci); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip2.UID, "->", "groupOrUser", ci.UID)
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

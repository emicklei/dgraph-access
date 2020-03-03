package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	dga "github.com/emicklei/dgraph-access"
	"google.golang.org/grpc"
)

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
		log.Fatal("alter schema failed ", err)
	}

	// insert data
	if err := dac.InTransactionDo(ctx, insertData); err != nil {
		log.Fatal(err)
	}

	// query data

	// which permissions does user [john.doe] have?
	// Method 1: find john then find its permissions
	dac = dac.ForReadOnly(ctx)

	john := new(CloudIdentity)
	if err := dac.FindEquals(john, "user", "john.doe"); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", john)

	pip := new(PermissionsInProject)
	if err := dac.FindEquals(pip, "identity", john); err != nil {
		log.Fatal(err)
	}
	log.Printf("(with node) %#v", pip)

	{ // if you only have the uid of John
		pip := new(PermissionsInProject)
		if err := dac.FindEquals(pip, "identity", john.UID); err != nil {
			log.Fatal(err)
		}
		log.Printf("(uid only) %#v", pip)
	}

	// which permissions does user [john.doe] have?
	// Method 2: find permissions filtering groupOrUser predicate
	query := `{
		q(func: type(PermissionsInProject)) @cascade {
				identity @filter(eq(user,"john.doe"))
				permissions
		}
	  }`
	data := map[string][]string{}
	if err := dac.RunQuery(&data, query, "q"); err != nil {
		log.Fatal(err)
	}
	log.Printf("(filter predicate) %#v", data["permissions"])

	// for which projects has service account [compute-default] permissions [role/editor] ?
	query = `{
		q(func: type(PermissionsInProject)) @filter(eq(permissions,"role/editor")) {
			identity @filter(eq(serviceAccount,"compute-default"))
			project {                          
				project_name
			}
		}
	  }`
	data2 := map[string]interface{}{}
	if err := dac.RunQuery(&data2, query, "q"); err != nil {
		log.Fatal(err)
	}
	log.Printf("(filter account and permissions) %#v", data2["project"])
}

func insertData(xs *dga.DGraphAccess) error {
	// serviceAcount(compute-default) has permission(role/editor) in project(my-project)
	sa := &CloudIdentity{
		ServiceAccount: "compute-default",
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
	if err := xs.CreateEdge(pip, "identity", sa); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip.UID, "->", "serviceAccount", sa.UID)

	if err := xs.CreateEdge(pip, "project", pr); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip.UID, "->", "project", pr.UID)

	// user(john.doe) has permission(role/viewer) in project(my-project)
	pip2 := &PermissionsInProject{
		Permissions: []string{"role/viewer"},
	}
	ci := &CloudIdentity{
		User: "john.doe",
	}
	if err := xs.CreateEdge(pip2, "identity", ci); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip2.UID, "->", "user", ci.UID)

	if err := xs.CreateEdge(pip2, "project", pr); err != nil {
		return err
	}
	fmt.Println("permissions-in-project", pip2.UID, "->", "project", pr.UID)
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

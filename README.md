# dgraph-access

This is a helper package to work with `github.com/dgraph-io/dgo`, a Go client for accessing a DGraph cluster.

```
// Create a graph with project nodes and edges to each account within that project.
func projectsAndAccounts() {
	da := dga.NewDGraphAccess(newClient()).ForReadWrite()
	defer da.DiscardTransaction()

	for _, each := range listProjects() {
		puid, ok := da.FindNodeHasEqualsString("_Project", "name", each.Name)
		if !ok {
			p := &AnnotatedProject{Name: each.Name}
			err := da.CreateNode(p)
			if err != nil {
				log.Fatal(err)
			}
			puid = p.GetUID()
		}
		for _, other := range listAccounts(each.Name) {
			auid := ensureServiceAccount(da, &other).GetUID()
			err := da.CreateEdge(auid, "project", puid)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	da.CommitTransaction()
}
```
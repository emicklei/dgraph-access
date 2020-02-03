module github.com/emicklei/dgraph-access/example

go 1.13

require (
	github.com/dgraph-io/dgo/v2 v2.1.0
	github.com/emicklei/dgraph-access v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.25.1
)

replace github.com/emicklei/dgraph-access => ../../

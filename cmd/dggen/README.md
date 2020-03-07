# dggen

A Go code generator for structs from types defined in a DGraph schema.

## usage

    dggen -s schema.txt -p main > gen.go && go fmt gen.go

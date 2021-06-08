# Cloud Project Permissions

## run dgraph local

    docker run --rm -it -p 8000:8000 -p 8080:8080 -p 9080:9080 dgraph/standalone:latest

## run example

    go run *.go

## some Ratel queries

    {
        node as var(func: type(PermissionsInProject)) 
        node_uid(func: uid(node)){
            expand(PermissionsInProject) {
                project_name
                group
                user
                serviceAccount
            }
        }
    }
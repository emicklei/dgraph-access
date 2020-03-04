# dgraph-access

[![Build Status](https://travis-ci.org/emicklei/dgraph-access.png)](https://travis-ci.org/emicklei/dgraph-access)
[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/dgraph-access)](https://goreportcard.com/report/github.com/emicklei/dgraph-access)
[![GoDoc](https://godoc.org/github.com/emicklei/dgraph-access?status.svg)](https://pkg.go.dev/github.com/emicklei/dgraph-access?tab=doc)

This is a helper package to work with `github.com/dgraph-io/dgo` (v2), a Go client for accessing a DGraph cluster.
See the examples folder for complete programs.

## status

This package is under development (see commits); the API and scope may change before a v1.0.0 release.

## motivation

This package was created to reduce the boilerplate code required to use the `raw` dgraph Go client.
`dgraph-access` adds the following features to the standard Go client:

- type UID and NQuad to create RDF triples w/o facets
- type Node to encapsulate an uid and graph.type for your own entities
- type DgraphAccess to handle transactions, JSON marshalling and populating entities
- type Mutation to encapsulate a dgraph mutations that contains a list of RDF triples (NQuad values)
- DgraphAccess can trace the queries, mutations and responses for debugging
- DgraphAccess also provides a fluent interface for the operations

## usage

    import (
        dga "github.com/emicklei/dgraph-access"
    )

## example

See [documented code examples](https://godoc.org/github.com/emicklei/dgraph-access)

See [examples](https://github.com/emicklei/dgraph-access/blob/master/examples)

© 2019-2020, [ernestmicklei.com](http://ernestmicklei.com).  MIT License. Contributions welcome.
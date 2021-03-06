# dgraph-access

[![Build Status](https://travis-ci.org/emicklei/dgraph-access.png)](https://travis-ci.org/emicklei/dgraph-access)
[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/dgraph-access)](https://goreportcard.com/report/github.com/emicklei/dgraph-access)
[![GoDoc](https://godoc.org/github.com/emicklei/dgraph-access?status.svg)](https://pkg.go.dev/github.com/emicklei/dgraph-access?tab=doc)
[![codecov](https://codecov.io/gh/emicklei/dgraph-access/branch/master/graph/badge.svg)](https://codecov.io/gh/emicklei/dgraph-access)

This is a helper package to work with `github.com/dgraph-io/dgo` (v2), a Go client for accessing a DGraph cluster.
See the examples folder for complete programs.

## status

This package is under development (see commits); the API and scope may change before a v1.0.0 release.

## motivation

This package was created to reduce the boilerplate code required to use the `raw` dgraph Go client.

## features

`dgraph-access` adds the following features to the standard Go client:

- type UID and NQuad to create RDF triples w/o facets
- type Node to encapsulate an uid and graph.type for your own entities
- type DgraphAccess to handle transactions, JSON marshalling and populating entities
- type Mutation to encapsulate a dgraph mutations that contains a list of RDF triples (NQuad values)
- UpsertNode, CreateNode, CreateEdge, RunQuery, FindEquals model common dgraph operations
- DgraphAccess can trace the queries, mutations and responses for debugging
- DgraphAccess also provides a Service interface for convenient use of the operations

This repository also includes the [dggen tool](https://github.com/emicklei/dgraph-access/tree/master/cmd/dggen) that takes a dgraph schema to generate Go types from Dgraph types.

## usage

    import (
        dga "github.com/emicklei/dgraph-access"
    )

## quick look

    d := dga.NewDGraphAccess(yourDgraphClient).ForReadWrite()    
    s := d.Service()
    err := s.Alterschema(`name: string @index(exact) .`)
    type Vegetable struct {
        dga.Node `json:",inline" // this is for UID and Type propagation
        Name string
        Color string
    }
    v1 := &Vegetable{Name:"Carrot"}
    err = s.CreateNode(v1)

    v2 := &Vegetable{Name:"Beet"}
    err = s.CreateEdge(v1,"similarTo",v2)

    v3 := new(Vegetable)
    err := s.FindEquals(v3,"name","Carrot")    

    v2.Color = "darkred"
    s.UpsertNode(v2,"name","Beet")

## examples

See [examples](https://github.com/emicklei/dgraph-access/blob/master/examples).

See [documented code examples](https://godoc.org/github.com/emicklei/dgraph-access)

© 2019+, [ernestmicklei.com](http://ernestmicklei.com).  MIT License. Contributions welcome.

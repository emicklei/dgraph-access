# dgraph-access

[![Build Status](https://travis-ci.org/emicklei/dgraph-access.png)](https://travis-ci.org/emicklei/dgraph-access)
[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/dgraph-access)](https://goreportcard.com/report/github.com/emicklei/dgraph-access)
[![GoDoc](https://godoc.org/github.com/emicklei/dgraph-access?status.svg)](https://godoc.org/github.com/emicklei/dgraph-access)

This is a helper package to work with `github.com/dgraph-io/dgo` (v2), a Go client for accessing a DGraph cluster.
See the example folder for a complete program.

## status

This package is under heavy development; the API and scope may change before a v1.0.0 release.

## usage

    import (
        dga "github.com/emicklei/dgraph-access"
    )

## example UID

    dga.BlankUID("help")
    dga.StringUID("me")
    dga.IntegerUID(42)
    dga.FunctionUID("v")
    dga.NewUID("_:root")

Produces

    _:help
    <me>
    <0x2a>
    uid(v)
    _:root

## example NQuad

    salesCategoryID := "web1.4"
    assortmentID := 42
    nq := dga.BlankNQuad(salesCategoryID, "HAS_CATEGORY_ASSORTMENT", dga.BlankUID(assortmentID))
    nq.RDF()

Produces

    _:web1.4 <HAS_CATEGORY_ASSORTMENT> _:42 .

## example Mutation

    m := Mutation{
        Operation: "delete",
        NQuads: []NQuad{
            NQuad{
                Subject:   StringUID("0xf11168064b01135b"),
                Predicate: "died",
                Object:    1998},
        },
    }
    m.RDF()

Produces

    {
        delete {
            <0xf11168064b01135b> <died> 1998 .
        }
    }


## example Query and Mutate using the client

See [example](https://github.com/emicklei/dgraph-access/blob/master/example/main.go)

© 2019, [ernestmicklei.com](http://ernestmicklei.com).  MIT License. Contributions welcome.
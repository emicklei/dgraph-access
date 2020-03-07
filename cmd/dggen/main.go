package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	dsp "github.com/emicklei/dgraph-parser"
)

var (
	schemaFile  = flag.String("s", "dgraph.schema", "dgraph.schema file")
	packageName = flag.String("p", "main", "Go package name using in the models")
)

// go run *.go -s ../../examples/permissions/schema.txt > ../gen.go
func main() {
	flag.Parse()
	data, err := ioutil.ReadFile(*schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	parser := dsp.NewParser(bytes.NewReader(data))
	schema, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
	outputFile(schema)
}

func outputFile(s *dsp.Schema) {
	f := FileData{PackageName: *packageName}
	f.Imports = append(f.Imports, `dga "github.com/emicklei/dgraph-access"`)
	for _, each := range s.Types {
		t := TypeData{Name: each.Name}
		for _, other := range each.Predicates {
			def := s.FindPredicate(other.Name)
			if def != nil {
				v := FieldData{
					Name:           strings.Title(def.Name),
					TypeDefinition: toGoType(&f, def),
					Annotation:     toGoTags(def)}
				t.Fields = append(t.Fields, v)
			} else {
				log.Println("no definition found for", other.Name)
			}
		}
		f.Types = append(f.Types, t)
	}
	if err := fileTemplate.Execute(os.Stdout, f); err != nil {
		log.Fatal(err)
	}
}

func toGoType(f *FileData, p *dsp.PredicateDef) string {
	if p.Typename == "uid" {
		if p.IsArray {
			return "[]dga.Node"
		} else {
			return "dga.Node"
		}
	}
	if p.Typename == "dateTime" {
		f.EnsureImport("time")
		return "time.Time"
	}
	return p.Typename
}

func toGoTags(p *dsp.PredicateDef) string {
	extra := ""
	if p.Typename == "string" {
		extra = ",omitempty"
	}
	return fmt.Sprintf("`json:\"%s%s\"`", p.Name, extra)
}

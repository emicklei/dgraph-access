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
	outputFile  = flag.String("o", "models.go", "Go file name containing all the types")
)

// go run *.go -s ../../examples/permissions/schema.txt -o ../gen.go
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
	generate(schema, *outputFile)
}

func generate(s *dsp.Schema, output string) {
	out, err := os.Create(output)
	if err != nil {
		log.Fatalln("unable to create output:", err)
	}
	defer out.Close()

	f := FileData{PackageName: *packageName}
	f.Imports = append(f.Imports, `dga "github.com/emicklei/dgraph-access"`)
	for _, each := range s.Types {
		t := TypeData{Name: each.Name}
		for _, other := range each.Predicates {
			def := s.FindPredicate(other.Name)
			// unless no definition
			// OR is an edge i.o. scalar
			if def != nil {
				if !isEdge(def) {
					v := FieldData{
						Name:           fieldName(each.Name, def.Name),
						TypeDefinition: toGoType(&f, def),
						Annotation:     toGoTags(def)}
					t.Fields = append(t.Fields, v)
				} else {
					log.Printf("skip Go field for predicate to non-scalar object [%s]\n", other.Name)
				}
			} else {
				log.Println("no definition found for", other.Name)
			}
		}
		f.Types = append(f.Types, t)
	}
	if err := fileTemplate.Execute(out, f); err != nil {
		log.Fatal(err)
	}
}

func isEdge(p *dsp.PredicateDef) bool {
	return p.Typename == "uid"
}

func fieldName(t, s string) string {
	if len(s) <= 3 {
		return strings.ToUpper(s)
	}
	snakeless := strings.ReplaceAll(s, "_", "")
	if strings.HasPrefix(snakeless, strings.ToLower(t)) &&
		strings.HasSuffix(s, "id") {
		return "ID"
	}
	return strings.Title(s)
}

func toGoType(f *FileData, p *dsp.PredicateDef) string {
	// JSON mutations/queries are not handled for reference to Nodes
	// if p.Typename == "uid" {
	// 	if p.IsArray {
	// 		return "[]dga.Node"
	// 	} else {
	// 		return "dga.Node"
	// 	}
	// }
	switch p.Typename {
	case "default":
		if p.IsArray {
			return "[]string"
		}
		return "string"
	case "binary":
		return "[]byte"
	case "int":
		if p.IsArray {
			return "[]int64"
		}
		return "int64"
	case "float":
		if p.IsArray {
			return "[]float64"
		}
		return "float64"
	case "bool":
		if p.IsArray {
			return "[]bool"
		}
		return "bool"
	case "datetime":
		f.EnsureImport("time")
		if p.IsArray {
			return "[]time.Time"
		}
		return "time.Time"
	case "geo":
		return "string"
	case "uid":
		if p.IsArray {
			return "[]dga.UID"
		}
		return "dga.UID"
	case "string", "password":
		if p.IsArray {
			return "[]string"
		}
		return "string"
	default:
		if p.IsArray {
			return "[]string"
		}
		return "string"
	}
}

func toGoTags(p *dsp.PredicateDef) string {
	extra := ""
	if p.Typename == "string" {
		extra = ",omitempty"
	}
	return fmt.Sprintf("`json:\"%s%s\"`", p.Name, extra)
}

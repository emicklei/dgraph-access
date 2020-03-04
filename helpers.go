package dga

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type PredicateCondition struct {
	Predicate string
	Object    interface{}
}

func findFilterContent(predicateName, value interface{}) (filterContent string) {
	if s, ok := value.(string); ok {
		filterContent = fmt.Sprintf("eq(%s,%q)", predicateName, s)
	} else if n, ok := value.(HasUID); ok {
		filterContent = fmt.Sprintf("uid_in(%s,%s)", predicateName, n.GetUID().Assigned())
	} else if u, ok := value.(UID); ok {
		filterContent = fmt.Sprintf("uid_in(%s,%s)", predicateName, u.Assigned())
	} else {
		// unhandled type, TODO
		filterContent = fmt.Sprintf("eq(%s,%v)", predicateName, s)
	}
	return filterContent
}

func simpleType(result interface{}) string {
	tokens := strings.Split(fmt.Sprintf("%T", result), ".")
	return tokens[len(tokens)-1]
}

func trace(msg ...interface{}) {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "[dgraph-access-trace]")
	for _, each := range msg {
		fmt.Fprintf(b, " %v", each)
	}
	log.Println(b.String())
}

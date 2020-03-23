package dga

import (
	"reflect"
	"strings"
)

const thisPkg = "github.com/emicklei/dgraph-access"

func ReflectNQuads(uid UID, value HasUID) (list []NQuad) {
	e := reflect.ValueOf(value).Elem()
	for i := 0; i < e.NumField(); i++ {
		varType := e.Type().Field(i).Type
		// skip the embedded Node
		if varType.PkgPath() == thisPkg && varType.Name() == "Node" {
			continue
		}
		varName := e.Type().Field(i).Tag.Get("json")
		if len(varName) == 0 {
			varName = e.Type().Field(i).Name
		} else {
			namePlus := strings.Split(varName, ",")
			varName = namePlus[0]
		}
		varValue := e.Field(i).Interface()
		if IsZeroOfUnderlyingType(varValue) {
			continue
		}
		list = append(list, NQuad{
			Subject:   uid,
			Predicate: varName,
			Object:    varValue,
		})
	}
	return
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

package common

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func PrettyPrint(obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	Must(err)
	fmt.Printf("%s\n", jsonBytes)
}

func GetField(object interface{}, fieldName string) interface{} {
	reflectObject := reflect.Indirect(reflect.ValueOf(object))
	return reflectObject.FieldByName(fieldName).Interface()
}

package common

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// PrettyPrint displays a nested structure in human readable JSON format.
func PrettyPrint(obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	Must(err)
	fmt.Printf("%s\n", jsonBytes)
}

// GetField takes an arbitary structure, and uses reflection to retrieve the
// field with specified name from it. Panics if the field does not exist.
func GetField(object interface{}, fieldName string) interface{} {
	reflectObject := reflect.Indirect(reflect.ValueOf(object))
	return reflectObject.FieldByName(fieldName).Interface()
}

// StructDecode takes a source structure, and uses reflection to translate it to the
// destination structure. Typically used to convert native Go structures to
// map[string]interface{}, and vice versa.
func StructDecode(src interface{}, dest interface{}) error {
	if bytes, err := json.Marshal(src); err != nil {
		return err
	} else if err = json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

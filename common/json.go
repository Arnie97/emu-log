package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

var (
	confOnce sync.Once
	conf     map[string]string
)

// Conf takes a key and return the configuration value for it.
// The configuration will be loaded from a JSON file on the filesystem when
// the function is being called for the first time.
func Conf(key string) string {
	confOnce.Do(func() {
		if err := loadConf(); err != nil {
			Must(fmt.Errorf("cannot read conf file: %v", err))
		}
	})
	return conf[key]
}

func loadConf() error {
	file, err := os.Open(AppPath() + "/emu-log.json")
	if err != nil {
		return err
	}
	defer file.Close()

	if bytes, err := ioutil.ReadAll(file); err != nil {
		return err
	} else if err = json.Unmarshal(bytes, &conf); err != nil {
		return err
	}
	return nil
}

// AppPath returns the relative path for the directory in which the binary
// executable of this application resides.
func AppPath() string {
	path, err := os.Executable()
	Must(err)
	return filepath.Dir(path)
}

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

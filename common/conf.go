package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

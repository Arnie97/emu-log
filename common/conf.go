package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml"
)

const (
	ConfPath = "/emu-log.toml"
)

type (
	GlobalConf struct {
		Request  *RequestConf           `toml:"request,omitempty"`
		Schedule ScheduleConf           `toml:"schedule,omitempty"`
		Adapters map[string]AdapterConf `toml:"adapters,omitempty"`
	}
	AdapterConf struct {
		Request     *RequestConf     `toml:"request,omitempty"`
		SearchSpace []GenerationRule `toml:"search,omitempty"`
	}
	RequestConf struct {
		Interval  Duration `toml:"interval,omitempty"`
		UserAgent string   `toml:"user-agent,omitempty"`
		SessionID string   `toml:"session-id,omitempty"`
	}
	ScheduleConf struct {
		StartTime Duration `toml:"start-time,omitempty"`
		EndTime   Duration `toml:"end-time,omitempty"`
	}
	GenerationRule struct {
		Format string `toml:"format,omitempty"`
		Min    *int   `toml:"min,omitempty"`
		Max    *int   `toml:"max,omitempty"`
		Step   *int   `toml:"step,omitempty"`
	}
)

func (rule *GenerationRule) Emit(serials chan<- string) {
	defaultValue := 1
	if rule.Min == nil {
		rule.Min = &defaultValue
	}
	if rule.Max == nil {
		rule.Max = rule.Min
	}
	if rule.Step == nil {
		rule.Step = &defaultValue
	}

	for i := *rule.Min; i <= *rule.Max; i += *rule.Step {
		serials <- fmt.Sprintf(rule.Format, i)
	}
}

var (
	confOnce sync.Once
	conf     GlobalConf
)

// Conf loads configuration from a TOML file on the filesystem when
// the function is being called for the first time.
func Conf() *GlobalConf {
	confOnce.Do(func() {
		if err := loadConf(); err != nil {
			Must(fmt.Errorf("cannot read conf file: %v", err))
		}
	})
	return &conf
}

func loadConf() error {
	file, err := os.Open(AppPath() + ConfPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if bytes, err := ioutil.ReadAll(file); err != nil {
		return err
	} else if err = toml.Unmarshal(bytes, &conf); err != nil {
		return err
	}
	if conf.Request == nil {
		conf.Request = new(RequestConf)
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

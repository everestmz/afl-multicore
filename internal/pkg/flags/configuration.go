package flags

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

type MulticoreConfiguration struct {
	SessionName  string `flag:"session"`
	NumWorkers   int    `flag:"workers"`
	DisableCheck bool   `flag:"nocheck"`
}

func GenerateMulticoreConfig(args []string) (*MulticoreConfiguration, error) {
	config := defaultMulticoreConfiguration()

	err := generateConfig(args, config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func defaultMulticoreConfiguration() *MulticoreConfiguration {
	return &MulticoreConfiguration{
		NumWorkers:   runtime.NumCPU(),
		SessionName:  "fuzzer",
		DisableCheck: false,
	}
}

type AFLConfiguration struct {
	// Required parameters
	In     string   `flag:"i"`
	Out    string   `flag:"o"`
	Binary []string `flag:"-"`
	// Execution control settings
	File    string `flag:"f"`
	Timeout string `flag:"t"`
	Memory  string `flag:"m"`
	Qemu    bool   `flag:"Q"`
	// Fuzzing behavior settings
	Dirty bool   `flag:"d"`
	Dumb  bool   `flag:"n"`
	Dict  string `flag:"x"`
	// Other stuff
	Banner      string `flag:"T"`
	Exploration bool   `flag:"C"`
}

func GenerateAFLConfig(args []string) (*AFLConfiguration, error) {
	config := &AFLConfiguration{}

	for i, arg := range args {
		if arg == "--" {
			config.Binary = args[i+1:]
		}
	}

	err := generateConfig(args, config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (config *AFLConfiguration) GenerateArgs(flags []Flag) []string {
	args := []string{}

	t := reflect.TypeOf(*config)
	v := reflect.ValueOf(*config)
	for _, f := range flags {
		tag := t.Field(f.Int()).Tag.Get("flag")
		switch v.Field(f.Int()).Kind() {
		case reflect.Bool:
			if v.Field(f.Int()).Bool() {
				args = append(args, flagify(tag))
			}
		case reflect.String:
			val := v.Field(f.Int()).String()
			if val != "" {
				args = append(args, flagify(tag), val)
			}
		}
	}

	return args
}

func generateConfig(args []string, config interface{}) error {
	if !reflect.ValueOf(config).CanInterface() {
		return errors.New("'config' cannot be used as an interface")
	}

	flags := map[string]int{}
	t := reflect.Indirect(reflect.ValueOf(config)).Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("flag")
		if tag == "-" {
			continue
		}
		flags[tag] = i
	}

	v := reflect.ValueOf(config)
	s := v.Elem()
	for i, arg := range args {
		deflagged := deflag(arg)
		if idx, ok := flags[deflagged]; ok {
			switch s.Field(idx).Kind() {
			case reflect.Bool:
				s.Field(idx).SetBool(true)
			case reflect.String:
				s.Field(idx).SetString(args[i+1])
			case reflect.Int:
				num, err := strconv.ParseInt(args[i+1], 10, 0)
				if err != nil {
					return fmt.Errorf(
						"argument to '%s' flag must be an integer", deflagged)
				}
				s.Field(idx).SetInt(num)
			}
		}
	}

	return nil
}

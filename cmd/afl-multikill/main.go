package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/everestmz/afl-multicore/internal/pkg/flags"
	"github.com/everestmz/afl-multicore/internal/pkg/helpers"
)

func main() {
	helpers.CheckHelp(os.Args, helpText)

	config, _, err := flags.Parse(os.Args[1:])
	if err != nil {
		helpers.Fail("", err)
	}

	filename := filepath.Join(os.TempDir(), config.SessionName)
	_, err = os.Stat(filename)
	if err != nil {
		helpers.Fail(fmt.Sprintf(
			"could not find active session '%s'", config.SessionName), nil)
	}
	f, err := os.Open(filename)
	if err != nil {
		helpers.Fail("could not open active session", err)
	}
	fbytes, err := ioutil.ReadAll(f)
	if err != nil {
		helpers.Fail("could not read session file", err)
	}
	pidStrs := strings.Split(string(fbytes), "\n")

	for _, pidStr := range pidStrs {
		if pidStr == "" {
			continue
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			helpers.Fail("could not read pid from session file", err)
		}
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			helpers.Fail("could not kill process group", err)
		}
	}
	f.Close()
	err = os.Remove(filename)
	if err != nil {
		helpers.Fail("error removing session file", err)
	}
}

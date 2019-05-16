package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/everestmz/afl-multicore/internal/pkg/afl"
	"github.com/everestmz/afl-multicore/internal/pkg/flags"
	"github.com/everestmz/afl-multicore/internal/pkg/helpers"
)

func main() {
	helpers.CheckHelp(os.Args, helpText)

	if os.Getenv("AFL_PATH") == "" {
		helpers.Fail("please set the AFL_PATH variable", nil)
	}

	config, aflConfig, err := flags.Parse(os.Args[1:])
	if err != nil {
		helpers.Fail("", err)
	}

	err = afl.Showmap(aflConfig)
	if err != nil {
		helpers.Fail("", err)
	}

	// Create a session and save it for later
	filename := filepath.Join(os.TempDir(), config.SessionName)
	_, err = os.Stat(filename)
	if err == nil {
		helpers.Fail(fmt.Sprintf("session '%s' already exists", config.SessionName), nil)
	}
	f, err := os.Create(filename)
	if err != nil {
		helpers.Fail("could not create tempfile for session", err)
	}
	defer f.Close()

	// Set up AFl workers
	workers := []*afl.WorkerOptions{{
		Master: true,
		Name:   fmt.Sprintf("%s_master", config.SessionName),
	}}
	for i := 1; i < config.NumWorkers; i++ {
		workers = append(workers, &afl.WorkerOptions{
			Name: fmt.Sprintf("%s_slave%v", config.SessionName, i),
		})
	}

	for _, worker := range workers {
		pid, err := afl.Fuzz(aflConfig, worker)
		if err != nil {
			helpers.Fail("", err)
		}

		_, err = f.WriteString(fmt.Sprintf("%s\n", strconv.Itoa(pid)))
		if err != nil {
			helpers.Fail("could not save session", err)
		}
	}

	fmt.Printf("%v AFl workers running! Use afl-multistats to check progress.\n", len(workers))
}

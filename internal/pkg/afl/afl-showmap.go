package afl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/everestmz/afl-multicore/internal/pkg/flags"
)

const (
	aflShowmap = "afl-showmap"
)

func Showmap(config *flags.AFLConfiguration) error {
	// First, create the tempfile we're testing with
	inputFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return fmt.Errorf("Could not create tempfile: %s", err.Error())
	}
	defer os.Remove(inputFile.Name())

	// Write a simple 1-char test
	_, err = inputFile.WriteString("z")
	if err != nil {
		return fmt.Errorf("could not write test file: %s", err.Error())
	}
	err = inputFile.Close()
	if err != nil {
		return fmt.Errorf("could not write test file: %s", err.Error())
	}

	// Then, create the file where we'll place the map
	mapFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return fmt.Errorf("Could not create tempfile: %s", err.Error())
	}
	defer os.Remove(mapFile.Name())

	// Finally, compile the args needed
	args := []string{"-o", mapFile.Name()}
	args = append(args, config.GenerateArgs([]flags.Flag{
		flags.Timeout, flags.Memory, flags.QEMU},
	)...)

	// If the binary has a @@ we need to replace it
	binary := []string{"--"}
	for _, arg := range config.Binary {
		if strings.Contains(arg, "@@") {
			binary = append(binary, strings.Replace(arg, "@@", inputFile.Name(), -1))
		} else {
			binary = append(binary, arg)
		}
	}

	args = append(args, binary...)

	cmd := exec.Command(aflShowmap, args...)
	outputBuf := &bytes.Buffer{}
	cmd.Stderr = outputBuf
	cmd.Stdout = outputBuf
	err = cmd.Run()
	if err != nil {
		fmt.Printf(outputBuf.String())
		return fmt.Errorf("afl-showmap failed: %s", err.Error())
	}
	return nil
}

func isSet(s string) bool {
	return s != ""
}

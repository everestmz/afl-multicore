package afl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/everestmz/afl-multicore/internal/pkg/flags"
)

const (
	aflFuzz = "afl-fuzz"
)

type WorkerOptions struct {
	Master bool
	Name   string
}

func (o WorkerOptions) Args() []string {
	var args []string
	if o.Master {
		args = []string{"-M"}
	} else {
		args = []string{"-S"}
	}
	return append(args, o.Name)
}

func Fuzz(config *flags.AFLConfiguration, options *WorkerOptions) (int, error) {
	args := config.GenerateArgs([]flags.Flag{
		flags.In, flags.Out, flags.File, flags.Timeout, flags.Memory, flags.QEMU,
		flags.Dirty, flags.Dumb, flags.Dict, flags.Banner, flags.Exploration,
	})
	// Parse args

	binary := append([]string{"--"}, config.Binary...)
	args = append(
		append(args, options.Args()...),
		binary...)

	os.Setenv("AFL_NO_UI", "1")
	cmd := exec.Command(aflFuzz, args...)
	outputBuf := &bytes.Buffer{}
	cmd.Stderr = outputBuf
	cmd.Stdout = outputBuf
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := cmd.Start()
	if err != nil {
		fmt.Printf(outputBuf.String())
		return 0, fmt.Errorf("afl-fuzz failed: %s", err.Error())
	}

	ticker := time.NewTicker(time.Second)
	stopChan := make(chan error)
	go func() {
		stopChan <- cmd.Wait()
	}()

	var output strings.Builder

	for {
		select {
		case <-ticker.C:
			bytes, _ := ioutil.ReadAll(outputBuf)
			output.Write(bytes)
			outStr := string(bytes)
			if strings.Contains(outStr, "All set and ready to roll!") {
				return cmd.Process.Pid, nil
			}
		case err := <-stopChan:
			bytes, _ := ioutil.ReadAll(outputBuf)
			output.Write(bytes)
			fmt.Printf(output.String())
			return 0, err
		}
	}
}

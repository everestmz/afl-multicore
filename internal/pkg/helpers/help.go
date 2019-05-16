package helpers

import (
	"fmt"
	"os"
)

func CheckHelp(args []string, text string) {
	if isHelp(args) {
		fmt.Println(text)
		os.Exit(0)
	}
}

func isHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		} else if len(args) == 1 {
			return true
		}
	}
	return false
}

func Fail(s string, err error) {
	os.Stderr.WriteString(s)
	if err != nil {
		if s != "" {
			os.Stderr.WriteString(": ")
		}
		os.Stderr.WriteString(err.Error())
	}
	os.Stderr.WriteString("\n")
	os.Exit(1)
}

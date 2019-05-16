package flags

import (
	"fmt"
	"strings"
)

func flagify(s string) string {
	return fmt.Sprintf("-%s", s)
}

func deflag(s string) string {
	if strings.HasPrefix(s, "--") {
		return strings.Replace(s, "--", "", 1)
	} else if strings.HasPrefix(s, "-") {
		return strings.Replace(s, "-", "", 1)
	}
	return s
}

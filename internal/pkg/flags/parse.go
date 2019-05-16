// We can't use the default flags.parse, since we want "unknown"
// flags to be passed through to AFL and only extract the flags we
// care about.

package flags

import (
	"fmt"
)

func Parse(args []string) (*MulticoreConfiguration, *AFLConfiguration, error) {

	// Pre-parse checks
	for _, elt := range args {
		switch elt {
		case "-M":
			return nil, nil, fmt.Errorf(
				"do not use '-M' when invoking afl-multicore")
		case "-S":
			return nil, nil, fmt.Errorf(
				"do not use '-S' when invoking afl-multicore")
		default:
			continue
		}
	}

	multiConfig, err := GenerateMulticoreConfig(args)
	if err != nil {
		return nil, nil, err
	}

	aflConfig, err := GenerateAFLConfig(args)
	if err != nil {
		return nil, nil, err
	}

	return multiConfig, aflConfig, nil
}

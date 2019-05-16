package main

import (
	"flag"
	"fmt"
	"os"

	multistats "github.com/everestmz/afl-multicore/internal/afl-multistats"
	"github.com/everestmz/afl-multicore/internal/pkg/helpers"
)

var (
	format   = flag.String("format", "json", "the format to output stats in")
	advanced = flag.Bool("advanced", false, "output all stats, not just summary")
	hud      = flag.Bool("hud", false, "display a persistent stats HUD")
)

func main() {
	helpers.CheckHelp(os.Args, helpText)

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		helpers.Fail("please ensure AFL sync dir is the only argument", nil)
	}

	syncDir := args[0]
	if *hud {
		multistats.Hud(syncDir)
		return
	}

	instanceStats, err := multistats.ReadStats(syncDir)
	if err != nil {
		helpers.Fail("could not read stats", err)
	}

	finalStats := multistats.MergeStats(instanceStats)

	// Now, output stats in some way
	if *advanced {
		bytes, err := finalStats.JSON()
		if err != nil {
			helpers.Fail("error serializing stats to JSON", err)
		}
		fmt.Printf(string(bytes))
		return
	}

	switch *format {
	case "json":
		bytes, err := finalStats.Basic().JSON()
		if err != nil {
			helpers.Fail("error serializing stats to JSON", err)
		}
		fmt.Printf(string(bytes))
	case "human":
		fmt.Printf(string(finalStats.Basic().Human()))
	}

}

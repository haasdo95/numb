package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/numb/run"

	"github.com/user/numb/bootstrap"
)

func printUsage() {
	fmt.Println("Wrong Usage")
}

func main() {

	args := os.Args

	if len(args) == 1 {
		printUsage()
		return
	}

	// Parsing cmdline arguments
	skipFirstFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	silent := skipFirstFlagSet.Bool("silent", false, "silent stdout & stderr")
	skipFirstFlagSet.Parse(os.Args[2:])
	var runconfig = map[string]interface{}{
		"silent": *silent,
	}

	subcmd := args[1] // the second
	switch subcmd {   // cannot assume that nmb.json exists
	case "init":
		bootstrap.Init()
	case "deinit":
		bootstrap.Deinit()
	default:
		nmbConfig := bootstrap.GetConfig()

		switch subcmd {
		case "test":
			for _, cmdline := range nmbConfig.Test {
				run.Test(cmdline, runconfig)
			}
		case "train":
			for _, cmdline := range nmbConfig.Train {
				run.Train(cmdline, runconfig)
			}
		}
	}
}

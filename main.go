package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/numb/run"

	"github.com/user/numb/bootstrap"
	"github.com/user/numb/analysis"
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
		collection, session := run.GetCollection()
		defer session.Close()
		switch subcmd {
		case "test":
			for _, cmdline := range nmbConfig.Test {
				run.Test(cmdline, runconfig, collection)
			}
		case "train":
			for _, cmdline := range nmbConfig.Train {
				run.Train(cmdline, runconfig, collection)
			}
		case "list":
			analysis.List(collection)
		default:
			printUsage()
		}
	}
}

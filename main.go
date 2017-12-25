package main

import (
	"gopkg.in/mgo.v2"
	"flag"
	"fmt"
	"os"

	"github.com/nasyxx/numb/run"
	"github.com/nasyxx/numb/utils"
	"github.com/nasyxx/numb/database"

	"github.com/nasyxx/numb/bootstrap"
	"github.com/nasyxx/numb/analysis"
	"github.com/nasyxx/numb/versioning"
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
	all := skipFirstFlagSet.Bool("all", false, "test all that is untested")
	skipFirstFlagSet.Parse(os.Args[2:])
	var runconfig = map[string]interface{}{
		"silent": *silent,
		"all": *all,
	}

	subcmd := args[1] // the second
	switch subcmd {   // cannot assume that nmb.json exists
	case "init":
		bootstrap.Init()
	case "deinit":
		bootstrap.Deinit()
	default:
		nmbConfig := bootstrap.GetConfig()
		session, err := mgo.Dial("127.0.0.1")
		utils.Check(err)
		collection := database.GetCollection(session)
		defer session.Close()
		switch subcmd {
		case "test":
			for _, cmdline := range nmbConfig.Test {
				run.Test(cmdline, runconfig, collection, nil)
			}
		case "train":
			for _, cmdline := range nmbConfig.Train {
				run.Train(cmdline, runconfig, collection)
			}
		case "list":
			analysis.List(collection)
		case "revert":
			// anything trailing
			if len(args) == 3 {
				versioning.Revert(collection, args[2])
			} else if len(args) == 2 { // nothing trailing
				versioning.Revert(collection, "")
			}
		case "queue":
			if len(args) == 4 {
				if args[2] == "init" {
					run.QueueInit(args[3])
				} else if args[2] == "run" {
					for _, cmdline := range nmbConfig.Train {
						run.QueueRun(cmdline, runconfig, collection, args[3])
					}
				} else {
					printUsage()
				}
			} else {
				printUsage()
			}
		case "report":
			if len(args) == 3 {
				analysis.Report(collection, args[2])
			} else {
				printUsage()
			}
		default:
			printUsage()
		}
	}
}

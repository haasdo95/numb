package run

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"github.com/user/numb/utils"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Train runs a command in train mode.
func Train(cmdline string, runconfig map[string]interface{}) {
	trainEnv := make(map[string]string)
	trainEnv["NUMB_MODE"] = "TRAIN"
	run(cmdline, trainEnv, runconfig, true)
}

// Test runs a command in train mode.
func Test(cmdline string, runconfig map[string]interface{}) {
	testEnv := make(map[string]string)
	testEnv["NUMB_MODE"] = "TEST"
	run(cmdline, testEnv, runconfig, false)
}

func runTrain(cmd *exec.Cmd, graphReader, paramReader, stateDictReader *os.File, collection *mgo.Collection) {
	utils.Check(cmd.Wait())

}

func runTest(cmd *exec.Cmd, graphReader, interactWriter *os.File, collection *mgo.Collection) {
	// capture signal
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGUSR1)

	// get comp graph first
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, graphReader)
	utils.Check(err)
	concreteGraph := string(buf.Bytes())
	query := collection.Find(bson.M{"ConcreteGraph": concreteGraph}).Sort("-Timestamp")
	if cnt, _ := query.Count(); cnt == 0 {
		cmd.Process.Signal(os.Interrupt)
		fmt.Println("The model you are testing doesn't even exist")
		return
	}

	<-sigs // blocks until signal comes

	// start prompting user:
	fmt.Println("This model has been trained with following parameters.")
	fmt.Println("Simply hit enter to use the latest. Or input the number to specify.")

	results := make([]Schema, 0)
	utils.Check(query.All(&results))
	for idx, r := range results {
		fmt.Printf("%d: %v", idx, r.Params)
	}

	fmt.Print("Use parameter: ")
	var choice = 0
	fmt.Scanln(&choice)

	savedStatedictFilename := results[choice].StateDictFilename
	savedStatedictFilename = ".nmb/" + savedStatedictFilename

	interactWriter.WriteString(savedStatedictFilename) // send the file name to python

	utils.Check(cmd.Wait())
}

func run(cmdline string, newEnv map[string]string, runconfig map[string]interface{}, isTrain bool) {
	cmdPath := strings.Split(cmdline, " ")
	cmd := exec.Command(cmdPath[0], cmdPath[1:]...)

	// set runtime config
	var silent = false
	if s1, keyok := runconfig["silent"]; keyok {
		if s2, typeok := s1.(bool); typeok {
			silent = s2
		}
	}
	if !silent {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// end: set runtime config

	oldEnv := os.Environ()
	strNewEnv := utils.Map2env(newEnv)
	env := append(oldEnv, strNewEnv...)
	cmd.Env = env

	// create pipes
	pGraphR, pGraphW, err := os.Pipe()
	utils.Check(err)
	defer pGraphR.Close()

	pParamR, pParamW, err := os.Pipe()
	utils.Check(err)
	defer pParamR.Close()

	pStateR, pStateW, err := os.Pipe()
	utils.Check(err)
	defer pStateR.Close()

	pInteractR, pInteractW, err := os.Pipe() // This is a writing pipe!
	utils.Check(err)
	defer pInteractW.Close()
	// end: create pipes

	// setting pipes in py script
	cmd.ExtraFiles = []*os.File{
		pGraphW,
		pParamW,
		pStateW,
		pInteractR, // will block python execution;
		// In python script a signal will be sent to the parent before blocking
	}

	collection, session := GetCollection()
	defer session.Close()

	utils.Check(cmd.Start())

	pGraphW.Close()
	pParamW.Close()
	pStateW.Close()
	pInteractR.Close()

	if !isTrain {
		runTest(cmd, pGraphR, pInteractW, collection)
	} else {
		runTrain(cmd, pGraphR, pParamR, pStateR, collection)
	}

}

// if !isTrain { // block for user to input stuff
// 	<-sigDone
// }

// pGraphW.Close() // keep it otherwise io.Copy will block indefinitely
// pParamW.Close()
// pStateW.Close()

// fmt.Println("Printing: Graph")
// buf := bytes.NewBuffer(nil)
// _, err = io.Copy(buf, pGraphR)
// utils.Check(err)
// fmt.Println(string(buf.Bytes()))

// fmt.Println("Printing: Param")
// buf = bytes.NewBuffer(nil)
// _, err = io.Copy(buf, pParamR)
// utils.Check(err)
// fmt.Println(string(buf.Bytes()))

// fmt.Println("Printing: State Dict")
// buf = bytes.NewBuffer(nil)
// size, err := io.Copy(buf, pStateR)
// utils.Check(err)
// fmt.Println("NUMBER OF BYTES RECEIVED AS STATE DICT: ", size)

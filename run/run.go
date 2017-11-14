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

	// capture signal
	sigs := make(chan os.Signal)
	sigDone := make(chan bool)
	if !isTrain { // this is only necessary when the user is running test mode
		signal.Notify(sigs, syscall.SIGUSR1)
		go func() {
			sig := <-sigs // blocks until signal comes
			// TODO: handle user input
			fmt.Println("Receiving: ", sig)
			pInteractW.WriteString("Hello From Dad!")
			pInteractW.Close() // let the Python script continue
			sigDone <- true    // let golang continue (not really necessary for now)
		}()
	}

	utils.Check(cmd.Start())

	if !isTrain { // block for user to input stuff
		<-sigDone
	}

	pGraphW.Close() // keep it otherwise io.Copy will block indefinitely
	pParamW.Close()
	pStateW.Close()

	fmt.Println("Printing: Graph")
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, pGraphR)
	utils.Check(err)
	fmt.Println(string(buf.Bytes()))

	fmt.Println("Printing: Param")
	buf = bytes.NewBuffer(nil)
	_, err = io.Copy(buf, pParamR)
	utils.Check(err)
	fmt.Println(string(buf.Bytes()))

	fmt.Println("Printing: State Dict")
	buf = bytes.NewBuffer(nil)
	size, err := io.Copy(buf, pStateR)
	utils.Check(err)
	fmt.Println("NUMBER OF BYTES RECEIVED AS STATE DICT: ", size)

	utils.Check(cmd.Wait())

}

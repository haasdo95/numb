package run

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

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
	run(cmdline, trainEnv, runconfig)
}

// Test runs a command in train mode.
func Test(cmdline string, runconfig map[string]interface{}) {
	testEnv := make(map[string]string)
	testEnv["NUMB_MODE"] = "TEST"
	run(cmdline, testEnv, runconfig)
}

func run(cmdline string, newEnv map[string]string, runconfig map[string]interface{}) {
	cmdPath := strings.Split(cmdline, " ")
	cmd := exec.Command(cmdPath[0], cmdPath[1:]...)

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

	oldEnv := os.Environ()
	strNewEnv := utils.Map2env(newEnv)
	env := append(oldEnv, strNewEnv...)
	cmd.Env = env

	pGraphR, pGraphW, err := os.Pipe()
	check(err)
	defer pGraphR.Close()

	pParamR, pParamW, err := os.Pipe()
	check(err)
	defer pParamR.Close()

	pStateR, pStateW, err := os.Pipe()
	check(err)
	defer pStateR.Close()

	cmd.ExtraFiles = []*os.File{
		pGraphW,
		pParamW,
		pStateW,
	}

	check(cmd.Start())

	pGraphW.Close() // keep it otherwise io.Copy will block indefinitely
	pParamW.Close()
	pStateW.Close()

	fmt.Println("Printing: Graph")
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, pGraphR)
	check(err)
	fmt.Println(string(buf.Bytes()))

	fmt.Println("Printing: Param")
	buf = bytes.NewBuffer(nil)
	_, err = io.Copy(buf, pParamR)
	check(err)
	fmt.Println(string(buf.Bytes()))

	fmt.Println("Printing: State Dict")
	buf = bytes.NewBuffer(nil)
	size, err := io.Copy(buf, pStateR)
	check(err)
	fmt.Println("NUMBER OF BYTES RECEIVED AS STATE DICT: ", size)

	check(cmd.Wait())

}

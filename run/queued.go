package run

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"github.com/nasyxx/numb/utils"
)

type QueueSpec struct {
	Queue []map[string]interface{} `json:"queue"`;
	Start int `json:"start"`;
	End int `json:"end"`;
}

// QueueInit creates the queue json file named by qfileName
func QueueInit(qfileName string) {
	qs := QueueSpec {
		Queue: []map[string]interface{}{},
		Start: 0,
		End: 0,
	}
	qsBytes, err := json.MarshalIndent(qs, "", "\t")
	utils.Check(err)
	qsFile, err := os.OpenFile(qfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer qsFile.Close()
	utils.Check(err)
	_, err = qsFile.Write(qsBytes)
	utils.Check(err)
}

func QueueRun(cmdline string, runconfig map[string]interface{}, collection *mgo.Collection, qfileName string) {
	var qs QueueSpec
	raw, err := ioutil.ReadFile(qfileName)
	utils.Check(err)
	err = json.Unmarshal(raw, &qs)
	utils.Check(err)
	queueEnv := make(map[string]string)
	queueEnv["NUMB_MODE"] = "QUEUE"
	// check if range is valid
	if qs.Start >= 0 && qs.Start < qs.End && qs.End <= len(qs.Queue) {
		for i := qs.Start; i < qs.End; i++ {
			entry := qs.Queue[i]
			jb, err := json.Marshal(entry)
			utils.Check(err)
			js := string(jb)
			run(cmdline, queueEnv, runconfig, QUEUE, collection, js, nil)
		}
	} else {
		println("Bad Range")
	}
}

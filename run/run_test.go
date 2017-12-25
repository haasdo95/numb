package run_test

import (
	"io/ioutil"
	"strings"
	"io"
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"github.com/haasdo95/numb/database"
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
	"log"
	"gopkg.in/ory-am/dockertest.v3"
	"github.com/haasdo95/numb/run"
	"github.com/libgit2/git2go"
	
)

var db *mgo.Session

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalln("Can't connect to docker")
	}
	resource, err := pool.Run("sameersbn/mongodb", "latest", []string{})
	if err != nil {
		log.Fatalln("Can't get resource")
	}
	if err := pool.Retry(func () error {
		var err error
		db, err = mgo.Dial(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	err = os.Chdir("github.com/haasdo95/numb/demo")
	if err != nil {
		log.Fatal(err.Error())
	}
	code := m.Run() // actually run testcases
	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf(err.Error())
	}
	os.Exit(code)
}

func isInNmb(result database.Schema, t *testing.T) {
	// test that the statedict file has been checked in .nmb
	files, _ := ioutil.ReadDir(".nmb")
	notFound := true
	for _, file := range files {
		if file.Name() == result.StateDictFilename {
			notFound = false
			break
		}
	}
	if notFound {
		t.Fatalf("State Dict File named %s is not found under .nmb", result.StateDictFilename)
	}
}

func TestTrain(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()	
	run.Train("python train.py", map[string]interface{}{
		"silent": true,
	}, collection)
	query := collection.Find(bson.M{})
	if cnt, err := query.Count(); cnt != 1 || err != nil {
		t.Fatal("Fail to retrieve only one result")
	}
	var results []database.Schema
	query.All(&results)
	result := results[0]
	// get HEAD of numb branch
	repo, err := git.OpenRepository(".git")
	numbBranch, err := repo.LookupBranch("numb", git.BranchLocal)
	if err != nil {
		t.Fatal("Error occurred when opening repo!")
	}
	head := numbBranch.Target()
	if head.String() != result.Versioning {
		t.Fatalf("Git hashed don't agree! numb branch has: %s; db has: %s", head.String(), result.Versioning)
	}
	// test that code snapshot has been successfully retrieved
	trainFile, err := os.Open("lenet.py")
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, trainFile)
	fileContent := buffer.String()
	if !strings.Contains(fileContent, result.Code) {
		t.Fatalf("DB entry & actual file don't match! File should contain:\n %s", result.Code)
	}
	// test that the statedict file has been checked in .nmb
	isInNmb(result, t)
}

func TestTest(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()
	defaultSetting := map[string]interface{}{
		"silent": false,
		"all": false,
	}
	run.Train("python train.py", defaultSetting, collection) // train it first
	mockInput := strings.NewReader("0\n") // inject a mock reader
	run.Test("python test.py", defaultSetting, collection, mockInput)
	// see if test result has been saved in database
	query := collection.Find(bson.M{})
	if cnt, err := query.Count(); cnt != 1 || err != nil {
		t.Fatal("Bad result count")
	}
	var results []database.Schema
	query.All(&results)
	result := results[0]
	if result.Test == "" {
		t.Fatal("Testing results weren't saved!")
	}
}

func TestQueue(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()
	run.QueueRun("python train.py", map[string]interface{}{}, collection, "testqueue.json")
	query := collection.Find(bson.M{})
	if cnt, err := query.Count(); cnt != 3 || err != nil {
		t.Fatal("Fail to retrieve exactly three results")
	}
	var results []database.Schema
	query.All(&results)
	// test that all state dicts are saved
	for _, result := range results {
		isInNmb(result, t)
	}
}

func TestTestAll(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()
	run.QueueRun("python train.py", map[string]interface{}{}, collection, "testqueue.json")
	query := collection.Find(bson.M{})
	if cnt, err := query.Count(); cnt != 3 || err != nil {
		t.Fatal("Fail to retrieve exactly three results")
	}
	var results []database.Schema
	query.All(&results)
	for _, result := range results {
		if result.Test != "" {
			t.Fatalf("Should be untested. But test result is set as: %s", result.Test)
		}
	}
	testAllSetting := map[string]interface{}{
		"silent": false,
		"all": true,
	}
	run.Test("python test.py", testAllSetting, collection, nil)
	// test if all three are filled out with test results
	query = collection.Find(bson.M{}) // need to re-query
	query.All(&results)
	if len(results) != 3 {
		t.Fatal("Wrong number of parameter record")
	}
	for _, result := range results {
		if result.Test == "" {
			t.Fatal("Test result is not filled in!")
		}
	}
}
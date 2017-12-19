package analysis_test

import (
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"io"
	"bytes"
	"github.com/nasyxx/numb/run"
	"github.com/nasyxx/numb/analysis"
	"github.com/nasyxx/numb/database"
	"regexp"
	"gopkg.in/ory-am/dockertest.v3"
	"gopkg.in/mgo.v2"
	"testing"
	"log"
	"fmt"
	"os"
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
	err = os.Chdir("github.com/nasyxx/numb/demo")
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

func maskID(benchmark string) string {
	var regexID = regexp.MustCompile(`ID:  ([0-9]+)`)
	return string(regexID.ReplaceAll([]byte(benchmark), []byte("")))
}

func TestList(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()
	run.QueueRun("python train.py", map[string]interface{}{}, collection, "testqueue.json")
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Failed to create pipe")
	}
	os.Stdout = w
	analysis.List(collection)
	w.Close()
	os.Stdout = origStdout
	// retrieve output from r
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r)
	output := buf.String()
	output = maskID(output)
	// test that output is the same as benchmark
	benchmark := openBenchmark("benchmark_beforetest.txt", t)
	benchmark = maskID(benchmark)
	if output != benchmark {
		t.Fatal("Doesn't match 'benchmark-before'")
	}
	r.Close()

	// see if test results can be reflected in numb list
	run.Test("python test.py", map[string]interface{}{
		"all": true,
	}, collection, nil)
	origStdout = os.Stdout
	r, w, err = os.Pipe()
	if err != nil {
		t.Fatal("Failed to create pipe")
	}
	os.Stdout = w
	analysis.List(collection)
	w.Close()
	os.Stdout = origStdout
	// retrieve output from r
	buf = bytes.NewBuffer(nil)
	io.Copy(buf, r)
	output = buf.String()
	output = maskID(output)
	r.Close()
	benchmark = openBenchmark("benchmark_aftertest.txt", t)
	benchmark = maskID(benchmark)
	if output != benchmark {
		t.Fatal("Doesn't match 'benchmark-after'")
	}
}

func openBenchmark(benchmarkName string, t *testing.T) string {
	benchmarkBeforeFile, err := os.Open(benchmarkName)
	if err != nil {
		t.Fatal("Failed to open benchmark file")
	}
	benchmarkBuffer := bytes.NewBuffer(nil)
	io.Copy(benchmarkBuffer, benchmarkBeforeFile)
	return benchmarkBuffer.String()
}

func TestReport(t *testing.T) {
	collection := database.GetCollection(db)
	defer collection.DropCollection()
	run.QueueRun("python train.py", map[string]interface{}{}, collection, "testqueue.json")
	run.Test("python test.py", map[string]interface{}{
		"all": true,
	}, collection, nil)
	query := collection.Find(bson.M{})
	var result database.Schema
	query.One(&result)
	ID := strconv.FormatInt(result.Timestamp, 10)
	t.Log(ID)
	wd, _ := os.Getwd()
	t.Log(wd)
	analysis.Report(collection, ID)
	// test if the dir has been created
	dirname := "report-" + ID
	_, err := os.Stat(wd + "/" + dirname)
	if os.IsNotExist(err) {
		t.Fatal("Dir not created!")
	}
	os.RemoveAll(wd + "/" + dirname)
}
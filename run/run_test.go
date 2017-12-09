package run_test

import (
	"os"
	"fmt"
	"gopkg.in/mgo.v2"
	// "os"
	"testing"
	"log"
	"gopkg.in/ory-am/dockertest.v3"
	"github.com/user/numb/run"
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
	code := m.Run() // actually run testcases
	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf(err.Error())
	}
	os.Exit(code)
}

func TestTrain(t *testing.T) {
	os.Chdir("testcases")
	run.Train("python train.py")
}

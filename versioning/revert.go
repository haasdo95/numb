package versioning

import (
	"github.com/user/numb/database"
	"strconv"
	"os"
	"os/exec"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"gopkg.in/mgo.v2"
)

// Revert allows you to get back to a previous stage
func Revert(collection *mgo.Collection, timestamp string) {
	if timestamp == "" {
		exec.Command("git", "checkout", "master").Run()
		return
	}
	timestampNum, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse input")
		return
	}
	query := collection.Find(bson.M{"timestamp": timestampNum})
	if cnt, _ := query.Count(); cnt != 1 {
		fmt.Println("Cannot find record")
		return
	}
	var result database.Schema
	query.One(&result)
	if result.Versioning == "" {
		fmt.Println("This model is not versioned.")
		fmt.Println("Probably it's generated with 'numb queue run'")
		return
	}
	gitHash := result.Versioning
	cmd := exec.Command("git", "checkout", gitHash)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	fmt.Println("Dear User:")
	fmt.Println("\tPlease try not to edit stuff here,")
	fmt.Println("\tbefore I figure out how to let you do so safely.")
	fmt.Println("A million thanks.")
}
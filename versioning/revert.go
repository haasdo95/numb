package versioning

import (
	"os/exec"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"gopkg.in/mgo.v2"
)

// Revert allows you to get back to a previous stage
func Revert(collection *mgo.Collection, gitHash string) {
	if gitHash == "" {
		exec.Command("git", "checkout", "master").Run()
		return
	}
	query := collection.Find(bson.M{"versioning": gitHash})
	if cnt, _ := query.Count(); cnt != 1 {
		fmt.Println("Cannot find record")
		return
	}
	exec.Command("git", "checkout", gitHash).Run()
	fmt.Println("Dear User:")
	fmt.Println("\tPlease try not to edit stuff here,")
	fmt.Println("\tbefore I figure out how to let you do so safely.")
	fmt.Println("A million thanks.")
}
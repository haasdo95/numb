package run

import (
	"github.com/user/numb/bootstrap"
	"github.com/user/numb/utils"
	"gopkg.in/mgo.v2"
)

func GetCollection() *mgo.Collection {
	session, err := mgo.Dial("127.0.0.1")
	utils.Check(err)
	session.SetMode(mgo.Monotonic, true)
	name := bootstrap.GetConfig().Name
	collection := session.DB(name).C(name)

	index := mgo.Index{
		Key: []string{
			"AbstractGraph",
			"ConcreteGraph",
			"StateDictFilename",
			"Params",
			"Test",
		},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	utils.Check(collection.EnsureIndex(index))

	return collection
}

package database

import (
	"github.com/haasdo95/numb//bootstrap"
	"github.com/haasdo95/numb//utils"
	"gopkg.in/mgo.v2"
)

func GetCollection(session *mgo.Session) (*mgo.Collection) {
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
			"Timestamp",
		},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	utils.Check(collection.EnsureIndex(index))

	return collection
}

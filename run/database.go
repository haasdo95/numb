// this file defines the database schema
package run

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Schema struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	AbstractGraph     string
	ConcreteGraph     string
	StateDictFilename string
	Params            map[string]interface{}
	Test              TestInfo
	Timestamp         time.Time
}

type TestInfo struct {
	Accuracy float64
	Extra    map[string]interface{}
	// TODO: more test records
}

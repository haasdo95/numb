// Package analysis provides the functionality to do numb list and numb insight
package analysis

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

const (
	// BOTH mode lists both tested and untested
	BOTH ListMode = iota
	// TESTED mode lists only tested model
	TESTED
	// SHORT mode lists 5 at most
	SHORT
	// REVERSE mode puts the oldest first
	REVERSE
)
type ListMode int
type ListModes []ListMode

func (modes ListModes) containMode(mode ListMode) bool {
	for _, currMode := range modes {
		if currMode == mode {
			return true
		}
	}
	return false
}

func groupBy(collection *mgo.Collection, key string) []map[string]interface{} {
	pipeline := []bson.M {
		bson.M {
			"$match": bson.M {}, // match all
		},
		bson.M {
			"$group": bson.M {
				"_id": "$" + key,
				"entries": bson.M {
					"$push": "$$ROOT",
				},
			},
		},
	}
	pipe := collection.Pipe(pipeline)
	results := []map[string]interface{}{}
	pipe.All(&results)
	return results
}

// List lists the models based on the mode given
func List(collection *mgo.Collection, modes ...ListMode)  {
	// var modesList ListModes = modes
	// var timeSpec string
	// if modesList.containMode(REVERSE) {
	// 	timeSpec = "timestamp"
	// } else {
	// 	timeSpec = "-timestamp"
	// }
	// alreadyTrainedQuery := collection.Find(bson.M{}).Sort(timeSpec)
	// if modesList.containMode(TESTED) {
	// 	alreadyTrainedQuery = alreadyTrainedQuery.Select(bson.M{"test": bson.M{"$ne": ""}})
	// }
	// stopAtFive := modesList.containMode(SHORT)
	// results := make([]run.Schema, 0)
	// utils.Check(alreadyTrainedQuery.All(&results))
	results := groupBy(collection, "concrete")
	println(results)
}

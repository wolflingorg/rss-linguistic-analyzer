// Add tasks for analiz
// Uses global variable "DB" and "config"
package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	tm "task-manager"
)

func AddTasksHandler() {
	c := db.C("news")
	var items []Item

	// try to find feeds to update
	err := c.Find(bson.M{"wordmap": bson.M{"$exists": false}}).Limit(config.Handler.Tasks).All(&items)
	if err != nil {
		panic(err)
	}

	// set items to work channel
	for i, value := range items {
		work := tm.WorkRequest{Id: i, Data: value}
		WorkQueue <- work
	}

	// TODO delete this
	fmt.Printf("\n%d tasks added\n", len(items))
}

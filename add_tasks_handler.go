// Add tasks for analiz
// Uses global variable "DB" and "config"
package main

import (
	"gopkg.in/mgo.v2/bson"
	tm "task-manager"
)

func AddTasksHandler() {
	c := db.C("news")
	var items []Item

	// try to find news
	limit := config.Handler.Tasks - tm.GetTasksCount()
	if limit <= 0 {
		LogInfo.Printf("Tasks didnt add. %d active tasks count\n", tm.GetTasksCount())
		return
	}

	err := c.Find(bson.M{
		"$or": []interface{}{
			bson.M{"errors": bson.M{"$lt": 3}},
			bson.M{"errors": bson.M{"$exists": false}},
		},
		"wordmap": bson.M{"$exists": false},
		"_id":     bson.M{"$nin": tm.GetTasksIds()},
	}).Limit(limit).All(&items)
	if err != nil {
		LogError.Fatalf("Couldnt get mongodb result %s\n", err)
	}

	// set items to work channel
	for _, value := range items {
		work := tm.WorkRequest{Id: value.Id, Data: value}
		tm.NewWork(work)
	}

	if len(items) > 0 {
		LogInfo.Printf("%d tasks added\n", len(items))
	}
}

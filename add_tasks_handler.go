// Add tasks for analiz
// Uses global variable "DB" and "config"
package main

import (
	"gopkg.in/mgo.v2/bson"
	tm "task-manager"
)

var tasks_count int = 0

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
		"status": 1,
		"errors": bson.M{"$lt": 3},
		"_id":    bson.M{"$nin": tm.GetTasksIds()},
	}).Limit(limit).All(&items)
	if err != nil {
		LogError.Fatalf("Couldnt get mongodb result %s\n", err)
	}

	// set items to work channel
	for _, value := range items {
		work := tm.WorkRequest{Id: value.Id, Data: value}
		tm.NewWork(work)
	}

	// create new version of dictionary
	if len(items) == 0 && tasks_count > 0 {
		CreateDictionaryVersion(dict_version)
		dict_version += 1
		tasks_count = 0
	} else {
		tasks_count += len(items)
		LogInfo.Printf("%d tasks added\n", len(items))
	}
}

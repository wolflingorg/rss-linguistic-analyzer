package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DictionaryVersion struct {
	Id        int `json:"id" bson:"_id,omitempty"`
	Date      time.Time
	Documents int
}

// Return last dictionary version
func GetLastDictionaryVersion() *DictionaryVersion {
	item := new(DictionaryVersion)
	dv := db.C("dictionary_versions")

	dv.Find(nil).Sort("-_id").Limit(1).One(&item)
	return item
}

// Create new version of dictionary
func CreateDictionaryVersion(version_id int) {
	var err error
	dv := db.C("dictionary_versions")
	n := db.C("news")
	d := db.C("dictionary")
	var items []Item

	// get news with this version id
	err = n.Find(bson.M{
		"status":      2,
		"dictversion": version_id,
	}).All(&items)
	if err != nil {
		LogError.Fatalf("Couldnt get mongodb result %s\n", err)
	}

	// update dictionary
	for _, news := range items {
		for _, value := range news.WordMap {
			var doc interface{}

			change := mgo.Change{
				Update: bson.M{
					"$inc": bson.M{"cnt": 1},
					"$set": bson.M{"ver": version_id, "lang": news.Lang},
				},
				ReturnNew: true,
				Upsert:    true,
			}
			d.Find(bson.M{"word": value.Word, "lang": news.Lang}).Apply(change, &doc)
		}
	}

	item := new(DictionaryVersion)
	item.Id = version_id
	item.Date = time.Now()
	item.Documents, err = n.Find(bson.M{"status": 2}).Count()
	if err != nil {
		LogError.Fatalf("Couldnt create dictionary version %s\n", err)
	}

	dv.Insert(&item)
}

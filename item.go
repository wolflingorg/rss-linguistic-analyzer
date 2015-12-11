// Structure of Item
package main

import (
	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string
	Summary      string
	Content      string
	Lang         string
	WordMap      []MapItem
	WordChecksum []string
}

type MapItem struct {
	Word  string
	Freq  int
	Morph string
}

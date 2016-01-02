// Item morph handler
// Get task and create WordMap and WordChecksums from it
package main

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strings"
	tm "task-manager"
)

func ItemMorphHandler(work tm.WorkRequest, worker_id int) {
	n := db.C("news")

	// check that work.Data equal Item interface
	if item, ok := work.Data.(Item); ok {
		host, _ := config.GetFreeLingHostByLang(item.Lang)
		c, err := connectToFreeLing(host)
		defer c.Close()
		if err != nil {
			LogError.Printf("Worker %d connection ERROR\n", worker_id)
			return
		}

		var content string
		if item.Title == item.Content {
			content = item.Title
		} else {
			content = item.Title + " " + item.Content
		}
		word_map := getWordMap(content, c)

		if len(word_map) == 0 {
			n.Update(bson.M{"_id": item.Id}, bson.M{"$set": bson.M{"errors": item.Errors + 1}})
			LogError.Printf("Worker %d FAILED with News id %s\n", worker_id, item.Id)
			return
		}

		// word checksum
		var word_checksum []string
		for _, value := range word_map {
			hasher := md5.New()
			hasher.Write([]byte(value.Word))
			word_checksum = append(word_checksum, hex.EncodeToString(hasher.Sum(nil)))
		}

		// status
		var status uint
		if len(word_checksum) < 5 {
			status = 0
		} else {
			status = 2
		}

		// update news info
		n.Update(bson.M{"_id": item.Id}, bson.M{"$set": bson.M{
			"wordmap":      word_map,
			"wordchecksum": word_checksum,
			"wordscount":   len(word_checksum),
			"status":       status,
			"dictversion":  dict_version,
		}})
	}
}

// get word map from string
// all words conwerts to their lemma
// we get only verbs, nouns, dates, adjectivs
func getWordMap(msg string, c net.Conn) (result []MapItem) {
	status, err := getMorphResult(msg, c)
	if err != nil {
		LogError.Printf("Couldnt get Freeling results %s", err)
		return nil
	}

	lines := strings.Split(status, "\n")
	for _, value := range lines {
		words := strings.Split(value, " ")

		if words != nil && len(words) > 2 {
			if findInWordMap(result, words[1]) == false {
				if strings.HasPrefix(words[2], "A") || strings.HasPrefix(words[2], "N") || strings.HasPrefix(words[2], "V") || strings.HasPrefix(words[2], "Q") || strings.HasPrefix(words[2], "W") {
					result = append(result, MapItem{
						Word:  words[1],
						Freq:  1,
						Morph: words[2],
					})
				}
			}
		}
	}

	return
}

// util
// try to find word in the word map
func findInWordMap(m []MapItem, s string) bool {
	for i, value := range m {
		if value.Word == s {
			m[i].Freq += 1
			return true
		}
	}

	return false
}

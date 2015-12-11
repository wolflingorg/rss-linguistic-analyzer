// Item morph handler
// Get task and create WordMap and WordChecksums from it
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strings"
	tm "task-manager"
)

func ItemMorphHandler(work tm.WorkRequest, worker_id int) {
	n := db.C("news")

	// check that work.Data equal Item interface
	if item, ok := work.Data.(Item); ok {
		c, err := getConnection(worker_id, item.Lang)
		if err != nil {
			// TODO delete this
			fmt.Printf("\tWorker %d connection ERROR\n", worker_id)
			return
		}

		word_map := getWordMap(item.Title+" "+item.Content, c)

		if len(word_map) == 0 {
			// TODO delete this
			fmt.Printf("\tWorker %d FAILED\n", worker_id)
			return
		}

		// word checksum
		var word_checksum []string
		for _, value := range word_map {
			hasher := md5.New()
			hasher.Write([]byte(value.Word))
			word_checksum = append(word_checksum, hex.EncodeToString(hasher.Sum(nil)))
		}

		// update news info
		n.Update(bson.M{"_id": item.Id}, bson.M{"$set": bson.M{
			"wordmap":      word_map,
			"wordchecksum": word_checksum,
		}})

		// TODO delete this
		fmt.Printf("\tWorker %d OK\n", worker_id)
	}
}

// get word map from string
// all words conwerts to their lemma
// we get only verbs, nouns, dates, adjectivs
func getWordMap(msg string, c net.Conn) (result []MapItem) {
	status, err := getMorphResult(msg, c)
	if err != nil {
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

// we create only one connection to FreeLing server per worker
// this function try to get connection or create it
func getConnection(worker_id int, lang string) (net.Conn, error) {
	var err error
	
	if FreeLingConnMap[worker_id][lang] == nil {
		if FreeLingConnMap[worker_id] == nil {
			FreeLingConnMap[worker_id] = make(map[string]net.Conn)
		}

		FreeLingConnMap[worker_id][lang], err = connectToFreeLing(FreeLingHostsByLang[lang])
		if err != nil {
			FreeLingConnMap[worker_id][lang] = nil
			return nil, err
		}
	}

	return FreeLingConnMap[worker_id][lang], nil
}

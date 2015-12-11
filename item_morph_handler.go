package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strings"
	tm "task-manager"
)

func ItemMorphHandler(work tm.WorkRequest, worker_id int) {
	n := db.C("news")
	var err error

	// check that work.Data equal Item interface
	if item, ok := work.Data.(Item); ok {
		// connect to FreeLing
		if FreeLingConnMap[worker_id][item.Lang] == nil {
			if FreeLingConnMap[worker_id] == nil {
				FreeLingConnMap[worker_id] = make(map[string]net.Conn)
			}

			FreeLingConnMap[worker_id][item.Lang], err = connectToFreeLing(FreeLingHostsByLang[item.Lang])
			if err == nil {
				fmt.Printf("Connect... %d\n", worker_id)
			} else {
				FreeLingConnMap[worker_id][item.Lang] = nil
				return
			}
		}

		word_map := getMorphResult(item.Title+" "+item.Content, FreeLingConnMap[worker_id][item.Lang])

		if len(word_map) == 0 {
			fmt.Printf("\tWorker %d FAILED\n", worker_id)
			return
		}

		// update news info
		n.Update(bson.M{"_id": item.Id}, bson.M{"$set": bson.M{
			"wordmap": word_map,
		}})

		// TODO delete this
		fmt.Printf("\tWorker %d OK\n", worker_id)
	}
}

func getMorphResult(msg string, c net.Conn) (result []MapItem) {
	fmt.Fprintf(c, "%s%c", msg, '\x00')
	status, err := bufio.NewReader(c).ReadString('\x00')
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lines := strings.Split(status, "\n")
	for _, value := range lines {
		words := strings.Split(value, " ")

		if words != nil && len(words) > 2 {
			if findInWordMap(result, words[1]) == false {
				if strings.HasPrefix(words[2], "A") || strings.HasPrefix(words[2], "N") || strings.HasPrefix(words[2], "V") || strings.HasPrefix(words[2], "Q") {
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

func findInWordMap(m []MapItem, s string) bool {
	for i, value := range m {
		if value.Word == s {
			m[i].Freq += 1
			return true
		}
	}

	return false
}

func connectToFreeLing(host string) (c net.Conn, err error) {
	c, err = net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

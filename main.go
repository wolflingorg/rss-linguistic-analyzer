package main

import (
	"flag"
	//"fmt"
	"gopkg.in/mgo.v2"
	"net"
	"strings"
	tm "task-manager"
	"time"
)

var (
	CONFIG_PATH         string                      // PATH to ini file
	config              = new(Config)               // Config struct
	db                  *mgo.Database               // Data Base
	FreeLingHostsByLang map[string]string           // FL hosts by lang
	FreeLingConnMap     map[int]map[string]net.Conn // Map of connections by worker_id and lang
)

func main() {
	// get flags
	flag.StringVar(&CONFIG_PATH, "c", "", "PATH to ini file")
	flag.Parse()

	// config
	LoadConfig(config, CONFIG_PATH)

	// connect to db
	session, err := mgo.Dial(strings.Join(config.Db.Host, ","))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db = session.DB("rss")

	// freeling hosts by lang
	FreeLingHostsByLang = make(map[string]string)
	for _, value := range config.FreeLing.Hosts {
		params := strings.Split(value, "@")
		FreeLingHostsByLang[params[0]] = params[1]
	}
	FreeLingConnMap = make(map[int]map[string]net.Conn)

	// start task manager
	tm.StartDispatcher(config.Handler.Workers, ItemMorphHandler)

	for {
		select {
		case <-time.After(config.Handler.Interval):
			AddTasksHandler()
		}
	}
}

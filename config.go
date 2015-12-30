// Structure of config.ini file
package main

import (
	"errors"
	"gopkg.in/ini.v1"
	"log"
	"strings"
	"time"
)

type Config struct {
	LogPath string
	Db      struct {
		Host []string
	}
	FreeLing struct {
		Hosts []string
	}
	Handler struct {
		Workers  int
		Tasks    int
		Interval time.Duration
	}
}

// Load and Map config from file
func LoadConfig(config *Config, CONFIG_PATH string) {
	err := ini.MapTo(config, CONFIG_PATH)
	if err != nil {
		log.Fatalf("Couldnt parse config file: %s\n", err)
	}
}

// Get FreeLing host by lang
func (this *Config) GetFreeLingHostByLang(lang string) (string, error) {
	for _, value := range this.FreeLing.Hosts {
		params := strings.Split(value, "@")
		if params[0] == lang {
			return params[1], nil
		}
	}

	return "", errors.New("Couldnt find Host by this lang")
}

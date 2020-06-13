package common

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var config = make(map[string]string)

func init() {
	loadConfig()
}

func loadConfig() {
	configFile := "app.conf"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}
	f, _ := os.Open(configFile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keyPair := strings.Split(line, "=")
		if len(keyPair) < 2 {
			continue
		}
		key := strings.TrimSpace(keyPair[0])
		val := strings.TrimSpace(keyPair[1])
		config[key] = val
		log.Printf("load config: %v=%v \n", key, val)
	}
}

func GetConfig(key, def string) string {
	val, ok := config[key]
	if ok {
		return val
	}else {
		return def
	}
}

func GetConfigBool(key string) bool {
	val, ok := config[key]
	if ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatal(key, "is not bool type!")
		}
		return v
	}
	return false
}



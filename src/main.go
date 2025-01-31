package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

type CalEntry struct {
	Title string
	Urls  []string
}
type CalMap map[string]CalEntry

var CfgPath, Cronjob, ProdID, ContentDir, Addr, Port string

func readEnv() error {
	var there bool
	keys := []string{"CFG_PATH", "CRON", "PROD_ID", "CONTENT_DIR", "ADDR", "PORT"}
	vars := []*string{&CfgPath, &Cronjob, &ProdID, &ContentDir, &Addr, &Port}

	for i, key := range keys {
		ptr := vars[i]

		*ptr, there = os.LookupEnv(key)
		if !there {
			return fmt.Errorf("missing %s", key)
		}
		log.Println(key, "=", *ptr)
	}

	return nil
}

func main() {
	log.Println("CalUnite started, reading environment variables")

	err := readEnv()
	if err != nil {
		panic(err)
	}

	calmap, err := parseYml(CfgPath)
	if err != nil {
		panic(err)
	}

	// to make relative paths work
	err = os.MkdirAll(ContentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.Chdir(ContentDir)
	if err != nil {
		panic(err)
	}

	// start first merge immediately
	merger := unite(calmap)
	merger()
	log.Println("Initial merge finished")

	c := cron.New()
	c.AddFunc(Cronjob, merger)
	c.Start()
	log.Println("Started merge cronjob")

	serve()
}

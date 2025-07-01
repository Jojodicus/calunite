package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
)

type CalEntry struct {
	Title string
	Urls  []string
}
type CalDatum struct {
	File  string
	Entry CalEntry
}
type CalData []CalDatum

var CfgPath, Cronjob, ProdID, ContentDir, FileNavigation, Addr, Port string
var cronRunner *cron.Cron

func readEnv() error {
	var there bool
	keys := []string{"CFG_PATH", "CRON", "PROD_ID", "CONTENT_DIR", "FILE_NAVIGATION", "ADDR", "PORT"}
	vars := []*string{&CfgPath, &Cronjob, &ProdID, &ContentDir, &FileNavigation, &Addr, &Port}

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

func createDirsAndCd() {
	// to make relative paths work
	err := os.MkdirAll(ContentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.Chdir(ContentDir)
	if err != nil {
		panic(err)
	}
}

func mergeAndSchedule(c *cron.Cron) {
	// parse immediately, as little downtime as possible
	calmap, err := parseYml(CfgPath)
	if err != nil {
		panic(err)
	}

	entries := c.Entries()
	if len(entries) != 0 {
		if len(entries) != 1 {
			panic(fmt.Errorf("unexpected number of entries: %v", entries))
		}

		c.Remove(entries[0].ID)
		<-c.Stop().Done() // avoid potential race condition
		log.Println("Stopped previous cronjob")

		// clean state
		err := os.RemoveAll(ContentDir)
		if err != nil {
			panic(err)
		}
		createDirsAndCd()
	}
	// --- only section during reload with downtime ---

	// start first merge immediately
	merger := unite(calmap)
	merger()
	log.Println("Initial merge finished")

	c.AddFunc(Cronjob, merger)
	c.Start()
	log.Println("Started merge cronjob")
}

func main() {
	log.Println("CalUnite started, reading environment variables")

	err := readEnv()
	if err != nil {
		panic(err)
	}

	// watch config for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(CfgPath)
	if err != nil {
		log.Fatal(err)
	}
	go watchYml(watcher)

	// also sets PWD for serve()
	createDirsAndCd()

	cronRunner = cron.New()
	mergeAndSchedule(cronRunner)

	serve()
}

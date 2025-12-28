package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
)

type CalEntry struct {
	Title *string
	Urls  []string
}
type CalDatum struct {
	File  string
	Entry CalEntry
}
type CalData []CalDatum

var CfgPath, Cronjob, ProdID, ContentDir, FileNavigation, DotPrivate, Addr, Port string
var cronRunner *cron.Cron

func readEnv() error {
	var there bool
	keys := []string{"CFG_PATH", "CRON", "PROD_ID", "CONTENT_DIR", "FILE_NAVIGATION", "DOT_PRIVATE", "ADDR", "PORT"}
	vars := []*string{&CfgPath, &Cronjob, &ProdID, &ContentDir, &FileNavigation, &DotPrivate, &Addr, &Port}

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

func createDirsAndCd() error {
	// to make relative paths work
	err := os.MkdirAll(ContentDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create new directories: %v", err)
	}
	err = os.Chdir(ContentDir)
	if err != nil {
		return fmt.Errorf("could not change working directory: %v", err)
	}
	return nil
}

func mergeAndSchedule(cronjobs *cron.Cron) error {
	// parse immediately, as little downtime as possible
	calmap, err := parseYml(CfgPath)
	if err != nil {
		return fmt.Errorf("could not parse yml: %v", err)
	}

	entries := cronjobs.Entries()
	if len(entries) != 0 {
		for i := range len(entries) {
			cronjobs.Remove(entries[i].ID)
		}

		<-cronjobs.Stop().Done() // avoid potential race condition
		log.Println("Stopped previous cronjob")

		// clean state
		err := os.RemoveAll(ContentDir)
		if err != nil {
			return fmt.Errorf("could not clean old directory: %v", err)
		}
		err = createDirsAndCd()
		if err != nil {
			return err
		}
	}
	// --- only section during reload with downtime ---

	// start first merge immediately
	merger := unite(calmap)
	merger()
	log.Println("Initial merge finished")

	cronjobs.AddFunc(Cronjob, merger)
	cronjobs.Start()
	log.Println("Started merge cronjob")
	return nil
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
	err = watcher.Add(path.Dir(CfgPath))
	if err != nil {
		log.Fatal(err)
	}
	go watchYml(watcher)

	// also sets PWD for serve()
	createDirsAndCd()

	cronRunner = cron.New()
	err = mergeAndSchedule(cronRunner)
	if err != nil {
		log.Println("ERROR", err)
		log.Println("No files will be served, check config or permissions!")
	}

	serve()
}

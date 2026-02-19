package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
)

var logLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type CalEntry struct {
	Title       *string  `yaml:"title"`
	EventFormat *string  `yaml:"event_format"`
	Urls        []string `yaml:"urls"`
}
type CalDatum struct {
	File  string
	Entry CalEntry
}
type CalData []CalDatum

var CfgPath, Cronjob, ProdID, ContentDir, FileNavigation, DotPrivate, Addr, Port, LogLevel string
var cronRunner *cron.Cron

func readEnv() error {
	slog.Info("Reading configured environment variables")

	var there bool
	keys := []string{"CFG_PATH", "CRON", "PROD_ID", "CONTENT_DIR", "FILE_NAVIGATION", "DOT_PRIVATE", "ADDR", "PORT", "LOG_LEVEL"}
	vars := []*string{&CfgPath, &Cronjob, &ProdID, &ContentDir, &FileNavigation, &DotPrivate, &Addr, &Port, &LogLevel}

	for i, key := range keys {
		ptr := vars[i]

		*ptr, there = os.LookupEnv(key)
		if !there {
			return fmt.Errorf("missing %s", key)
		}
		slog.Debug(fmt.Sprint(key, "=", *ptr))
	}

	return nil
}

func createDirsAndCd() error {
	slog.Debug("Creating directories and changing CWD")

	// to make relative paths work
	err := os.MkdirAll(ContentDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not create new directories: %v", err)
	}
	err = os.Chdir(ContentDir)
	if err != nil {
		return fmt.Errorf("Could not change working directory: %v", err)
	}
	return nil
}

func mergeAndSchedule(cronjobs *cron.Cron) error {
	// parse immediately, as little downtime as possible
	calmap, err := parseYml(CfgPath)
	if err != nil {
		return fmt.Errorf("Could not parse yml: %v", err)
	}

	entries := cronjobs.Entries()
	if len(entries) != 0 {
		for i := range len(entries) {
			entry := entries[i]
			slog.Debug(fmt.Sprintf("Removing cron entry %v", entry))
			cronjobs.Remove(entry.ID)
		}

		slog.Info("Stopping previous cronjob")
		<-cronjobs.Stop().Done() // avoid potential race condition

		// clean state
		slog.Info("Removing old calendars")
		err := os.RemoveAll(ContentDir)
		if err != nil {
			return fmt.Errorf("Could not clean old directory: %v", err)
		}
		err = createDirsAndCd()
		if err != nil {
			return err
		}
	}
	// --- only section during reload with downtime ---

	// start first merge immediately
	slog.Info("Starting initial merge, this can take a while...")
	merger := unite(calmap)
	merger()
	slog.Info("Initial merge finished, starting merge cronjob")

	_, err = cronjobs.AddFunc(Cronjob, merger)
	if err != nil {
		return fmt.Errorf("Unable to register cronjob: %v", err)
	}
	cronjobs.Start()
	slog.Debug("Started merge cronjob")
	return nil
}

func main() {
	// initialize logger - debug as initial level
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelDebug)
	logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))

	slog.Info("Welcome to CalUnite!")

	err := readEnv()
	if err != nil {
		Fatal("Error while reading environment variables", err)
	}

	newLevel := logLevels[strings.ToLower(LogLevel)]
	slog.Debug(fmt.Sprintf("Setting log-level to %v", newLevel))
	logLevel.Set(newLevel)

	// watch config for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Fatal("Error while creating file watcher", err)
	}
	defer watcher.Close()
	err = watcher.Add(path.Dir(CfgPath))
	if err != nil {
		Fatal("Error while configuring file watcher", err)
	}
	go watchYml(watcher)

	// also sets PWD for serve()
	createDirsAndCd()

	cronRunner = cron.New()
	err = mergeAndSchedule(cronRunner)
	if err != nil {
		slog.Warn("No files will be served, check config or permissions!", ErrAttr(err))
	}

	serve()
}

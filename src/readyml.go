package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type CalMap map[string]CalEntry

func (calmap CalMap) String() string {
	var sb strings.Builder
	sb.WriteRune('|')
	for s := range calmap {
		sb.WriteString(s)
		sb.WriteRune('|')
	}
	return sb.String()
}

func (caldata CalData) contains(file string) bool {
	for _, e := range caldata {
		if e.File == file {
			return true
		}
	}
	return false
}

func topoSort(calmap CalMap) (CalData, error) {
	// sort out cyclic recursive definitions

	// ordered list that satisfies all recursive dependencies
	caldata := make(CalData, 0, len(calmap))

	// continuously add "fine" entries until done
	worklist := calmap
	for len(worklist) > 0 {
		changed := false
		updatedWorklist := worklist

		// find candidate
		for k, e := range worklist {
			ok := true
			// recursive definition?
			for _, calendar := range e.Urls {
				_, there := calmap[calendar]
				if there && !caldata.contains(calendar) {
					// recursive definition, but not yet sorted
					ok = false
					break
				}
			}

			if ok {
				// definition fine, add to output
				caldata = append(caldata, CalDatum{k, e})
				changed = true
				delete(updatedWorklist, k)
			}
		}

		// cycle detected
		if !changed {
			return caldata, fmt.Errorf("config contains cycles: %s", worklist)
		}

		worklist = updatedWorklist
	}

	return caldata, nil
}

func parseYml(filename string) (CalData, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	log.Println(filename, "read successfully")

	var calmap CalMap

	err = yaml.Unmarshal(file, &calmap)
	if err != nil {
		return nil, err
	}
	log.Println("YAML parsing done")

	return topoSort(calmap)
}

func watchYml(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Name == CfgPath {
				log.Println("Modified", event.Name, "reloading...")
				err := mergeAndSchedule(cronRunner)
				if err != nil {
					// kinda copy paste from main()...
					log.Println("ERROR", err)
					log.Println("No files will be served, check config or permissions!")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Fatal(err)
		}
	}
}

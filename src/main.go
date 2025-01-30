package main

import (
	"os"

	"github.com/robfig/cron/v3"
)

type CalEntry struct {
	Title string
	Urls  []string
}
type CalMap map[string]CalEntry

const CfgPath = "./test.yml"
const ContentDir = "./wwwdata"
const Addr = "0.0.0.0"
const Port = 8080
const Cronjob = "@every 5s" // format: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format

func main() {
	calmap, err := parseYml(CfgPath)
	if err != nil {
		panic(err)
	}

	// to make relative paths work
	os.Chdir(ContentDir)

	c := cron.New()
	c.AddFunc(Cronjob, unite(calmap))
	c.Start()

	serve()
}

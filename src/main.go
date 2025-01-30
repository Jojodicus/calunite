package main

import "github.com/robfig/cron/v3"

type Calmap map[string]struct {
	Title string
	Urls  []string
}

const CfgName = "./config.yml"
const ContentDir = "./wwwdata"
const Addr = ""
const Port = 8080
const Cronjob = "@every 5s"

func main() {
	calmap, err := parseYml(CfgName)
	if err != nil {
		panic(err)
	}

	c := cron.New()
	c.AddFunc(Cronjob, unite(calmap))
	c.Start()

	serve()
}

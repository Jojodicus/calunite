package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func parseYml(filename string) (CalMap, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var calmap CalMap

	err = yaml.Unmarshal(file, &calmap)
	if err != nil {
		return nil, err
	}

	return calmap, nil
}

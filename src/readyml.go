package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func parseYml(filename string) (CalMap, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println(filename, "read successfully")

	var calmap CalMap

	err = yaml.Unmarshal(file, &calmap)
	if err != nil {
		return nil, err
	}
	fmt.Println("YAML parsing done")

	return calmap, nil
}

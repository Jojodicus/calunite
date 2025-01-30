package main

import (
	"fmt"
)

func unite(calmap Calmap) func() {
	// closures are nice!
	return func() {
		for calendar, entry := range calmap {
			fmt.Println(calendar, entry.Title)
			for _, url := range entry.Urls {
				fmt.Println(url)
			}
			fmt.Println()
		}
	}
}

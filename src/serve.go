package main

import (
	"fmt"
	"log"
	"net/http"
)

func serve() {
	fmt.Println("Running webserver on", fmt.Sprintf("%s:%s", Addr, Port))
	http.Handle("/", http.FileServer(http.Dir("."))) // we are already in the content directory
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", Addr, Port), nil))
}

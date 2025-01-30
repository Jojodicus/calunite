package main

import (
	"fmt"
	"log"
	"net/http"
)

func serve() {
	http.Handle("/", http.FileServer(http.Dir(ContentDir)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", Addr, Port), nil))
}

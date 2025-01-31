package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type justFilesFilesystem struct {
	Fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)

	if err != nil {
		return nil, err
	}

	stat, _ := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func serve() {
	var fileSystem http.FileSystem = http.Dir(".")
	nav, err := strconv.ParseBool(FileNavigation)
	if err == nil && !nav {
		fileSystem = justFilesFilesystem{fileSystem}
	}

	log.Println("Running webserver on", fmt.Sprintf("%s:%s", Addr, Port))
	http.Handle("/", http.FileServer(fileSystem)) // we are already in the content directory
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", Addr, Port), nil))
}

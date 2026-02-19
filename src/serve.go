package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type justFilesFilesystem struct {
	Fs         http.FileSystem
	DotPrivate bool
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	// don't log file name, scared because of log4shell

	if fs.DotPrivate && len(name) > 2 && name[1] == '.' {
		// "/.calendar.ics"
		slog.Debug("Attempt to open private file")
		return nil, os.ErrNotExist
	}

	f, err := fs.Fs.Open(name)

	if err != nil {
		slog.Debug("Failed to open file", ErrAttr(err))
		return nil, err
	}

	stat, _ := f.Stat()
	if stat.IsDir() {
		slog.Debug("Attempt to list directory contents")
		return nil, os.ErrNotExist
	}

	return f, nil
}

func serve() {
	var fileSystem http.FileSystem = http.Dir(".")
	nav, err := strconv.ParseBool(FileNavigation)
	if err == nil && !nav {
		slog.Debug("File navigation is turned off, creating custom file system wrapper")
		dotPriv, err := strconv.ParseBool(DotPrivate)
		if err == nil {
			fileSystem = justFilesFilesystem{fileSystem, dotPriv}
		}
	}

	slog.Info(fmt.Sprintf("Running webserver on %s:%s", Addr, Port))
	http.Handle("/", http.FileServer(fileSystem)) // we are already in the content directory
	Fatal("Failed to start webserver", http.ListenAndServe(fmt.Sprintf("%s:%s", Addr, Port), nil))
}

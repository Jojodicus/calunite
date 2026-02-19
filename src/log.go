package main

import (
	"log/slog"
	"os"
)

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func Fatal(msg string, err error) {
	slog.Error(msg, ErrAttr(err))
	os.Exit(1)
}

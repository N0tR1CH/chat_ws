package main

import (
	"io"
	"log/slog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

func NewSlogLogger(w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, nil)
	logger := slog.New(handler)
	return logger
}

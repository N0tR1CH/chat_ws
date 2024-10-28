package main

import (
	"html/template"
	"net/http"
)

type Config struct {
	Host string
	Port string
}

func NewServer(
	logger Logger,
	templates *template.Template,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger, templates)
	var handler http.Handler = mux
	return handler
}

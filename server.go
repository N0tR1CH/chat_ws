package main

import "net/http"

type Config struct {
	Host string
	Port string
}

func NewServer(logger Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger)
	var handler http.Handler = mux
	return handler
}

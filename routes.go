package main

import "net/http"

func addRoutes(mux *http.ServeMux, logger Logger) {
	logger.Info("Message", "test", "test")
	mux.Handle("/", http.NotFoundHandler())
}

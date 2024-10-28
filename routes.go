package main

import (
	"html/template"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger Logger,
	templates *template.Template,
) {
	mux.Handle("/", handleChatGet(logger, templates))
	mux.Handle("/room", handleRoomWs(logger))
}

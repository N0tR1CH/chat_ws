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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("/", "host", r.Host)
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			logger.Error("rendering template error", "err", err.Error())
		}
	})
}

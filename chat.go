package main

import (
	"html/template"
	"net/http"
)

func handleChatGet(logger Logger, templates *template.Template) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Bad method", http.StatusBadRequest)
			}
			logger.Info("handleChatGet", "msg", "Getting chat")
			if err := templates.ExecuteTemplate(w, "chat.html", nil); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				logger.Error("rendering template error", "err", err.Error())
			}
		},
	)
}

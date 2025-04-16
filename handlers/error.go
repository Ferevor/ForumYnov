package handlers

import (
	"forumynov/database"
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int) {
	msg := http.StatusText(statusCode)
	if msg == "" {
		statusCode = http.StatusInternalServerError
		msg = "Unknown error"
	}

	w.WriteHeader(statusCode)

	tmpl, err := template.ParseFiles("./templates/error.html")
	if err != nil {
		log.Println("[handlers/error.go] [ErrorHandler] Error ParseFiles >>>", err)
		http.Error(w, msg, statusCode)
		return
	}

	data := database.ErrorPageData{
		Code:    statusCode,
		Message: msg,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("[handlers/error.go] [ErrorHandler] Error Execute >>>", err)
		http.Error(w, msg, statusCode)
		return
	}
}

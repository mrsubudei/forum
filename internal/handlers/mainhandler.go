package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ViewHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		// Errors(writer, "Page is Not Found", http.StatusNotFound)
		log.Fatal()
		return
	}
	if request.Method != http.MethodGet {
		log.Fatal()
		// Errors(writer, "Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
		// Errors(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = html.Execute(writer, nil)
	if err != nil {
		log.Fatal(err)
		// Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
}

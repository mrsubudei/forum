package main

import (
	"fmt"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ViewHandler)
	// mux.HandleFunc("/ascii-art", app.CreateHandler)
	mux.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	// mux.Handle("/templates/img/", http.StripPrefix("/templates/img/", http.FileServer(http.Dir("templates/img"))))
	fmt.Println("Starting server at post: 8087\n" +
		"http://localhost:8087/")
	err := http.ListenAndServe("localhost:8087", mux)
	log.Fatal(err)
}

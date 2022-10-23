package main

import (
	"errors"
	"forum/internal/database"
	"log"
	"os"
)

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handlers.ViewHandler)
	// mux.HandleFunc("/register", handlers.RegisterHandler)
	// mux.HandleFunc("/registration", handlers.RegistrationHandler)
	exist, err := Exists("forum.db")
	if err != nil {
		log.Fatal(err)
	}

	if !exist {
		database.CreateDB()
	}
	database.WriteToUsers()
	// mux.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	// fmt.Println("Starting server at post: 8087\n" +
	// 	"http://localhost:8087/")
	// err := http.ListenAndServe(":8087", mux)
	// log.Fatal(err)
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

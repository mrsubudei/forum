package http

import (
	"html/template"
	"net/http"
)

type Asc struct {
	Str      string
	Orig     string
	StatCode int
}

type Userr struct {
	Name  string
	Email string
}

var Ascii Asc

func ViewHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
	if request.Method != http.MethodGet {
		Errors(writer, "Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		Errors(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = html.Execute(writer, nil)
	if err != nil {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
}

func RegisterHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/register" {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
	if request.Method != http.MethodGet {
		Errors(writer, "Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	html, err := template.ParseFiles("templates/register.html")
	if err != nil {
		Errors(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// name := request.FormValue("name")
	// email := request.FormValue("email")

	err = html.Execute(writer, nil)
	if err != nil {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
}

func RegistrationHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/registration" {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}

	if request.Method != http.MethodPut {
		Errors(writer, "Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}

	html, err := template.ParseFiles("templates/register.html")
	if err != nil {
		Errors(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = html.Execute(writer, nil)
	if err != nil {
		Errors(writer, "Page is Not Found", http.StatusNotFound)
		return
	}
}

func Errors(writer http.ResponseWriter, er string, code int) {
	writer.WriteHeader(code)
	Ascii.Str = er
	Ascii.StatCode = code
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(writer, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = html.Execute(writer, Ascii)
	if err != nil {
		http.Error(writer, "404: Not Found", 404)
		return
	}
}

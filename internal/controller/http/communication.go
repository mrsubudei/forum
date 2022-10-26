package http

import (
	"forum/internal/usecase"
)

type CommunicationRoutes struct {
	c usecase.Communication
}

// func NewCommunicationRoute(w http.ResponseWriter, r *http.Request, c usecase.Communication) {
// 	route := &CommunicationRoutes{c}
// 	if r.URL.Path != "/" {
// 		Errors(w, "Page is Not Found", http.StatusNotFound)
// 		return
// 	}
// 	if r.Method != http.MethodGet {
// 		Errors(w, "Method is not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	html, err := template.ParseFiles("templates/index.html")
// 	if err != nil {
// 		Errors(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// 	err = html.Execute(w, nil)
// 	if err != nil {
// 		Errors(w, "Page is Not Found", http.StatusNotFound)
// 		return
// 	}
// }

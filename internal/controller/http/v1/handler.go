package v1

import (
	"forum/internal/entity"
	"forum/internal/usecase"
	"net/http"
	"text/template"
)

type Handler struct {
	usecases *usecase.UseCases
}

type Content struct {
	Authorized   bool
	Unauthorized bool
	Posts        []entity.Post
}

func NewHandler(services *usecase.UseCases) *Handler {
	return &Handler{
		usecases: services,
	}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)

	if r.URL.Path != "/" {
		http.Error(w, "404: Page is Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	posts, err := h.usecases.Posts.GetAllPosts()
	if err != nil {
		http.Error(w, "404: Not Found", 404)
		return
	}
	content := Content{
		Authorized:   authorized,
		Unauthorized: !authorized,
		Posts:        posts,
	}
	err = html.Execute(w, content)
	if err != nil {
		http.Error(w, "404: Not Found", 404)
		return
	}
}

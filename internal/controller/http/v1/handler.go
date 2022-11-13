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
	Admin        bool
	Posts        []entity.Post
	Users        []entity.User
	Comments     []entity.Comment
	ErrorMsg     ErrMessage
}

func NewHandler(services *usecase.UseCases) *Handler {
	return &Handler{
		usecases: services,
	}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	foundUser := h.getExistedSession(w, r)
	content := Content{}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized
	errors := ErrMessage{}

	if r.URL.Path != "/" {
		errors.Code = http.StatusNotFound
		errors.Message = pageNotFound
		h.Errors(w, errors)
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
	content.Posts = posts

	err = html.Execute(w, content)
	if err != nil {
		http.Error(w, "404: Not Found", 404)
		return
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {

	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)
}

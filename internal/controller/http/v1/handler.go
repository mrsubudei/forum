package v1

import (
	"forum/internal/entity"
	"forum/internal/usecase"
	"net/http"
	"strings"
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

	if r.URL.Path != "/" {
		errors.Code = http.StatusNotFound
		errors.Message = errPageNotFound
		h.Errors(w, errors)
		return
	}
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
	posts, err := h.usecases.Posts.GetAllPosts()
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}
	content.Posts = posts

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SearchPageHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	foundUser := h.getExistedSession(w, r)
	content := Content{}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/search.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}
	r.ParseForm()
	if len(r.Form["search"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}

	searchRequest := strings.ToLower(r.Form["search"][0])
	posts, err := h.usecases.Posts.GetAllPosts()
	filtered := h.filterPosts(posts, searchRequest)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}

	authorized := h.checkSession(w, r)
	foundUser := h.getExistedSession(w, r)
	content := Content{
		Posts: filtered,
	}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)
}

func (h *Handler) filterPosts(posts []entity.Post, request string) []entity.Post {
	var filtered []entity.Post
	if len(posts) == 0 {
		return filtered
	}
	for _, post := range posts {
		found := false
		if strings.Contains(strings.ToLower(post.Title), request) {
			filtered = append(filtered, post)
			continue
		} else if strings.Contains(strings.ToLower(post.Content), request) {
			filtered = append(filtered, post)
			continue
		} else if strings.Contains(strings.ToLower(post.User.Name), request) {
			filtered = append(filtered, post)
			continue
		}
		for _, val := range post.Categories {
			if strings.Contains(strings.ToLower(val), request) {
				filtered = append(filtered, post)
				found = true
				break
			}
		}
		if !found {
			for _, val := range post.Comments {
				if strings.Contains(strings.ToLower(val.Content), request) {
					filtered = append(filtered, post)
					break
				}
			}
		}
	}
	return filtered
}

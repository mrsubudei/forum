package v1

import (
	"fmt"
	"forum/internal/entity"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	if r.URL.Path != "/" {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	key := Key("content")

	content, ok := r.Context().Value(key).(Content)
	if !ok {
		log.Printf("v1 - IndexHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - IndexHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	posts, err := h.usecases.Posts.GetAllPosts()
	if err != nil {
		log.Println(fmt.Errorf("v1 - IndexHandler - GetAllPosts: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	content.Posts = posts

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SearchPageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - SearchPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/search.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchPageHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchPageHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}
	r.ParseForm()
	if len(r.Form["search"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	searchRequest := strings.ToLower(r.Form["search"][0])
	posts, err := h.usecases.Posts.GetAllPosts()
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchHandler - GetAllPosts: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	filtered := h.filterPosts(posts, searchRequest)
	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - SearchHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	content.Posts = filtered

	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreateCategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CreateCategoryPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	if !content.Admin {
		errors.Code = http.StatusForbidden
		errors.Message = ErrLowAccessLevel
		h.Errors(w, errors)
		return
	}

	html, err := template.ParseFiles("templates/create_category.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCategoryPageHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCategoryPageHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	r.ParseForm()
	if len(r.Form["category"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}

	data := r.Form["category"][0]
	categories := strings.Split(data, "\r\n")

	err := h.usecases.Posts.CreateCategories(categories)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCategoryHandler - CreateCategories: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) SearchByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	category := path[len(path)-1]
	if r.URL.Path != "/categories/"+category {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - SearchByCategoryHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	posts, err := h.usecases.Posts.GetAllByCategory(category)
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchByCategoryHandler - GetAllByCategory: %w", err))
		if strings.Contains(err.Error(), entity.ErrPostNotFound.Error()) {
			errors.Code = http.StatusBadRequest
			errors.Message = ErrBadRequest
			h.Errors(w, errors)
			return
		}
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
	content.Posts = posts
	html, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchByCategoryHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - SearchByCategoryHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
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

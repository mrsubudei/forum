package v1

import (
	"forum/internal/entity"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (h *Handler) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/posts/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = errPageNotFound
		h.Errors(w, errors)
		return
	}

	authorized := h.checkSession(w, r)
	foundUser := h.getExistedSession(w, r)
	if authorized {
		err := h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	}
	content := ContentSingle{}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	html, err := template.ParseFiles("templates/post.html")
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	post, err := h.usecases.Posts.GetById(int64(id))
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}
	post.ContentWeb = strings.Split(post.Content, "\\n")

	content.Post = post

	err = html.Execute(w, content)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	authorized := h.checkSession(w, r)
	if !authorized {
		errors.Code = http.StatusForbidden
		errors.Message = errStatusNotAuthorized
		h.Errors(w, errors)
		return
	}
	foundUser := h.getExistedSession(w, r)
	if authorized {
		err := h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	}
	content := ContentSingle{}
	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	categories, err := h.usecases.Posts.GetAllCategories()
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
	content.Post.Categories = categories
	html, err := template.ParseFiles("templates/create_post.html")
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

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	authorized := h.checkSession(w, r)
	if !authorized {
		errors.Code = http.StatusForbidden
		errors.Message = errStatusNotAuthorized
		h.Errors(w, errors)
		return
	}

	foundUser := h.getExistedSession(w, r)
	if authorized {
		err := h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	}

	r.ParseForm()
	if len(r.Form["title"]) == 0 || len(r.Form["content"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}

	content := ContentSingle{}
	valid := true

	postTitle := r.Form["title"][0]
	postContent := r.Form["content"][0]

	categories := r.Form["categories"]

	if len(categories) == 0 {
		content.ErrorMsg.Message = postCategoryRequired
		valid = false
	}

	newPost := entity.Post{}
	newPost.Title = postTitle
	newPost.Content = strings.ReplaceAll(postContent, "\r\n", "\\n")
	newPost.Categories = categories
	newPost.User = foundUser

	if !valid {
		html, err := template.ParseFiles("templates/create_post.html")
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
		categories, err := h.usecases.Posts.GetAllCategories()
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
		content.Post.Categories = categories
		err = html.Execute(w, content)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	} else {
		err := h.usecases.Posts.CreatePost(newPost)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *Handler) PostPutLikeHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	if !authorized {
		errors.Code = http.StatusForbidden
		errors.Message = errStatusNotAuthorized
		h.Errors(w, errors)
		return
	}
	foundUser := h.getExistedSession(w, r)
	if authorized {
		err := h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/put_post_like/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	content := ContentSingle{}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = foundUser.Id

	err = h.usecases.Posts.MakeReaction(post, commandPutLike)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) PostPutDislikeHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	if !authorized {
		errors.Code = http.StatusForbidden
		errors.Message = errStatusNotAuthorized
		h.Errors(w, errors)
		return
	}
	foundUser := h.getExistedSession(w, r)
	if authorized {
		err := h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			errors.Code = http.StatusInternalServerError
			errors.Message = errInternalServer
			h.Errors(w, errors)
			return
		}
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/put_post_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	content := ContentSingle{}

	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

	post := entity.Post{
		Id: int64(id),
	}
	post.User.Id = foundUser.Id

	err = h.usecases.Posts.MakeReaction(post, commandPutDislike)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

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
	if r.URL.Path != "/posts/"+path[len(path)-1] || err != nil || id == 0 || id < 0 {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}

	authorized := h.checkSession(w, r)
	foundUser := h.getExistedSession(w, r)
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
	content := ContentSingle{}
	if foundUser.Id == 1 {
		content.Admin = true
	}
	content.Authorized = authorized
	content.Unauthorized = !authorized

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

func (h *Handler) PostPutLikeHandler(w http.ResponseWriter, r *http.Request) {
	authorized := h.checkSession(w, r)
	if !authorized {
		errors.Code = http.StatusForbidden
		errors.Message = errStatusNotAuthorized
		h.Errors(w, errors)
		return
	}
	foundUser := h.getExistedSession(w, r)

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/put_post_like/"+path[len(path)-1] || err != nil || id == 0 || id < 0 {
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

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/put_post_dislike/"+path[len(path)-1] || err != nil || id == 0 || id < 0 {
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

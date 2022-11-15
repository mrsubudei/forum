package v1

import (
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

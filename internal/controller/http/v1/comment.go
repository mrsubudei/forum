package v1

import (
	"forum/internal/entity"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (h *Handler) CreateCommentPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = errMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/create_comment_page/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = errPageNotFound
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
	content.User.Id = foundUser.Id
	content.Post.Id = int64(id)

	html, err := template.ParseFiles("templates/create_comment.html")
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

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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
	if len(r.Form["content"]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/create_comment/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = errPageNotFound
		h.Errors(w, errors)
		return
	}

	commentContent := r.Form["content"][0]
	newComment := entity.Comment{}
	newComment.Content = strings.ReplaceAll(commentContent, "\r\n", "\\n")
	newComment.User = foundUser
	newComment.PostId = int64(id)

	err = h.usecases.Comments.WriteComment(newComment)
	if err != nil {
		errors.Code = http.StatusInternalServerError
		errors.Message = errInternalServer
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (h *Handler) CommentPutLikeHandler(w http.ResponseWriter, r *http.Request) {
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
	if r.URL.Path != "/put_comment_like/"+path[len(path)-1] || err != nil || id <= 0 {
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

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = foundUser.Id

	err = h.usecases.Comments.MakeReaction(comment, commandPutLike)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (h *Handler) CommentPutDislikeHandler(w http.ResponseWriter, r *http.Request) {
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
	if r.URL.Path != "/put_comment_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
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

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = foundUser.Id

	err = h.usecases.Comments.MakeReaction(comment, commandPutDislike)
	if err != nil {
		errors.Code = http.StatusBadRequest
		errors.Message = errBadRequest
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

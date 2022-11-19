package v1

import (
	"forum/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) CommentPutLikeHandler(w http.ResponseWriter, r *http.Request) {
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
	if r.URL.Path != "/put_comment_like/"+path[len(path)-1] || err != nil || id == 0 || id < 0 {
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

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if r.URL.Path != "/put_comment_dislike/"+path[len(path)-1] || err != nil || id == 0 || id < 0 {
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

package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"forum/internal/entity"
)

func (h *Handler) CreateCommentPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCommentPageHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/create_comment_page/"+path[len(path)-1] || err != nil || id <= 0 {
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CreateCommentPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	content.Post.Id = int64(id)

	html, err := template.ParseFiles("templates/create_comment.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCommentPageHandler - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCommentPageHandler - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
		h.Errors(w, errors)
		return
	}

	r.ParseForm()
	if len(r.Form["content"][0]) == 0 {
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCommentHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/create_comment/"+path[len(path)-1] || err != nil || id <= 0 {
		log.Println(fmt.Errorf("v1 - CreateCommentHandler - URL.Path: %w", err))
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CreateCommentHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	commentContent := r.Form["content"][0]
	newComment := entity.Comment{}
	newComment.Content = strings.ReplaceAll(commentContent, "\r\n", "\\n")
	newComment.User = content.User
	newComment.PostId = int64(id)

	err = h.usecases.Comments.WriteComment(newComment)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CreateCommentHandler - WriteComment: %w", err))
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (h *Handler) CommentPutLikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - CommentPutLikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_comment_like/"+path[len(path)-1] || err != nil || id <= 0 {
		log.Println(fmt.Errorf("v1 - CommentPutLikeHandler - URL.Path: %w", err))
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CommentPutLikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = content.User.Id

	err = h.usecases.Comments.MakeReaction(comment, CommandPutLike)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CommentPutLikeHandler - MakeReaction: %w", err))
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer")+"#"+strconv.Itoa(int(comment.Id)), http.StatusFound)
}

func (h *Handler) CommentPutDislikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		log.Println(fmt.Errorf("v1 - CommentPutDislikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_comment_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
		log.Println(fmt.Errorf("v1 - CommentPutDislikeHandler - URL.Path: %w", err))
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		log.Printf("v1 - CommentPutDislikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content)
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return
	}

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = content.User.Id

	err = h.usecases.Comments.MakeReaction(comment, CommandPutDislike)
	if err != nil {
		log.Println(fmt.Errorf("v1 - CommentPutDislikeHandler - MakeReaction: %w", err))
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
		h.Errors(w, errors)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer")+"#"+strconv.Itoa(int(comment.Id)), http.StatusFound)
}

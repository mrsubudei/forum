package v1

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"forum/internal/entity"
)

func (h *Handler) CreateCommentPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentPageHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/create_comment_page/"+path[len(path)-1] || err != nil || id <= 0 {
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentPageHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	content.Post.Id = int64(id)
	content.Uri = strconv.Itoa(id)

	err = h.ParseAndExecute(w, content, "templates/create_comment.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentPageHandler - ParseAndExecute - %w", err))
	}
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Errors(w, http.StatusMethodNotAllowed)
		return
	}

	var content Content
	var err error
	var ok bool

	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/create_comment/"+path[len(path)-1] || err != nil || id <= 0 {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - URL.Path: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(ImageSizeInt << 20)
	if err != nil {
		h.Errors(w, http.StatusBadRequest)
		return
	}

	content.Uri = strconv.Itoa(id)

	imagePath, err := h.GetImage(w, r)
	if err != nil {
		if strings.Contains(err.Error(), imageTypeForbidden) ||
			strings.Contains(err.Error(), imageTooLarge) {
			w.WriteHeader(http.StatusBadRequest)
			content.ErrorMsg.Message = err.Error()
			err = h.ParseAndExecute(w, content, "templates/create_comment.html")
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - ParseAndExecute #1: %w", err))
			}
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - GetImage: %w", err))
			h.Errors(w, http.StatusInternalServerError)
		}
		return
	}

	if len(r.MultipartForm.Value["content"]) == 0 || len(r.MultipartForm.Value["content"][0]) == 0 {
		h.Errors(w, http.StatusBadRequest)
		return
	}

	content, ok = r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	commentContent := r.MultipartForm.Value["content"][0]
	newComment := entity.Comment{}
	newComment.Content = strings.ReplaceAll(commentContent, "\r\n", "\\n")
	newComment.User = content.User
	newComment.PostId = int64(id)
	newComment.ImagePath = "/" + imagePath

	err = h.Usecases.Comments.WriteComment(newComment)
	if err != nil {
		err = os.Remove(imagePath)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - Remove: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - CreateCommentHandler - WriteComment: %w", err))
			h.Errors(w, http.StatusBadRequest)
			return
		}
	}
	http.Redirect(w, r, "/posts/"+strconv.Itoa(id), http.StatusFound)
}

func (h *Handler) CommentPutLikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutLikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_comment_like/"+path[len(path)-1] || err != nil || id <= 0 {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutLikeHandler - URL.Path: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutLikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = content.User.Id

	err = h.Usecases.Comments.MakeReaction(comment, CommandPutLike)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutLikeHandler - MakeReaction: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer")+"#"+strconv.Itoa(int(comment.Id)), http.StatusFound)
}

func (h *Handler) CommentPutDislikeHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutDislikeHandler - Atoi: %w", err))
	}
	if r.URL.Path != "/put_comment_dislike/"+path[len(path)-1] || err != nil || id <= 0 {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutDislikeHandler - URL.Path: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}

	content, ok := r.Context().Value(Key("content")).(Content)
	if !ok {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutDislikeHandler - TypeAssertion:"+
			"got data of type %T but wanted v1.Content", content))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	comment := entity.Comment{
		Id: int64(id),
	}
	comment.User.Id = content.User.Id

	err = h.Usecases.Comments.MakeReaction(comment, CommandPutDislike)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CommentPutDislikeHandler - MakeReaction: %w", err))
		h.Errors(w, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer")+"#"+strconv.Itoa(int(comment.Id)), http.StatusFound)
}

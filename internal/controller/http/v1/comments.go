package v1

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"forum/internal/entity"

	"forum/pkg/auth"
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

func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	typeSl := strings.Split(mimeType, "/")
	imageType := typeSl[1]

	if imageType != "jpeg" && imageType != "png" && imageType != "gif" {
		return "", errors.New(imageTypeForbidden)
	}

	if _, err := os.Stat("templates/img/storage"); os.IsNotExist(err) {
		os.MkdirAll("templates/img/storage", os.ModePerm)
	}

	manager := auth.NewManager(h.Cfg)
	generated, err := manager.NewToken()
	if err != nil {
		return "", fmt.Errorf("newToken: %w", err)
	}

	path := "templates/img/storage/" + generated + "." + imageType
	targetFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	defer targetFile.Close()

	written, err := io.Copy(targetFile, file)
	if err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}
	if written >= (ImageSizeInt << 20) {
		err := os.Remove(path)
		if err != nil {
			return "", fmt.Errorf("remove: %w", err)
		}
		return "", errors.New(imageTooLarge)
	}

	return path, nil
}

func (h *Handler) CheckSizeExceeded(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open: %w", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return false, fmt.Errorf("decodeConfig: %w", err)
	}
	if image.Width > ImageWidthInt || image.Height > ImageHeightInt {
		return true, nil
	}
	return false, nil
}

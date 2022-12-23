package v1

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"forum/internal/config"
	"forum/internal/usecase"
	"forum/pkg/logger"
)

type Handler struct {
	Usecases *usecase.UseCases
	Cfg      config.Config
	l        *logger.Logger
	Mux      *http.ServeMux
}

func NewHandler(usecases *usecase.UseCases, cfg config.Config, logger *logger.Logger) *Handler {
	mux := http.NewServeMux()
	return &Handler{
		Usecases: usecases,
		Cfg:      cfg,
		l:        logger,
		Mux:      mux,
	}
}

func (h *Handler) ParseAndExecute(w http.ResponseWriter, content Content, path string) error {
	root := getRootPath()
	html, err := template.ParseFiles(root + path)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("parseFiles: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return err
	}
	err = html.Execute(w, content)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("execute: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handler) Errors(w http.ResponseWriter, status int) {
	root := getRootPath()

	errors := ErrMessage{}
	switch status {
	case http.StatusBadRequest:
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
	case http.StatusUnauthorized:
		errors.Code = http.StatusUnauthorized
		errors.Message = ErrStatusNotAuthorized
	case http.StatusForbidden:
		errors.Code = http.StatusForbidden
		errors.Message = ErrLowAccessLevel
	case http.StatusNotFound:
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
	case http.StatusMethodNotAllowed:
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
	case http.StatusNotAcceptable:
		errors.Code = http.StatusNotAcceptable
		errors.Message = UserNotExist
	case http.StatusInternalServerError:
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
	}
	html, err := template.ParseFiles(root + "templates/errors.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - Errors - ParseFiles: %w", err))
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errors.Code)
	err = html.Execute(w, errors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	// main
	router.Handle("/", h.AssignStatus(http.HandlerFunc(h.IndexHandler)))

	// users routes
	router.Handle("/signin_page/", h.AssignStatus(http.HandlerFunc(h.SignInPageHandler)))
	router.Handle("/signup_page/", h.AssignStatus(http.HandlerFunc(h.SignUpPageHandler)))
	router.Handle("/signin/", h.AssignStatus(http.HandlerFunc(h.SignInHandler)))
	router.Handle("/signup/", h.AssignStatus(http.HandlerFunc(h.SignUpHandler)))
	router.Handle("/signout/", h.CheckAuth(http.HandlerFunc(h.SignOutHandler)))
	router.Handle("/edit_profile_page/", h.CheckAuth(http.HandlerFunc(h.EditProfilePageHandler)))
	router.Handle("/edit_profile/", h.CheckAuth(http.HandlerFunc(h.EditProfileHandler)))
	router.Handle("/users/", h.AssignStatus(http.HandlerFunc(h.UserPageHandler)))
	router.Handle("/all_users_page/", h.AssignStatus(http.HandlerFunc(h.AllUsersPageHandler)))
	router.Handle("/find_reacted_users/", h.CheckAuth(http.HandlerFunc(h.FindReactedUsersHandler)))

	// searching routes
	router.Handle("/search_page/", h.AssignStatus(http.HandlerFunc(h.SearchPageHandler)))
	router.Handle("/search/", h.AssignStatus(http.HandlerFunc(h.SearchHandler)))

	// posts routes
	router.Handle("/create_category_page/", h.CheckAuth(http.HandlerFunc(h.CreateCategoryPageHandler)))
	router.Handle("/create_category/", h.CheckAuth(http.HandlerFunc(h.CreateCategoryHandler)))
	router.Handle("/categories/", h.AssignStatus(http.HandlerFunc(h.SearchByCategoryHandler)))
	router.Handle("/posts/", h.AssignStatus(http.HandlerFunc(h.PostPageHandler)))
	router.Handle("/create_post_page/", h.CheckAuth(http.HandlerFunc(h.CreatePostPageHandler)))
	router.Handle("/create_post/", h.CheckAuth(http.HandlerFunc(h.CreatePostHandler)))
	router.Handle("/find_posts/", h.CheckAuth(http.HandlerFunc(h.FindPostsHandler)))
	router.Handle("/put_post_like/", h.CheckAuth(http.HandlerFunc(h.PostPutLikeHandler)))
	router.Handle("/put_post_dislike/", h.CheckAuth(http.HandlerFunc(h.PostPutDislikeHandler)))

	// comments routes
	router.Handle("/create_comment_page/", h.CheckAuth(http.HandlerFunc(h.CreateCommentPageHandler)))
	router.Handle("/create_comment/", h.CheckAuth(http.HandlerFunc(h.CreateCommentHandler)))
	router.Handle("/put_comment_like/", h.CheckAuth(http.HandlerFunc(h.CommentPutLikeHandler)))
	router.Handle("/put_comment_dislike/", h.CheckAuth(http.HandlerFunc(h.CommentPutDislikeHandler)))

	// fileserver
	router.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	router.Handle("/templates/img/", http.StripPrefix("/templates/img/", http.FileServer(http.Dir("templates/img"))))
}

func getRootPath() string {
	_, basePath, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(filepath.Dir(basePath), "/")
	rootSlice := pathSlice[:len(pathSlice)-4]
	root := strings.Join(rootSlice, "/")
	return root + "/"
}

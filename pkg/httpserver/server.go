package httpserver

import (
	"context"
	v1 "forum/internal/controller/http/v1"
	"net/http"
	"time"
)

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 5 * time.Second
	DefaultAddr            = ":8087"
	DefaultShutdownTimeout = 3 * time.Second
	ShutdownTimeout        = 5 * time.Second
)

type Server struct {
	httpServer *http.Server
	handler    *v1.Handler
}

func NewServer(handler *v1.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         DefaultAddr,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		},
		handler: handler,
	}
}

func (s *Server) Run() error {
	http.Handle("/", s.handler.AssignStatus(http.HandlerFunc(s.handler.IndexHandler)))
	http.Handle("/signin_page/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SignInPageHandler)))
	http.Handle("/signup_page/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SignUpPageHandler)))
	http.Handle("/signin/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SignInHandler)))
	http.Handle("/signup/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SignUpHandler)))
	http.Handle("/signout/", s.handler.CheckAuth(http.HandlerFunc(s.handler.SignOutHandler)))
	http.Handle("/edit_profile_page/", s.handler.CheckAuth(http.HandlerFunc(s.handler.EditProfilePageHandler)))
	http.Handle("/edit_profile/", s.handler.CheckAuth(http.HandlerFunc(s.handler.EditProfileHandler)))
	http.Handle("/users/", s.handler.AssignStatus(http.HandlerFunc(s.handler.UserPageHandler)))
	http.Handle("/all_users_page/", s.handler.AssignStatus(http.HandlerFunc(s.handler.AllUsersPageHandler)))
	http.Handle("/search_page/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SearchPageHandler)))
	http.Handle("/search/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SearchHandler)))
	http.Handle("/create_category_page/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreateCategoryPageHandler)))
	http.Handle("/create_category/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreateCategoryHandler)))
	http.Handle("/categories/", s.handler.AssignStatus(http.HandlerFunc(s.handler.SearchByCategoryHandler)))
	http.Handle("/posts/", s.handler.AssignStatus(http.HandlerFunc(s.handler.PostPageHandler)))
	http.Handle("/create_post_page/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreatePostPageHandler)))
	http.Handle("/create_post/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreatePostHandler)))
	http.Handle("/find_posts/", s.handler.CheckAuth(http.HandlerFunc(s.handler.FindPostsHandler)))
	http.Handle("/put_post_like/", s.handler.CheckAuth(http.HandlerFunc(s.handler.PostPutLikeHandler)))
	http.Handle("/put_post_dislike/", s.handler.CheckAuth(http.HandlerFunc(s.handler.PostPutDislikeHandler)))
	http.Handle("/create_comment_page/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreateCommentPageHandler)))
	http.Handle("/create_comment/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CreateCommentHandler)))
	http.Handle("/put_comment_like/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CommentPutLikeHandler)))
	http.Handle("/put_comment_dislike/", s.handler.CheckAuth(http.HandlerFunc(s.handler.CommentPutDislikeHandler)))
	http.Handle("/find_reacted_users/", s.handler.CheckAuth(http.HandlerFunc(s.handler.FindReactedUsersHandler)))
	http.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	http.Handle("/templates/img/", http.StripPrefix("/templates/img/", http.FileServer(http.Dir("templates/img"))))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

package httpserver

import (
	"context"
	"net/http"
	"time"

	v1 "forum/internal/controller/http/v1"
)

type Server struct {
	httpServer *http.Server
	h          *v1.Handler
}

const DefaultTime = int(time.Second)

func NewServer(handler *v1.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         handler.Cfg.Server.Port,
			ReadTimeout:  time.Duration(handler.Cfg.Server.ReadTimeout * DefaultTime),
			WriteTimeout: time.Duration(handler.Cfg.Server.WriteTimeout * DefaultTime),
		},
		h: handler,
	}
}

func (s *Server) Run() error {
	http.Handle("/", s.h.AssignStatus(http.HandlerFunc(s.h.IndexHandler)))
	http.Handle("/signin_page/", s.h.AssignStatus(http.HandlerFunc(s.h.SignInPageHandler)))
	http.Handle("/signup_page/", s.h.AssignStatus(http.HandlerFunc(s.h.SignUpPageHandler)))
	http.Handle("/signin/", s.h.AssignStatus(http.HandlerFunc(s.h.SignInHandler)))
	http.Handle("/signup/", s.h.AssignStatus(http.HandlerFunc(s.h.SignUpHandler)))
	http.Handle("/signout/", s.h.CheckAuth(http.HandlerFunc(s.h.SignOutHandler)))
	http.Handle("/edit_profile_page/", s.h.CheckAuth(http.HandlerFunc(s.h.EditProfilePageHandler)))
	http.Handle("/edit_profile/", s.h.CheckAuth(http.HandlerFunc(s.h.EditProfileHandler)))
	http.Handle("/users/", s.h.AssignStatus(http.HandlerFunc(s.h.UserPageHandler)))
	http.Handle("/all_users_page/", s.h.AssignStatus(http.HandlerFunc(s.h.AllUsersPageHandler)))
	http.Handle("/search_page/", s.h.AssignStatus(http.HandlerFunc(s.h.SearchPageHandler)))
	http.Handle("/search/", s.h.AssignStatus(http.HandlerFunc(s.h.SearchHandler)))
	http.Handle("/create_category_page/", s.h.CheckAuth(http.HandlerFunc(s.h.CreateCategoryPageHandler)))
	http.Handle("/create_category/", s.h.CheckAuth(http.HandlerFunc(s.h.CreateCategoryHandler)))
	http.Handle("/categories/", s.h.AssignStatus(http.HandlerFunc(s.h.SearchByCategoryHandler)))
	http.Handle("/posts/", s.h.AssignStatus(http.HandlerFunc(s.h.PostPageHandler)))
	http.Handle("/create_post_page/", s.h.CheckAuth(http.HandlerFunc(s.h.CreatePostPageHandler)))
	http.Handle("/create_post/", s.h.CheckAuth(http.HandlerFunc(s.h.CreatePostHandler)))
	http.Handle("/find_posts/", s.h.CheckAuth(http.HandlerFunc(s.h.FindPostsHandler)))
	http.Handle("/put_post_like/", s.h.CheckAuth(http.HandlerFunc(s.h.PostPutLikeHandler)))
	http.Handle("/put_post_dislike/", s.h.CheckAuth(http.HandlerFunc(s.h.PostPutDislikeHandler)))
	http.Handle("/create_comment_page/", s.h.CheckAuth(http.HandlerFunc(s.h.CreateCommentPageHandler)))
	http.Handle("/create_comment/", s.h.CheckAuth(http.HandlerFunc(s.h.CreateCommentHandler)))
	http.Handle("/put_comment_like/", s.h.CheckAuth(http.HandlerFunc(s.h.CommentPutLikeHandler)))
	http.Handle("/put_comment_dislike/", s.h.CheckAuth(http.HandlerFunc(s.h.CommentPutDislikeHandler)))
	http.Handle("/find_reacted_users/", s.h.CheckAuth(http.HandlerFunc(s.h.FindReactedUsersHandler)))
	http.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	http.Handle("/templates/img/", http.StripPrefix("/templates/img/", http.FileServer(http.Dir("templates/img"))))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.h.Cfg.Server.ShutDownTimeout*DefaultTime))
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

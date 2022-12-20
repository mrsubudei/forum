package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forum/internal/config"
	v1 "forum/internal/controller/http/v1"
	"forum/internal/repository"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
	"forum/pkg/httpserver"
	"forum/pkg/logger"
	"forum/pkg/sqlite3"
)

func Run(cfg config.Config) {
	// Logger
	l := logger.New()

	// Sqlite
	sq, err := sqlite3.New("database/forum.db")
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - sqlite3.New: %w", err))
		return
	}
	defer sq.Close()

	// Repository
	repo := repository.NewRepositories(sq)
	err = sqlite.CreateDB(sq)
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - NewRepositories: %w", err))
		return
	}

	// Dependencies
	hasher := hasher.NewBcryptHasher()
	tokenManager := auth.NewManager(cfg)

	// Usecases
	postsUseCase := usecase.NewPostsUseCase(repo.Posts, repo.Users, repo.Comments)
	usersUseCase := usecase.NewUsersUseCase(repo.Users, hasher, tokenManager, repo.Posts, repo.Comments)
	commentsUseCase := usecase.NewCommentsUseCase(repo.Comments, repo.Posts, repo.Users)
	useCases := usecase.NewUseCases(postsUseCase, usersUseCase, commentsUseCase)

	// Http
	handler := v1.NewHandler(useCases, cfg, l)
	server := httpserver.NewServer(handler)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	fmt.Printf("Server started at http://%s%s\n", cfg.Server.Host, cfg.Server.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	err = server.Shutdown()
	if err != nil {
		l.WriteLog(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

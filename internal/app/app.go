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
	"forum/pkg/sqlite3"
)

func Run(cfg config.Config) {
	// Logger
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Println(fmt.Errorf("app - Run - os.OpenFile: %w", err))
		return
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Sqlite
	sq, err := sqlite3.New()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - sqlite3.New: %w", err))
		return
	}
	defer sq.Close()

	// Repository
	repositories := repository.NewRepositories(sq)
	err = sqlite.CreateDB(sq)
	if err != nil {
		log.Println(fmt.Errorf("app - Run - NewRepositories: %w", err))
		return
	}

	// Dependencies
	hasher := hasher.NewBcryptHasher()
	tokenManager := auth.NewManager(cfg)

	// Usecases
	useCases := usecase.NewUseCases(usecase.Dependencies{
		Repos:        repositories,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	// Http
	handler := v1.NewHandler(useCases, cfg)
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
		fmt.Printf("app - Run - httpServer.Shutdown: %s\n", err.Error())
	}
}

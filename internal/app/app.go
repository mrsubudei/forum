package app

import (
	"errors"
	"fmt"
	v1 "forum/internal/controller/http/v1"
	"forum/internal/repository"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
	"forum/pkg/httpserver"
	"forum/pkg/sqlite3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	sq, err := sqlite3.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sq.Close()
	repositories := repository.NewRepositories(sq)
	err = sqlite.CreateDB(sq)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	hasher := hasher.NewBcryptHasher()
	tokenManager, err := auth.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	useCases := usecase.NewUseCases(usecase.Dependencies{
		Repos:        repositories,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})
	handler := v1.NewHandler(useCases)
	server := httpserver.NewServer(handler)
	fmt.Println("starting server..")

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	fmt.Println("Server started at port 8080")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	err = server.Shutdown()
	if err != nil {
		fmt.Printf("app - Run - httpServer.Shutdown: %s\n", err.Error())
	}
}

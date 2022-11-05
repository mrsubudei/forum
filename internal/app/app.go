package app

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"forum/internal/repository/sqlite"
	"forum/internal/usecase"
	"forum/pkg/auth"
	"forum/pkg/hasher"
	"forum/pkg/sqlite3"
	"log"
	"strings"
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
	hasher := hasher.NewBcryptHasher()
	tokenManager, err := auth.NewManager("s")
	if err != nil {
		log.Fatal(err)
	}
	useCases := usecase.NewUseCases(usecase.Dependencies{
		Repos:        repositories,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	user := entity.User{
		Name:     "uff",
		Email:    "uff",
		Password: "uff",
	}
	err = useCases.Users.SignUp(user)
	if err != nil {
		sl := strings.Split(err.Error(), ":")
		for i := 0; i < len(sl); i++ {
			fmt.Println(strings.Trim(sl[i], " "))

		}
	}

}

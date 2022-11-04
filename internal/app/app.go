package app

import (
	"forum/internal/usecase/repo/sqlite"
	"forum/pkg/sqlite3"

	"log"
)

func Run() {
	sq, err := sqlite3.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sq.Close()
	CommunicationUseCase := sqlite.New(sq)

	err = CommunicationUseCase.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
}

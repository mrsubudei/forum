package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func WriteToUsers() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(
		`INSERT INTO users(name, email, password, date, city, sex) 
			values(?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	name := "Sui"
	email := "Sbi@m.fo"
	pass := "vivseurodi"
	date := "19.06.1989"
	city := "Astana"
	sex := "male"

	cached := cachePassword(pass)

	_, err = stmt.Exec(name, email, cached, date, city, sex)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func cachePassword(str string) []byte {
	cached, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return cached
}

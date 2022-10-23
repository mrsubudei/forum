package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(str1, str2 string) {
	db, err := sql.Open("sqlite3", "./forum.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// sqlStmt := `
	// create table users (id integer not null primary key, name text, email text);
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into users(name, email) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(str1, str2)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	users := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT,
		date TEXT,
		city TEXT,
		sex TEXT
		);
	`
	_, err = db.Exec(users)
	if err != nil {
		log.Printf("%q: %s\n", err, users)
		return nil
	}

	posts := `
	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		date TEXT NOT NULL,
		message TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(posts)
	if err != nil {
		log.Printf("%q: %s\n", err, posts)
		return nil
	}

	comments := `
	CREATE TABLE comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		date TEXT NOT NULL,
		message TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(comments)
	if err != nil {
		log.Printf("%q: %s\n", err, comments)
		return nil
	}

	postlikes := `
	CREATE TABLE postlikes (
		post_id INTEGER,
		user_id INTEGER,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(postlikes)
	if err != nil {
		log.Printf("%q: %s\n", err, postlikes)
		return nil
	}

	commentlikes := `
	CREATE TABLE commentlikes (
		comment_id INTEGER,
		user_id INTEGER,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(commentlikes)
	if err != nil {
		log.Printf("%q: %s\n", err, commentlikes)
		return nil
	}

	return db
}

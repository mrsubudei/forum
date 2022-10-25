package usecases

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
		reg_date TEXT,
		date_of_birth TEXT,
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
		content TEXT NOT NULL,
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
		content TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(comments)
	if err != nil {
		log.Printf("%q: %s\n", err, comments)
		return nil
	}

	postLikes := `
	CREATE TABLE post_likes (
		post_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(postLikes)
	if err != nil {
		log.Printf("%q: %s\n", err, postLikes)
		return nil
	}

	postDisLikes := `
	CREATE TABLE post_dislikes (
		post_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(postDisLikes)
	if err != nil {
		log.Printf("%q: %s\n", err, postDisLikes)
		return nil
	}

	commentLikes := `
	CREATE TABLE comment_likes (
		comment_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(commentLikes)
	if err != nil {
		log.Printf("%q: %s\n", err, commentLikes)
		return nil
	}

	commentDisLikes := `
	CREATE TABLE comment_dislikes (
		comment_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(commentDisLikes)
	if err != nil {
		log.Printf("%q: %s\n", err, commentDisLikes)
		return nil
	}

	topics := `
	CREATE TABLE topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
		);
	`
	_, err = db.Exec(topics)
	if err != nil {
		log.Printf("%q: %s\n", err, topics)
		return nil
	}

	referencetopic := `
	CREATE TABLE referencetopic (
		post_id TEXT,
		topic_id TEXT,
		PRIMARY KEY(post_id, topic_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (topic_id) REFERENCES topics(id)
		);
	`
	_, err = db.Exec(referencetopic)
	if err != nil {
		log.Printf("%q: %s\n", err, referencetopic)
		return nil
	}

	return db
}

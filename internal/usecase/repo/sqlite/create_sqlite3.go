package sqlite

import (
	"forum/pkg/sqlite3"
)

type CommunicationRepo struct {
	*sqlite3.Sqlite
}

func New(sq *sqlite3.Sqlite) *CommunicationRepo {
	return &CommunicationRepo{sq}
}

func (c *CommunicationRepo) CreateDB() error {

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
	_, err := c.DB.Exec(users)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(posts)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(comments)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(postLikes)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(postDisLikes)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(commentLikes)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(commentDisLikes)
	if err != nil {
		return err
	}

	topics := `
	CREATE TABLE topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
		);
	`
	_, err = c.DB.Exec(topics)
	if err != nil {
		return err
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
	_, err = c.DB.Exec(referencetopic)
	if err != nil {
		return err
	}
	return nil
}

package sqlite

import (
	"forum/pkg/sqlite3"
)

func CreateDB(s *sqlite3.Sqlite) error {

	users := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT,
		reg_date TEXT,
		date_of_birth TEXT,
		city TEXT,
		sex TEXT,
		session_token TEXT,
		session_ttl TEXT
		);
	`

	_, err := s.DB.Exec(users)
	if err != nil {
		return err
	}

	posts := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = s.DB.Exec(posts)
	if err != nil {
		return err
	}

	comments := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		date TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	_, err = s.DB.Exec(comments)
	if err != nil {
		return err
	}

	postLikes := `
	CREATE TABLE IF NOT EXISTS post_likes (
		post_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = s.DB.Exec(postLikes)
	if err != nil {
		return err
	}

	postDisLikes := `
	CREATE TABLE IF NOT EXISTS post_dislikes (
		post_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	_, err = s.DB.Exec(postDisLikes)
	if err != nil {
		return err
	}
	commentLikes := `
	CREATE TABLE IF NOT EXISTS comment_likes (
		comment_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	_, err = s.DB.Exec(commentLikes)
	if err != nil {
		return err
	}
	commentDisLikes := `
	CREATE TABLE IF NOT EXISTS comment_dislikes (
		comment_id INTEGER,
		user_id INTEGER,
		date TEXT,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	_, err = s.DB.Exec(commentDisLikes)
	if err != nil {
		return err
	}

	topics := `
	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
		);
	`
	_, err = s.DB.Exec(topics)
	if err != nil {
		return err
	}
	referencetopic := `
	CREATE TABLE IF NOT EXISTS reference_topic (
		post_id INTEGER,
		topic TEXT,
		PRIMARY KEY (post_id, topic),
		FOREIGN KEY (post_id) REFERENCES posts(id)
		);
	`
	_, err = s.DB.Exec(referencetopic)
	if err != nil {
		return err
	}

	moderators := `
	CREATE TABLE IF NOT EXISTS moderators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = s.DB.Exec(moderators)
	if err != nil {
		return err
	}

	return nil
}

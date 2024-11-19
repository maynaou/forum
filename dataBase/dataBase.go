package dataBase

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("db error 1") // still need better error handling
		return
	}

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = Db.Exec(createUsersTable)
	if err != nil {
		fmt.Println("Error creating users table:", err)
	}

	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER NOT NULL,
		session_id TEXT PRIMARY KEY,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err = Db.Exec(createSessionsTable)
	if err != nil {
		fmt.Println("Error creating sessions table:", err)
	}

	const CreatPosts string = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_creator TEXT NOT NULL,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_creator) REFERENCES users(username)
	);`

	if _, err := Db.Exec(CreatPosts); err != nil {
		fmt.Println(err)
		return
	}

	const CommentsTable string = `
	CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_body TEXT NOT NULL,
		comment_writer TEXT NOT NULL,
		post_commented_id INTEGER NOT NULL,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_commented_id) REFERENCES posts(id)
	)`
	if _, err := Db.Exec(CommentsTable); err != nil {
		fmt.Println(err)
		return
	}

	const CreatCategories string = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		categorie TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id)
	)`

	if _, err := Db.Exec(CreatCategories); err != nil {
		fmt.Println(err)
		return
	}

	const CreateCommentsIndex string = `
	CREATE INDEX IF NOT EXISTS idx_post_commented_id ON comments(post_commented_id);`

	_, err = Db.Exec(CreateCommentsIndex)
	if err != nil {
		fmt.Println("Error creating index on post_commented_id:", err)
		return
	}
}

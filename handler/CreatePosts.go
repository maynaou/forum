package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	data "main/dataBase"
)

func CreatPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create_post" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed Method", http.StatusMethodNotAllowed)
		return
	}
	CurrentUser := r.URL.Query().Get("user")
	Post_id, _ := strconv.Atoi(r.URL.Query().Get("postid"))
	// fmt.Println("Create post function post id is :", Post_id+1, "and the writer is :", CurrentUser)
	title := r.FormValue("title")
	body := r.FormValue("body")
	// categories := r.FormValue("categories")
	categories := r.Form["categories"]
	// fmt.Println(categories)
	if len(categories) == 0 {
		categories = append(categories, "All")
	}
	if title == "" || body == "" {
		http.Error(w, "bad request empty post", http.StatusBadRequest)
		return
	}
	row := data.Db.QueryRow("SELECT username FROM users WHERE username = ?", CurrentUser)
	var username string
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("this user don't exist")
			http.Error(w, "you are in the guest session", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("we can't retrive data")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	// Add The post to the posts table
	_, err = data.Db.Exec("INSERT INTO posts(post_creator, title, body) VALUES (?, ?, ?)", CurrentUser, title, body)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	for _, categorie := range categories {
		_, err = data.Db.Exec("INSERT INTO categories(post_id, categorie) VALUES (?, ?)", Post_id+1, categorie)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Internal server error", 500)
			return
		}
	}

	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

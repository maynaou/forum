package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"main/dataBase"
)

func CreateComments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		commentBody := r.FormValue("comment_body")
		postID, err := strconv.Atoi(r.FormValue("post_id"))
		commentWriter := r.FormValue("comment_writer")
		fmt.Println(commentBody, commentWriter, postID)
		if commentBody == "" {
			fmt.Println("hhhhh")
			http.Error(w, "Bad Request: Comment cannot be empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		_, err = dataBase.Db.Exec("INSERT INTO comments(comment_body, comment_writer, post_commented_id) VALUES (?, ?, ?)", commentBody, commentWriter, postID)
		if err != nil {
			log.Println("Error inserting comment:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/forum?user="+commentWriter, http.StatusSeeOther)
	}
}

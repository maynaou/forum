package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	data "main/dataBase"
)

type Post struct {
	Postid            int
	Usernamepublished string
	CurrentUsser      string
	Title             string
	Body              string
	Time              any
	Categorie         string
	Comments          []Comment
}

type Comment struct {
	Comment_id        int
	Comment_body      string
	Comment_writer    string
	Post_commented_id int
	Comment_time      any
}

func Forum(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Panic récupéré : %v", rec)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}()

	if r.URL.Path != "/forum" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/forum.html")
	if err != nil {
		http.Error(w, "Internal Server Error with forum html page", http.StatusInternalServerError)
		return
	}
	var CurrentUser, CurrentSession string
	var session_id string
	cat_to_filter := r.FormValue("categories")
	cookie1, err1 := r.Cookie("session_token")
	cookie2, err2 := r.Cookie("user_token")

	if err1 != nil || err2 != nil {
		cookie3, err3 := r.Cookie("guest_token")
		if err3 != nil {
			log.Printf("Cookie 'guest_token' absent : %v", err3)
			// Passer en mode invité
			CurrentUser = "Guest"
			CurrentSession = "0"
		}
		CurrentUser = cookie3.Value
		CurrentSession = "0"
	} else {
		CurrentUser = cookie2.Value
		CurrentSession = cookie1.Value
		err = data.Db.QueryRow("SELECT session_id FROM sessions WHERE session_id = ?", CurrentSession).Scan(&session_id)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	posts_toshow := GetPosts(cat_to_filter, tmpl, w, CurrentUser)
	err = tmpl.Execute(w, struct {
		Currenuser string
		Posts      []Post
	}{
		Currenuser: CurrentUser,
		Posts:      posts_toshow,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
}

func GetPosts(cat_to_filter string, tmpl *template.Template, w http.ResponseWriter, CurrentUser string) []Post {
	var post_rows *sql.Rows
	var err error
	if cat_to_filter != "all" && cat_to_filter != "" {
		// post_rows, err = data.Db.Query("SELECT post_id FROM categories WHERE categorie = ?;", cat_to_filter)
		post_rows, err = data.Db.Query(`
			SELECT posts.* FROM posts
			JOIN categories ON posts.id = categories.post_id
			WHERE categories.categorie = ?`, cat_to_filter)
	} else {
		post_rows, err = data.Db.Query("SELECT * FROM posts;")
	}
	if err != nil {
		if err == sql.ErrNoRows {
			if err := tmpl.Execute(w, CurrentUser); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Template execution error: %v", err)
			}
			return nil
		}
		fmt.Println(err)
		http.Error(w, "ana hna Internal server error", http.StatusInternalServerError)
		return nil
	}
	defer post_rows.Close()
	// post_rows, err = data.Db.Query("SELECT * FROM posts WHERE id = ?;", cat_to_filter)
	var posts_toshow []Post
	var post Post
	for post_rows.Next() {

		var id int
		var title, body, usernamepublished string
		var time any
		if err := post_rows.Scan(&id, &usernamepublished, &title, &body, &time); err != nil {
			fmt.Println(err)
			continue
		}
		post.Comments = GetComments(id, 10, 0)
		// fmt.Println("comments id= ", comment_id, "post id= ", post_id)
		posts_toshow = append(posts_toshow, Post{
			Postid:            id,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
			Title:             title,
			Body:              body,
			Time:              time,
			Comments:          post.Comments,
		})
	}

	if err := post_rows.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(posts_toshow)-1; i++ {
		for j := i + 1; j < len(posts_toshow); j++ {
			posts_toshow[i], posts_toshow[j] = posts_toshow[j], posts_toshow[i]
		}
	}
	return posts_toshow
}

func GetComments(postID, limit, offset int) []Comment {
	if postID <= 0 || limit <= 0 || offset < 0 {
		log.Printf("Invalid parameters: postID=%d, limit=%d, offset=%d", postID, limit, offset)
		return nil
	}

	rows, err := data.Db.Query("SELECT comment_body, comment_writer, post_commented_id, time FROM comments WHERE post_commented_id = ? LIMIT ? OFFSET ?;", postID, limit, offset)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Comment_body, &comment.Comment_writer, &comment.Post_commented_id, &comment.Comment_time)
		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}

		if len(comment.Comment_body) > 500 {
			log.Printf("Commentaire trop long ignoré : %v", comment.Comment_body)
			continue
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating comments: %v", err)
	}
	return comments
}

package main

import (
	"fmt"
	"net/http"

	handler "main/handler"
)

var port = "9532"

func main() {
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/guest", handler.Guest)
	http.HandleFunc("/register", handler.HandleRegistration)
	http.HandleFunc("/forum", handler.Forum) // this is where the forum would be handled after the login
	http.HandleFunc("/create_post", handler.CreatPost)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/create_comment", handler.CreateComments)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

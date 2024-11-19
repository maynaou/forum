package handler

import (
	"fmt"
	"net/http"

	data "main/dataBase"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	CurrentCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// fmt.Println("Cookie", CurrentCookie.Value)
	CV := CurrentCookie.Value
	cookie := &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	http.SetCookie(w,cookie)
	_, err = data.Db.Exec("DELETE FROM sessions WHERE session_id = ?", CV)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

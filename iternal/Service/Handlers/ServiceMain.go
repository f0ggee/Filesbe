package Handlers

import (
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("KEY"))

//func F(w http.ResponseWriter, r *http.Request) {
//
//	session, err := store.Get(r, "token6")
//	if err != nil {
//		slog.Error("cookie don't send", err)
//		http.Error(w, "Cookie dont sent", http.StatusUnauthorized)
//		return
//	}
//
//	fmt.Println(session)
//
//	_, ok := session.Values["cookie"]
//	if !ok {
//		http.Redirect(w, r, "login", http.StatusFound)
//		return
//
//	} else {
//		///Here will be logic transfer to home page
//	}
//
//}

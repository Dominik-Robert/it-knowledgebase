package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Username string
	Password string
	// 0 = Schreiberlaubnis || 1 = Admin
	AccessLevel int
}

var mert User = User{Username: "mertayg", Password: "seinpasswort", AccessLevel: 1}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//timon := User{Username: "timonsrm", Password: "meinpasswort", AccessLevel: 1}
	t := template.Must(template.ParseFiles("templates/login.gohtml"))
	t.Execute(w, mert)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> Hello World </h1>")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "Post" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	uname := r.FormValue("username")
	upw := r.FormValue("password")
	fmt.Fprintf(w, "<p>Das ist der angegebene Username "+uname+" und das angegebene Passwort ist "+upw+"</p>")
	if upw == mert.Password {
		fmt.Fprintf(w, "<p>das passwort ist korrekt</p>")
	} else {
		fmt.Fprintf(w, "<p>das passwort ist falsch</p>")
	}
}
func feedHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/feed.gohtml"))
	t.Execute(w, mert)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/detail.gohtml"))
	t.Execute(w, mert)
}

func main() {
	http.HandleFunc("/", feedHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/loginauth/", authHandler)
	http.HandleFunc("/detail/", detailHandler)
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))

	http.ListenAndServe(":80", nil)
}

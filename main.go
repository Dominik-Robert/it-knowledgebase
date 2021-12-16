package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://timon:timonTKR@23.88.103.113:30001/timon"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	http.HandleFunc("/", feedHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/loginauth/", authHandler)
	http.HandleFunc("/detail/", detailHandler)
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("assets/images"))))

	http.ListenAndServe(":80", nil)
}

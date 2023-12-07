package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Article struct {
	Title, Anonse, FullText string
}

var posts = []Article{}
var showPost = Article{}

func init() {
	var err error
	connStr := getDBConnectionString()
	db, err = sql.Open("postgres", connStr)
	fmt.Println("Connection String:", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func getDBConnectionString() string {
	return fmt.Sprintf("user=postgres dbname=articles password=admin host=localhost sslmode=disable")
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	res, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Title, &post.Anonse, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}
	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	anonse := r.FormValue("anonse")
	full_text := r.FormValue("full_text")
	if title == "" || anonse == "" || full_text == "" {
		fmt.Fprintf(w, "Please fill the blans")
	} else {

		statement := "INSERT INTO articles (title, anonse, full_text) VALUES($1, $2, $3)"
		_, err := db.Exec(statement, title, anonse, full_text)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE title = '%s'", vars["title"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Title, &post.Anonse, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post
	}
	t.ExecuteTemplate(w, "show", showPost)

}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_article", save_article).Methods("POST")
	router.HandleFunc("/post/{title}", show_post).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleFunc()
}

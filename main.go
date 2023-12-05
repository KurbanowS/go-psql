package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	connStr := getDBConnectionString()
	// Initialize the database connection pool.
	db, err = sql.Open("postgres", connStr)
	fmt.Println("Connection String:", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func getDBConnectionString() string {
	// user := os.Getenv("postgres")
	// password := os.Getenv("admin")
	// host := os.Getenv("localhost:5432")
	// dbname := os.Getenv("articles")

	return fmt.Sprintf("user=postgres dbname=articles password=admin host=localhost sslmode=disable")
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "index", nil)
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
	fullText := r.FormValue("full_text")

	statement := "INSERT INTO articles (title, anonse, full_text) VALUES($1, $2, $3)"
	_, err := db.Exec(statement, title, anonse, fullText)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleFunc()
}

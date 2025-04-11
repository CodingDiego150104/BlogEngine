package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `validate:"required"`
	Content string `validate:"required"`
	Date    time.Time
}

var (
	posts    []Post
	mu       sync.Mutex
	validate *validator.Validate
	funcMap  = template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
	db *gorm.DB
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home.html").Funcs(funcMap).ParseFiles("templates/home.html"))

	// Pagina corrente
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	postsPerPage := 6
	mu.Lock()
	start := (page - 1) * postsPerPage
	end := start + postsPerPage
	if end > len(posts) {
		end = len(posts)
	}
	paginatedPosts := posts[start:end]

	totalPages := (len(posts) + postsPerPage - 1) / postsPerPage
	if totalPages == 0 {
		totalPages = 1
	}
	mu.Unlock()

	tmpl.Execute(w, struct {
		Posts       []Post
		CurrentPage int
		TotalPages  int
	}{
		Posts:       paginatedPosts,
		CurrentPage: page,
		TotalPages:  totalPages,
	})
}

func newPostFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("new.html").Funcs(funcMap).ParseFiles("templates/new.html"))
	tmpl.Execute(w, nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Errore nel parsing del form", http.StatusBadRequest)
		return
	}

	newPost := Post{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}

	err = validate.Struct(newPost)
	if err != nil {
		http.Error(w, "Titolo e contenuto sono obbligatori", http.StatusBadRequest)
		return
	}

	// Salva il nuovo post nel database
	db.Create(&newPost)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("Errore apertura DB:" + err.Error())
	}

	// Migrazione del database
	err = db.AutoMigrate(&Post{})
	if err != nil {
		panic("Errore nella migrazione del database: " + err.Error())
	}
}

func main() {
	database, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("[error] failed to initialize database, got error %v", err)
	}

	db.Migrate(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", models.HomeHandler)

	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

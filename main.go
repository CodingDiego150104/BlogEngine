package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Post struct {
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

	postsPerPage := 5
	mu.Lock()
	start := (page - 1) * postsPerPage
	end := start + postsPerPage
	if end > len(posts) {
		end = len(posts)
	}
	paginatedPosts := posts[start:end]
	totalPages := (len(posts) + postsPerPage - 1) / postsPerPage
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

	mu.Lock()
	posts = append([]Post{newPost}, posts...)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	validate = validator.New()

	r := chi.NewRouter()

	// File statici (CSS, JS, immagini, ecc.)
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	// Rotte
	r.Get("/", homeHandler)
	r.Get("/new", newPostFormHandler)
	r.Post("/create", createPostHandler)

	http.ListenAndServe(":8080", r)
}

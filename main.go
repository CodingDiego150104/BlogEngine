package main

import (
	"blog/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite" // questo è importante per usare sqlite senza CGO
)

var ( // variabili globali
	// db è il puntatore al database
	db       *gorm.DB
	validate *validator.Validate
	funcMap  = template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
)

type Post struct { // struttura del post
	ID      uint   `gorm:"primaryKey"`
	Title   string `validate:"required"`
	Content string `validate:"required"`
	Date    time.Time
}

var templates *template.Template

func initDB() { // funzione che inizializza il database
	var err error
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Errore apertura DB:", err)
	}

	err = db.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Errore nella migrazione:", err)
	}
}

func initTemplates() {
	// Carica tutti i template HTML nella cartella "templates"
	var err error
	templates, err = template.New("base").Funcs(funcMap).ParseFiles(
		"templates/home.html",
		"templates/new.html",
		"templates/post.html", // Aggiungi anche il template post.html
	)
	if err != nil {
		log.Fatal("Errore nel caricamento dei template:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) { // funzione che gestisce la pagina principale
	// Controlla se il metodo è GET
	tmpl := template.Must(template.New("home.html").Funcs(funcMap).ParseFiles("templates/home.html"))

	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	postsPerPage := 6
	var posts []Post

	// Recupera post ordinati per data decrescente
	result := db.Order("date DESC").Offset((page - 1) * postsPerPage).Limit(postsPerPage).Find(&posts)
	if result.Error != nil {
		http.Error(w, "Errore nel recupero dei post", http.StatusInternalServerError)
		return
	}

	// Conta il numero totale di post
	var total int64
	db.Model(&Post{}).Count(&total)
	totalPages := int((total + int64(postsPerPage) - 1) / int64(postsPerPage))
	if totalPages == 0 {
		totalPages = 1
	}

	tmpl.Execute(w, struct {
		Posts       []Post
		CurrentPage int
		TotalPages  int
	}{
		Posts:       posts,
		CurrentPage: page,
		TotalPages:  totalPages,
	})
}

// POST HANDLER
func newPostFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("new.html").Funcs(funcMap).ParseFiles("templates/new.html"))
	tmpl.Execute(w, nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) { // funzione che gestisce la creazione di un nuovo post
	// Controlla se il metodo è POST
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

	result := db.Create(&newPost)
	if result.Error != nil {
		http.Error(w, "Errore nel salvataggio del post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HANDLER COMMENTI
func showPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID del post non valido", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := db.Preload("Comments").First(&post, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "post.html", post)
	if err != nil {
		http.Error(w, "Errore rendering template", http.StatusInternalServerError)
	}
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID del post non valido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Errore nei dati del form", http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		PostID:    uint(id),
		Author:    r.FormValue("author"),
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		http.Error(w, "Errore salvataggio commento", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

// REINDIRIZZAMENTO
func open(url string) error { //funzione che permette di aprire un link nel browser
	var cmd string
	var param []string

	switch runtime.GOOS { //prende il valore del sistema operativo
	case "windows":
		cmd = "cmd"
		param = []string{"/c", "start"}
	case "darwin": //macos
		cmd = "open"
	default: //linux e similari
		cmd = "xdg-open"
	}
	param = append(param, url)
	return exec.Command(cmd, param...).Start() //esegue il comando
}

// MAIN
func main() {
	// Inizializza il database e il validatore
	initDB()
	initTemplates()
	validate = validator.New()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	r.Get("/", homeHandler)
	r.Get("/new", newPostFormHandler)
	r.Post("/create", createPostHandler)

	// Rotte per i singoli post e commenti
	r.Get("/post/{id}", showPostHandler)
	r.Post("/post/{id}/comment", addCommentHandler)

	log.Println("Server avviato su http://localhost:8080")
	open("http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

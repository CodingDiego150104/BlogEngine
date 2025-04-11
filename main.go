package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	posts    []models.Post
	mu       sync.Mutex
	validate *validator.Validate
	db       *gorm.DB
)

func main() {
	// Log to see if the main function is executed
	log.Println("Starting server...")

	initDB()         // Initialize the DB
	defer db.Close() // Ensure the DB connection is closed at the end
	validate = validator.New()

	// Create the router
	r := chi.NewRouter()

	// Middleware
	r.Use(chi.Middleware.Logger)    // Enable logging middleware
	r.Use(chi.Middleware.Recoverer) // Enable recovery middleware

	// Static file handling (CSS, JS, images)
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	// Define routes
	r.Get("/", homeHandler)
	r.Get("/new", newPostFormHandler)
	r.Post("/create", createPostHandler)

	// Log server startup
	log.Println("Server is running at http://localhost:8080")

	// Start the server
	http.ListenAndServe(":8080", r)
}

func initDB() {
	var err error
	// Open the SQLite database file (will create it if it doesn't exist)
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Migrate the schema (create tables if they don't exist)
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("Error migrating schema: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Fetch all posts from the database
	if err := db.Find(&posts).Error; err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Render the home page with the posts
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, posts)
}

func newPostFormHandler(w http.ResponseWriter, r *http.Request) {
	// Render the form for creating a new post
	tmpl := template.Must(template.ParseFiles("templates/new_post.html"))
	tmpl.Execute(w, nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	r.ParseForm()
	title := r.FormValue("title")
	body := r.FormValue("body")

	// Validate the form data
	post := models.Post{
		Title: title,
		Body:  body,
	}

	err := validate.Struct(post)
	if err != nil {
		// Validation failed, show the form again with the errors
		log.Printf("Validation failed: %v", err)
		http.Redirect(w, r, "/new", http.StatusFound)
		return
	}

	// Save the new post to the database
	if err := db.Create(&post).Error; err != nil {
		log.Printf("Error saving post: %v", err)
		http.Error(w, "Error saving post", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page after successful post creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"blog/models"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usa: blogctl [start|shutdown|restart|dbreset|backup]")
		return
	}

	switch os.Args[1] {
	case "start":
		runGoModTidyIfNeeded()
		startServer()
	case "shutdown":
		shutdownServer()
	case "restart":
		shutdownServer()
		time.Sleep(2 * time.Second) // Pausa tra stop e start
		runGoModTidyIfNeeded()
		startServer()
	case "dbreset":
		shutdownServer()
		dbclean_updated()
		runGoModTidyIfNeeded()
		startServer()
	case "backup":
		shutdownServer()
		dbBackup()
		runGoModTidyIfNeeded()
		startServer()
	default:
		fmt.Println("Comando non valido.")
	}
}

func runGoModTidyIfNeeded() { // Se le dipendenze non sono installate, esegui 'go mod tidy'
	// Controlla se il file go.sum esiste. Se non esiste, esegue 'go mod tidy'
	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		fmt.Println("Prima esecuzione rilevata, eseguo 'go mod tidy'...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Errore durante 'go mod tidy': %v", err)
		}
	}
}

func startServer() { // Avvia il server Go
	// Controlla se il file main.go esiste. Se non esiste, mostra un messaggio di errore
	fmt.Println("Avvio del server...")
	switch runtime.GOOS {
	case "windows":
		// Su Windows, avvia il server in una nuova finestra di terminale usando cmd.exe
		cmd := exec.Command("cmd", "/c", "start", "go", "run", "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatalf("Errore durante l'avvio del server: %v", err)
		}
		fmt.Println("Server avviato in una nuova finestra del terminale.")
	default:
		// Su Linux/macOS avvia il server in background
		cmd := exec.Command("nohup", "go", "run", "main.go", "&")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatalf("Errore durante l'avvio del server: %v", err)
		}
		fmt.Println("Server avviato in background.")
	}
}

func shutdownServer() { // Spegne il server Go
	// Per Windows, usa taskkill
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/IM", "main.exe") // Assicurati che il nome dell'eseguibile sia corretto
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Errore durante lo spegnimento del server:", err)
		} else {
			fmt.Println("Server spento.")
		}
	} else {
		// Per Linux/macOS, usa pkill
		cmd := exec.Command("pkill", "-f", "main.go") // Questo funziona su Linux/macOS
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Errore durante lo spegnimento del server:", err)
		} else {
			fmt.Println("Server spento.")
		}
	}
}

func dbclean() {
	if runtime.GOOS == "windows" {
		err := os.Remove("blog.db")
		if err != nil {
			log.Println("Errore durante la cancellazione del database:", err)
		} else {
			fmt.Println("Database cancellato con successo")
		}
	} else {
		cmd := exec.Command("rm", "blog.db")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Errore durante la cancellazione del database:", err)
		} else {
			fmt.Println("Database cancellato con successo")
		}
	}
}

func dbBackup() {

	source, err := os.Open("blog.db")
	if err != nil {
		fmt.Println("Errore nell'apertura del database")
	}

	date := time.Now()
	formatted_date := date.Format("01-02-2006 15-04-05")
	formatted_date = formatted_date + ".db"
	fmt.Println(formatted_date)

	dest, err := os.Create(formatted_date)
	if err != nil {
		fmt.Println("Errore nella creazione del file di destinazione")
	}

	_, err = io.Copy(dest, source)
	if err != nil {
		fmt.Println("Errore nella creazione del backup")
	} else {
		fmt.Println("Backup creato con successo")
	}
}

func dbclean_updated() {
	var err error
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Errore in apertura database", err)
	}

	err = db.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Errore nella migrazione", err)
	}

	var MaxId int
	db.Model(&models.Post{}).Select("MAX(id)").Scan(&MaxId)

	for i := 1; i <= MaxId; i++ {
		db.Delete(models.Post{}, i)
	}
}

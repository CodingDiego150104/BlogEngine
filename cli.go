package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usa: blogctl [start|shutdown|restart]")
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

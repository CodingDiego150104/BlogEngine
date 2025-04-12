# BlogEngine

BlogEngine è un motore di blog scritto in Go, con un frontend semplice in HTML/CSS.  
L’idea è avere qualcosa di leggero, facile da capire e modificare, senza mille dipendenze.

---

## ✨ Cosa fa

- Serve pagine blog scritte in Go
- Ha un frontend con due temi: chiaro e scuro
- C’è un pulsante per cambiare tema al volo
- Valida i dati inseriti nei form (es. per creare post)
- Salva tutto in un database SQLite (nessuna configurazione extra necessaria)

---

## 🧠 Tecnologie usate

### Backend (Go)

Il backend è scritto in Go usando librerie standard + alcune esterne:

- `html/template`: per generare HTML dinamico
- `net/http`: server web base
- `strconv`, `time`: utility varie
- [`chi`](https://github.com/go-chi/chi): router leggero e modulare
- [`validator`](https://github.com/go-playground/validator): validazione dei dati
- [`gorm`](https://gorm.io/): ORM per interagire con il database
- [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite): driver SQLite che funziona senza CGO

> ⚠️ Per far funzionare `modernc.org/sqlite`, è richiesto un compilatore C (`gcc`).

---

### Frontend

- HTML + CSS puro
- Due temi: `style.css` (chiaro) e `dark.css` (scuro)
- Toggle del tema gestito da `theme-toggle.js`

---

## ⚙️ Requisiti

- [Go](https://go.dev/) 1.18 o superiore
- Compilatore C (`gcc`)
  - Su Windows si consiglia **MSYS2**:
    ```bash
    pacman -S mingw-w64-x86_64-gcc
    ```

---

## 🚀 Come si usa

Clona il progetto, installa le dipendenze e avvia il server:

```bash
git clone https://github.com/CodingDiego150104/BlogEngine.git
cd BlogEngine
go mod tidy      # per scaricare tutte le dipendenze
go run main.go   # per avviare il server

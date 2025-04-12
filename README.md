
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
git clone https://github.com/D13GOOOO/BlogEngine.git
cd BlogEngine
go mod tidy      # per scaricare tutte le dipendenze
```

### Avviare, spegnere e riavviare il server con `blogctl`

Nel progetto è incluso un comando a riga di comando chiamato `blogctl`, che ti permette di gestire facilmente il server.

1. **Avviare il server**:
   ```bash
   blogctl start
   ```

2. **Spegnere il server**:
   ```bash
   blogctl shutdown
   ```

3. **Riavviare il server**:
   ```bash
   blogctl restart
   ```

Il comando `blogctl start` avvia il server in una nuova finestra del terminale su **Windows** e in background su **Linux/macOS**, mentre `shutdown` e `restart` fermano il server in esecuzione.

---

## 🐞 Risoluzione dei problemi

Se il comando `blogctl` non funziona correttamente, assicurati di avere Go installato correttamente, e che il file `cli.go` sia stato compilato come eseguibile. Se non è stato compilato, puoi farlo con:

```bash
go build cli.go
```

Il comando `blogctl` sarà quindi disponibile come eseguibile nel tuo terminale.

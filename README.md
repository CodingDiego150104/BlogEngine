# BlogEngine

BlogEngine è un motore di blog scritto in Go, con un frontend semplice in HTML/CSS.  
L’idea è avere qualcosa di leggero, facile da capire e modificare, senza mille dipendenze.

---

## Cosa fa

- Serve pagine blog scritte in Go
- Ha un frontend con due temi: chiaro e scuro
- C’è un pulsante per cambiare tema al volo
- Valida i dati inseriti nei form (es. per creare post) in modo pulito

---

## Backend (Go)

Il backend è scritto in Go usando alcune librerie standard e un paio di esterne:

- `html/template`: per generare le pagine HTML dinamiche
- `net/http`: server web base
- `sync`, `strconv`, `time`: utility varie
- [`chi`](https://github.com/go-chi/chi): un router leggero e modulare, ottimo per gestire le rotte in modo chiaro
- [`validator`](https://github.com/go-playground/validator): usato per controllare che i dati ricevuti (tipo da un form) siano validi — molto utile e pulito da integrare

---

## Frontend

- Scritto in HTML e CSS
- Ci sono due file CSS: `style.css` (tema chiaro) e `dark.css` (tema scuro)
- C’è anche uno script `theme-toggle.js` che permette di passare da un tema all’altro con un bottone

---

## Come si usa

Serve Go installato (versione recente, tipo 1.18+)

Clona il repo ed esegui:

```bash
git clone https://github.com/CodingDiego150104/BlogEngine.git
cd BlogEngine
go run main.go

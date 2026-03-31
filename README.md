# Micro-Vulnerability Scanner Simulation

Roadmap step-by-step per costruire il progetto in modo ordinato (senza saltare fasi).

## Stato avanzamento (aggiornato)

Check basati sui file presenti nel workspace al momento.

- [x] Modulo Go inizializzato in `go.mod` (`module micro-vuln-scanner`).
- [x] Struttura cartelle base creata: `internal/{config,domain,httpapi,logger,service,simulator,store}`.
- [x] Cartelle `pkg` e `test` create.
- [ ] Presenza di entrypoint root con `package main` in `main.go`.
- [x] Entrypoint applicativo in `cmd/api/main.go` da allineare a `package main` + funzione `main`.
- [x] Primo avvio applicazione riuscito (`go run ./cmd/api`).
- [x] Domain model `Vulnerability` implementato (con Severity, Status, validazioni, JSON tag).
- [x] Store concorrente in-memory implementato (RWMutex, Add, GetAll, GetBySeverity, retention).
- [x] Endpoint `GET /api/vulnerabilities` implementato.
- [x] Simulazione real-time con ticker implementata (generazione attiva con graceful shutdown).

## 0) Stato iniziale e setup

Obiettivo: sistemare la base del progetto Go e partire pulito.

1. Verifica versione Go locale (`go version`).
2. Conferma modulo in `go.mod` (ok: `micro-vuln-scanner`).
3. Crea una branch dedicata (consigliato) per lavorare in sicurezza.
4. Correggi il file entrypoint: il file eseguibile deve avere package `main`.
5. Decidi struttura definitiva:

- Opzione A: tenere `api/main.go`.
- Opzione B (piu standard): usare `cmd/api/main.go`.

6. Esegui un primo avvio minimale (`go run ./api` oppure `go run ./cmd/api`).

Done quando: il comando di run parte senza errori di package.

## 1) Struttura cartelle backend

Obiettivo: separare chiaramente responsabilita.

1. Crea/usa queste cartelle:

- `internal/domain`
- `internal/store`
- `internal/service`
- `internal/simulator`
- `internal/httpapi`
- `internal/config`
- `internal/logger`
- `pkg/response` (utility riusabili)

2. Definisci naming coerente dei file (uno scopo chiaro per file).
3. Tieni `main.go` solo come bootstrap (niente logica business dentro).

Done quando: la struttura e pronta e ogni package ha un ruolo preciso.

## 2) Domain model e validazioni

Obiettivo: modellare bene i dati prima della logica.

1. Definisci `Vulnerability` con:

- `ID`
- `ImageName`
- `Severity` (Critical, High, Medium, Low)
- `Status` (Fixed, Unfixed)
- `CreatedAt`

2. Definisci severity/status come valori chiusi (enum-like), non stringhe libere sparse.
3. Crea funzioni di parse/validazione severity e status.
4. Standardizza data in UTC e serializzazione coerente JSON.

Done quando: input invalidi vengono rifiutati subito e in modo prevedibile.

## 3) Store in-memory thread-safe

Obiettivo: evitare race condition tra API e simulatore.

1. Implementa store in memoria con lock (`sync.RWMutex`).
2. Metodi minimi:

- add vulnerability
- list all vulnerabilities
- list by severity

3. In lettura restituisci copie dei dati, non riferimenti mutabili interni.
4. (Opzionale ma consigliato) retention: mantieni solo ultime N vulnerabilities.

Done quando: lo store e sicuro in concorrenza.

## 4) Service layer

Obiettivo: centralizzare business logic.

1. Crea un service che usa lo store.
2. Sposta qui logica di filtro, ordinamento e regole applicative.
3. Definisci errori noti (es. severity non valida) in modo chiaro.
4. Mantieni handlers HTTP sottili: niente logica pesante negli endpoint.

Done quando: il service e il punto unico di logica applicativa.

## 5) Simulazione real-time

Obiettivo: generare vulnerabilita automaticamente ogni 5 secondi.

1. Crea una goroutine con ticker configurabile (default 5s).
2. A ogni tick genera una vulnerability casuale ma realistica.
3. Inseriscila nello store tramite service.
4. Gestisci stop pulito con context cancellation e `ticker.Stop()`.

Done quando: i dati crescono nel tempo senza memory leak di goroutine.

## 6) API HTTP

Obiettivo: esporre endpoint pulito e filtrabile.

1. Implementa `GET /api/vulnerabilities`.
2. Supporta query param `severity` (case-insensitive).
3. Se severity non valida: risposta `400` con errore JSON coerente.
4. Se ok: risposta `200` con lista (anche vuota).
5. Aggiungi header corretti (`Content-Type: application/json`).

Done quando: endpoint stabile e contratto API chiaro.

## 7) Bonus backend: context timeout

Obiettivo: mostrare maturita architetturale.

1. Nel layer HTTP crea un context con timeout per la richiesta.
2. Passa il context al service/store.
3. Gestisci timeout/cancel con risposta errore comprensibile.
4. Rendi timeout configurabile via env.

Done quando: la request non resta bloccata indefinitamente.

## 8) Config e logging

Obiettivo: servizio configurabile e osservabile.

1. Configura via env vars:

- porta
- intervallo ticker
- timeout API
- max elementi retention (se usata)

2. Aggiungi logging strutturato per:

- startup/shutdown
- nuova vulnerability generata
- richieste API (path, filtro, status, latenza)
- errori

Done quando: puoi capire facilmente cosa succede in runtime.

## 9) Test backend

Obiettivo: dimostrare qualita tecnica in colloquio.

1. Unit test su parse/validazione severity-status.
2. Unit test su filtri del service.
3. Test store concorrente.
4. Test endpoint:

- senza filtro
- filtro valido
- filtro invalido

5. Esegui race detector: `go test -race ./...`.

Done quando: test verdi + nessuna race rilevata.

## 10) Frontend React + TypeScript

Obiettivo: dashboard professionale e aggiornata automaticamente.

1. Crea layout:

- Sidebar filtri severity
- Main area tabella/card

2. Crea custom hook `useVulnerabilities` per:

- fetch dati
- stato loading/error
- filtro selezionato
- polling ogni 5-10s

3. Aggiungi badge/colori severity (critical rosso, high arancio, medium giallo, low verde/azzurro).
4. Gestisci stati UI: loading, empty, error.

Done quando: UI reattiva e leggibile, senza reload pagina.

## 11) Integrazione FE-BE

Obiettivo: collegamento robusto tra servizi.

1. Configura URL API via env frontend.
2. Gestisci CORS lato backend (solo per sviluppo locale).
3. Verifica mapping campi JSON e formati data.
4. Testa filtro sidebar -> query `?severity=...`.

Done quando: filtro frontend riflette correttamente i dati backend.

## 12) Docker / Docker Compose

Obiettivo: avvio rapido in stile aziendale.

1. Crea Dockerfile backend.
2. Crea Dockerfile frontend.
3. Crea `docker-compose.yml` con 2 servizi:

- backend
- frontend

4. Esponi porte e variabili ambiente necessarie.
5. Avvia e verifica l'intero stack da compose.

Done quando: tutto parte con un solo comando.

## 13) Checklist finale pre-colloquio

1. Il backend genera una vulnerability ogni 5s.
2. L'endpoint API filtra severity correttamente.
3. La dashboard aggiorna dati con polling.
4. Colori e stati UI sono chiari.
5. Test principali passano, incluso race detector.
6. Avvio completo con Docker Compose.
7. Sai spiegare in 60 secondi:

- separazione layer
- concorrenza con mutex
- context timeout
- custom hook frontend

---

## Ordine consigliato di lavoro (rapido)

1. Step 0-1 (setup + struttura)
2. Step 2-6 (backend core)
3. Step 7-9 (bonus + test)
4. Step 10-11 (frontend + integrazione)
5. Step 12-13 (docker + rifinitura colloquio)

Se segui questo ordine, arrivi a una demo solida in tempi brevi e con una narrazione tecnica convincente.

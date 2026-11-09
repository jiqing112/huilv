package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fx/internal"
)

//go:embed all:frontend/dist
var distFS embed.FS

//go:embed internal/seeddata/rates.json
var seedJSON []byte

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dataDir := flag.String("data", "data", "sqlite data directory")
	flag.Parse()

	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		log.Fatal(err)
	}
	store, err := internal.Open(filepath.Join(*dataDir, "fx.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// First boot offline still works: the embedded seed carries ~17 months
	// of real ECB daily rates.
	if store.RateCount() == 0 {
		if err := loadSeed(store); err != nil {
			log.Fatal(err)
		}
		log.Printf("seeded %d rates", store.RateCount())
	}
	if err := store.Prune(); err != nil {
		log.Printf("prune: %v", err)
	}
	refresh(store)
	go func() {
		for range time.Tick(time.Hour) {
			refresh(store)
		}
	}()

	api := (&internal.Server{Store: store}).Routes()
	sub, err := fs.Sub(distFS, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/", spa(sub))

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func loadSeed(store *internal.Store) error {
	var seed struct {
		Rates map[string]map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(seedJSON, &seed); err != nil {
		return err
	}
	for date, rates := range seed.Rates {
		if err := store.UpsertDay(date, rates); err != nil {
			return err
		}
	}
	return store.MetaSet("seeded_at", time.Now().Format(time.RFC3339))
}

// refresh pulls the latest day plus a short backfill (ECB occasionally
// revises recent days). Failures are non-fatal: stale data still serves.
func refresh(store *internal.Store) {
	defer func() {
		store.MetaSet("last_refresh", time.Now().Format(time.RFC3339))
	}()
	start := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	days, err := internal.FetchBackfill(start, end)
	if err != nil {
		log.Printf("refresh backfill: %v", err)
		return
	}
	for date, rates := range days {
		if err := store.UpsertDay(date, rates); err != nil {
			log.Printf("refresh upsert %s: %v", date, err)
		}
	}
}

// spa serves the embedded frontend: hashed assets get immutable caching,
// everything else falls back to index.html.
func spa(fsys fs.FS) http.Handler {
	fileServer := http.FileServerFS(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}
		f, err := fsys.Open(p)
		if err != nil {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()
		if strings.HasPrefix(p, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fileServer.ServeHTTP(w, r)
	})
}

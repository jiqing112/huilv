package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	Store *Store
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *Server) handleCurrencies(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Currencies)
}

// handleLatest: /api/latest?base=USD -> {date, base, rates}
func (s *Server) handleLatest(w http.ResponseWriter, r *http.Request) {
	date, day := s.Store.Latest()
	if day == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no data"})
		return
	}
	base := r.URL.Query().Get("base")
	if base == "" {
		base = "USD"
	}
	baseRate, ok := day[base]
	if !ok || baseRate == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown base " + base})
		return
	}
	rates := make(map[string]float64, len(day))
	for code, rate := range day {
		rates[code] = rate / baseRate
	}
	writeJSON(w, http.StatusOK, map[string]any{"date": date, "base": base, "rates": rates})
}

// handleHistory: /api/history?from=USD&to=CNY&days=365
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from, to := q.Get("from"), q.Get("to")
	days, _ := strconv.Atoi(q.Get("days"))
	if days <= 0 {
		days = 365
	}
	if days > 1825 {
		days = 1825
	}
	if CurrencyByCode[from].Code == "" || CurrencyByCode[to].Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown currency"})
		return
	}
	type point struct {
		D string  `json:"d"`
		R float64 `json:"r"`
	}
	var points []point
	for _, row := range s.Store.Series(days) {
		if rate, ok := Cross(row.Rates, from, to); ok {
			points = append(points, point{row.Date, rate})
		}
	}
	if points == nil {
		points = []point{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"from": from, "to": to, "days": days, "points": points})
}

// handleStatus: freshness info for the UI data dot and ops checks.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	date, _ := s.Store.Latest()
	writeJSON(w, http.StatusOK, map[string]any{
		"date":    date,
		"updated": s.Store.MetaGet("last_refresh"),
		"count":   s.Store.RateCount(),
		"now":     time.Now().Format(time.RFC3339),
	})
}

// handleBoard: /api/board?base=USD&days=30 -> relative series for every
// currency at once (one request feeds the whole rates list).
func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	base := q.Get("base")
	days, _ := strconv.Atoi(q.Get("days"))
	if days <= 0 {
		days = 30
	}
	if days > 1825 {
		days = 1825
	}
	if base == "" {
		base = "USD"
	}
	if CurrencyByCode[base].Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown base"})
		return
	}
	type point struct {
		D string  `json:"d"`
		R float64 `json:"r"`
	}
	series := map[string][]point{}
	for _, row := range s.Store.Series(days) {
		b, ok := row.Rates[base]
		if !ok || b == 0 {
			continue
		}
		for code, rate := range row.Rates {
			if code == base {
				continue
			}
			series[code] = append(series[code], point{row.Date, rate / b})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"base": base, "days": days, "series": series})
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/currencies", s.handleCurrencies)
	mux.HandleFunc("GET /api/latest", s.handleLatest)
	mux.HandleFunc("GET /api/history", s.handleHistory)
	mux.HandleFunc("GET /api/board", s.handleBoard)
	mux.HandleFunc("GET /api/status", s.handleStatus)
	return mux
}

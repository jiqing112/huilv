package internal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Store keeps one table: USD-based daily rates. Cross rates are derived
// (from -> to = usd[to] / usd[from]) so the whole dataset stays in one
// small table with full history.
type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	for _, stmt := range []string{
		`CREATE TABLE IF NOT EXISTS rates (date TEXT NOT NULL, code TEXT NOT NULL, rate REAL NOT NULL, PRIMARY KEY (date, code))`,
		`CREATE TABLE IF NOT EXISTS meta  (k TEXT PRIMARY KEY, v TEXT NOT NULL)`,
	} {
		if _, err := db.Exec(stmt); err != nil {
			return nil, fmt.Errorf("init schema: %w", err)
		}
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Prune drops rates for codes that are no longer in the metadata set
// (e.g. BGN, which the ECB stopped publishing when Bulgaria adopted the euro).
func (s *Store) Prune() error {
	codes := make([]string, 0, len(CurrencyByCode))
	for code := range CurrencyByCode {
		codes = append(codes, code)
	}
	in := "'" + strings.Join(codes, "','") + "'"
	_, err := s.db.Exec(`DELETE FROM rates WHERE code NOT IN (` + in + `)`)
	return err
}

func (s *Store) MetaGet(k string) string {
	var v string
	s.db.QueryRow(`SELECT v FROM meta WHERE k = ?`, k).Scan(&v)
	return v
}

func (s *Store) MetaSet(k, v string) error {
	_, err := s.db.Exec(`INSERT INTO meta (k, v) VALUES (?, ?) ON CONFLICT (k) DO UPDATE SET v = excluded.v`, k, v)
	return err
}

func (s *Store) RateCount() int {
	var n int
	s.db.QueryRow(`SELECT COUNT(*) FROM rates`).Scan(&n)
	return n
}

// UpsertDay writes one day of USD-based rates.
func (s *Store) UpsertDay(date string, rates map[string]float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO rates (date, code, rate) VALUES (?, ?, ?) ON CONFLICT (date, code) DO UPDATE SET rate = excluded.rate`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for code, rate := range rates {
		if _, err := stmt.Exec(date, code, rate); err != nil {
			return err
		}
	}
	// The source omits the base currency itself; keep USD present in every
	// day so every cross rate (and USD queries) always resolve.
	if _, ok := rates["USD"]; !ok && len(rates) > 0 {
		if _, err := stmt.Exec(date, "USD", 1.0); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Latest returns the most recent date and its USD-based rates.
func (s *Store) Latest() (string, map[string]float64) {
	var date string
	if err := s.db.QueryRow(`SELECT MAX(date) FROM rates`).Scan(&date); err != nil || date == "" {
		return "", nil
	}
	return date, s.Day(date)
}

// Day returns USD-based rates for one date.
func (s *Store) Day(date string) map[string]float64 {
	rows, err := s.db.Query(`SELECT code, rate FROM rates WHERE date = ?`, date)
	if err != nil {
		return nil
	}
	defer rows.Close()
	m := make(map[string]float64, 32)
	for rows.Next() {
		var code string
		var rate float64
		if rows.Scan(&code, &rate) == nil {
			m[code] = rate
		}
	}
	return m
}

// Cross returns the from -> to rate on a date, or false if data is missing.
func Cross(day map[string]float64, from, to string) (float64, bool) {
	f, ok1 := day[from]
	t, ok2 := day[to]
	if !ok1 || !ok2 || f == 0 {
		return 0, false
	}
	return t / f, true
}

// Series returns USD-based daily rates over the trailing n calendar days.
func (s *Store) Series(days int) []struct {
	Date  string
	Rates map[string]float64
} {
	if days < 1 {
		days = 1
	}
	rows, err := s.db.Query(`SELECT date, code, rate FROM rates WHERE date >= ? ORDER BY date`, time.Now().AddDate(0, 0, -days).Format("2006-01-02"))
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []struct {
		Date  string
		Rates map[string]float64
	}
	var cur string
	var m map[string]float64
	flush := func() {
		if cur != "" && len(m) > 0 {
			out = append(out, struct {
				Date  string
				Rates map[string]float64
			}{cur, m})
		}
	}
	for rows.Next() {
		var date, code string
		var rate float64
		if rows.Scan(&date, &code, &rate) != nil {
			continue
		}
		if date != cur {
			flush()
			cur, m = date, make(map[string]float64, 32)
		}
		m[code] = rate
	}
	flush()
	return out
}

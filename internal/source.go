package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Frankfurter serves the ECB reference rates without an API key.
const sourceBase = "https://api.frankfurter.dev/v1"

type frankfurterResp struct {
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type frankfurterSeries struct {
	Start string                       `json:"start_date"`
	End   string                       `json:"end_date"`
	Rates map[string]map[string]float64 `json:"rates"`
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

func fetchJSON(url string, into any) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("source %s: HTTP %d", url, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(into)
}

// FetchLatest returns the newest USD-based rate day from the source.
func FetchLatest() (string, map[string]float64, error) {
	var r frankfurterResp
	if err := fetchJSON(sourceBase+"/latest?base=USD", &r); err != nil {
		return "", nil, err
	}
	return r.Date, r.Rates, nil
}

// FetchBackfill returns USD-based rates for [start, end] inclusive.
func FetchBackfill(start, end string) (map[string]map[string]float64, error) {
	var r frankfurterSeries
	if err := fetchJSON(fmt.Sprintf("%s/%s..%s?base=USD", sourceBase, start, end), &r); err != nil {
		return nil, err
	}
	return r.Rates, nil
}

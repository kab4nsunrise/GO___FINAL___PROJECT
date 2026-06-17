package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}
	parts := strings.Fields(repeat)
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("bad d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid days")
		}
		next := start
		for !next.After(now) {
			next = next.AddDate(0, 0, days)
		}
		return next.Format(dateFormat), nil

	case "y":
		next := start
		for !next.After(now) {
			next = next.AddDate(1, 0, 0)
		}
		return next.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("unsupported rule: %s", rule)
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now", http.StatusBadRequest)
			return
		}
	}
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}

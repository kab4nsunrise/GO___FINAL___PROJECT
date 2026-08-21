package api

import (
	"fmt"
	"log"
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
		for {
			next = next.AddDate(0, 0, days)
			log.Printf("d: next=%s, now=%s", next.Format(dateFormat), now.Format(dateFormat))
			if next.After(now) {
				break
			}
		}
		result := next.Format(dateFormat)
		log.Printf("d returning: %s", result)
		return result, nil

	case "y":
		next := start
		for {
			next = next.AddDate(1, 0, 0)
			log.Printf("y: next=%s, now=%s", next.Format(dateFormat), now.Format(dateFormat))
			if next.Month() == time.February && next.Day() == 29 && next.Year()%4 != 0 {
				next = time.Date(next.Year(), time.March, 1, 0, 0, 0, 0, next.Location())
				log.Printf("y corrected to: %s", next.Format(dateFormat))
			}
			if next.After(now) {
				break
			}
		}
		result := next.Format(dateFormat)
		log.Printf("y returning: %s", result)
		return result, nil

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
	if _, err := w.Write([]byte(next)); err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
	}
}

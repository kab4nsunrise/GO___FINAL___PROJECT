package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		errorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorResponse(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		errorResponse(w, "Title required", http.StatusBadRequest)
		return
	}

	if err := normalizeTaskDate(&task); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		errorResponse(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Task added: id=%d, title=%s", id, task.Title)
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse(w, "id required", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorResponse(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		errorResponse(w, "id required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		errorResponse(w, "Title required", http.StatusBadRequest)
		return
	}
	if err := normalizeTaskDate(&task); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		errorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse(w, "id required", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		errorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func normalizeTaskDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	_, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}
	today := now.Format(dateFormat)

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("repeat rule error: %v", err)
		}

		if task.Date < today {
			task.Date = next
		}
	} else {
		if task.Date < today {
			task.Date = today
		}
	}
	return nil
}

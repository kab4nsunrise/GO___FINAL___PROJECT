package api

import (
	"net/http"
	"todo/pkg/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	tasks, err := db.Tasks(limit)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJSON(w, map[string]interface{}{"tasks": tasks}, http.StatusOK)
}

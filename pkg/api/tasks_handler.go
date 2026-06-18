package api

import (
	"net/http"
	"todo/pkg/db"
)

const defaultTasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(defaultTasksLimit)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJSON(w, map[string]interface{}{"tasks": tasks}, http.StatusOK)
}

package api

import (
    "net/http"
    "time"
    "todo/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
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
    if task.Repeat == "" {
        err = db.DeleteTask(id)
    } else {
        next, err := NextDate(time.Now(), task.Date, task.Repeat)
        if err != nil {
            errorResponse(w, err.Error(), http.StatusBadRequest)
            return
        }
        err = db.UpdateDate(id, next)
    }
    if err != nil {
        errorResponse(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

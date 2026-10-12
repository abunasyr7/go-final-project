package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abunasyr7/go-final-project/pkg/db"
)

func doneTaskHandler (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if len(id) == 0 {
		writeError(w, fmt.Errorf("no ID"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	if len(task.Repeat) == 0 {
		if err = db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError)
		}
		writeJSON(w, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err = db.UpdateDate(next, id); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeError(w, fmt.Errorf("no ID"), http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{})
}
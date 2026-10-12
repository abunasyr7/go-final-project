package api

import (
	"net/http"

	"github.com/abunasyr7/go-final-project/pkg/db"
)

const tasksLiimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(search, tasksLiimit)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
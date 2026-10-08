package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/abunasyr7/go-final-project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)

	if err != nil {
		return fmt.Errorf("invalid date %q: %w", task.Date, err)
	}

	var next string
	
	if len(task.Repeat) > 0 {
		next, err = NextDate(now, task.Date, task.Repeat)

		if err != nil {
			return fmt.Errorf("invalid repeat rule %q: %w", task.Repeat, err)
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}


func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)

	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		writeError(w, err, http.StatusBadRequest)
	} 

	if len(task.Title) == 0 {
		writeError(w, fmt.Errorf("no title of task"), http.StatusBadRequest)
		return
	}

	if err = checkDate(&task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)

	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{"id": id})
}
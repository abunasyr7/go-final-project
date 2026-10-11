
func getTaskHandler (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeError(w, fmt.Errorf("no ID"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err !=nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, task)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if len(task.ID) == 0 {
		writeError(w, fmt.Errorf("no ID"), http.StatusBadRequest)
		retrn
	}

	if len(task.Title) == 0 {
		writeError(w, fmt.Errorf("no title"), http.StatusBadRequest)
	}

	if err = db.UpdateTask(&task); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[strinf]any{})
}
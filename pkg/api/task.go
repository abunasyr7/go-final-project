
func getTaskHandler (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeError(w, fmt.Errorf("no ID"))
		return
	}

	task, err := db.GetTask(id)
	if err !=nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, task)
}
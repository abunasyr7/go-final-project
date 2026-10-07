package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	w.Write(resp)
}
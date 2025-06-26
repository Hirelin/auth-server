package handlers

import (
	"encoding/json"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "pong"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}
}

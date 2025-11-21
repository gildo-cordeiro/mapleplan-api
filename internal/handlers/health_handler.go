package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthCheck struct {
}

func (handler *HealthCheck) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "healthy",
		"message": "API is running smoothly",
	}

	json.NewEncoder(w).Encode(response)
}

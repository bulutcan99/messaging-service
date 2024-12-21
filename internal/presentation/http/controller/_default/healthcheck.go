package _default

import (
	"github.com/goccy/go-json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "up"}
	json.NewEncoder(w).Encode(response)
}

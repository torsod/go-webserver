package handler

import (
	"net/http"
	"time"
)

var startTime = time.Now()

// HealthHandler returns server health status
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"uptime":  time.Since(startTime).String(),
			"version": "1.0.0",
		})
	}
}

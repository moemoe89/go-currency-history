package routes

import "net/http"

// HealthCheckHandler handler for health check.
func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

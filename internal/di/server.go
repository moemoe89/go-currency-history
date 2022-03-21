package di

import "net/http"

// GetHTTPServer gets the parameter of Server for running an HTTP server.
func GetHTTPServer() http.Server {
	return http.Server{
		Addr:    ":8080",
		Handler: GetHTTPHandler(),
	}
}

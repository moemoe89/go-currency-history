package di

import (
	"net/http"

	"github.com/moemoe89/go-currency-history/internal/routes"
)

// GetHTTPHandler gets the Handler responds to an HTTP request.
func GetHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	v1 := "/v1"

	mux.HandleFunc("/", routes.HealthCheckHandler)
	mux.HandleFunc("/health", routes.HealthCheckHandler)
	mux.HandleFunc(v1+"/currency", routes.CurrencyHandler)
	mux.HandleFunc(v1+"/currency/history", routes.CurrencyHistoryHandler)

	return mux
}

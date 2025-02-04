package middleware

import (
	"net/http"
)

// Creates a basic middleware to control CORS policy.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, PUT, GET")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// identifies the IP of the request,
// either from the request Header or the remote address.
func GetIP(r *http.Request) string {
	return ""
}

// Logger is a middleware that should be injected directly into the chi router
// to capture and log all http events automatically.
// It will automatically create 2 log entries for each HTTP call:
// 1. Step = INPUT => is the start of all processing. A transactionId is generated for this call;
// 2. Step = OUTPUT => is the end of processing;
func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}

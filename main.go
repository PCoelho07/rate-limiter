package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"ratelimiter/ratelimiter"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ratelimiter.HTTPResponse{Message: "this is an /api response"})
	})

	s := http.Server{
		Addr:    ":8080",
		Handler: logging(ratelimiter.RateLimit(mux)),
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("cannot start the server: %v", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

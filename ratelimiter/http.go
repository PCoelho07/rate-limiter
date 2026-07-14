package ratelimiter

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"sync"
)

const (
	Count       = 5
	Capacity    = 5
	RefillValue = 1
	RefillRate  = 5
)

type HTTPResponse struct {
	Message any `json:"message"`
}

func RateLimit(next http.Handler) http.Handler {
	var mu sync.Mutex
	ips := make(map[string]*Bucket)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		mu.Lock()
		bucket, exists := ips[ip]
		if !exists {
			bucket = NewBucket(Capacity, Count, RefillValue, RefillRate)
			ips[ip] = bucket
		}

		bucket.Refill()
		err = bucket.UseToken()
		mu.Unlock()

		if err != nil {
			slog.Error("too many requests.", "remoteAddr", ip)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(HTTPResponse{Message: "too many requests"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

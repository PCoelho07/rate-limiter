# rate-limiter

A per-IP HTTP rate limiter for Go, built with the token bucket algorithm.
It ships as standard `net/http` middleware that you wrap around any handler.

## How it works

The rate limiter keeps one token bucket per client IP address. Each bucket
starts with a fixed number of tokens. Every request consumes one token:

- If the bucket has a token, the request passes through to the next handler.
- If the bucket is empty, the limiter responds with
  `429 Too Many Requests` and a JSON body.

Buckets refill over time. When at least one second has passed since the last
refill, the bucket gains `RefillRate` tokens. Client IPs are tracked in an
in-memory map guarded by a mutex, so the limiter is safe for concurrent use.

## Requirements

- Go 1.23.3 or later

## Run the server

1. Clone the repository and change into the project directory.
2. Start the server:

   ```sh
   go run .
   ```

3. Send a request to the example endpoint:

   ```sh
   curl -i http://localhost:8080/api
   ```

The server listens on port `8080` and exposes a single `/api` endpoint that
returns a JSON message. Every request is logged with its method and path.

## Try the rate limit

The default configuration allows a small burst of requests before the bucket
empties. Send several requests in quick succession to trigger the limit:

```sh
for i in $(seq 1 10); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/api; done
```

After the initial tokens are used, you see `429` responses until the bucket
refills.

## Use the middleware

Wrap any `http.Handler` with `ratelimiter.RateLimit` to apply the limit:

```go
package main

import (
	"net/http"

	"ratelimiter/ratelimiter"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: ratelimiter.RateLimit(mux),
	}

	server.ListenAndServe()
}
```

## Configuration

The limiter reads its settings from constants in
`ratelimiter/http.go`. Adjust these values to change the limit:

| Constant      | Default | Description                                    |
| ------------- | ------- | ---------------------------------------------- |
| `Count`       | `5`     | Number of tokens a new bucket starts with.     |
| `RefillRate`  | `5`     | Tokens added each time the bucket refills.     |
| `Capacity`    | `5`     | Stored on the bucket for future use.           |
| `RefillValue` | `1`     | Stored on the bucket for future use.           |

Refills happen once at least one second has passed since the previous refill.

> **Note:** The current refill logic does not cap `Count` at `Capacity`, so a
> long-idle bucket can accumulate more tokens than `Capacity`.

## Project structure

- `main.go` — Starts the HTTP server and adds request logging.
- `ratelimiter/http.go` — The `RateLimit` middleware and configuration.
- `ratelimiter/bucket.go` — The token bucket type and its refill logic.

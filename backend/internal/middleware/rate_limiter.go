// rate_limiter.go
package middleware

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiterMiddleware() func(http.Handler) http.Handler {
	rate, err := limiter.NewRateFromFormatted("50-M")
	if err != nil {
		log.Fatalf("Failed to create rate limiter: %v", err)
	}

	store := memory.NewStore()

	instance := limiter.New(store, rate)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("Failed to parse IP from RemoteAddr: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			limiterContext, err := instance.Get(r.Context(), ip)
			if err != nil {
				log.Printf("Error retrieving rate limit info: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(int64(limiterContext.Limit), 10))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(limiterContext.Remaining), 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(limiterContext.Reset), 10))

			if limiterContext.Reached {
				log.Printf("Rate limit exceeded for IP: %s", ip)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

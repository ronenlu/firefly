package data

import (
	"bufio"
	"context"
	_ "embed"
	"strings"

	"golang.org/x/time/rate"
)

//go:embed endg-urls.txt
var urlsFile string

// RateLimiter wraps the golang.org/x/time/rate.Limiter
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter that allows up to requestsPerSecond
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	// rate.Limit defines requests per second
	// burst allows for burst of requests up to that number
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
	}
}

// Wait blocks until a request can proceed according to the rate limit
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// LoadEssayURLs loads and parses the embedded URLs file, returning a slice of trimmed, non-empty URLs.
func LoadEssayURLs() []string {
	scanner := bufio.NewScanner(strings.NewReader(urlsFile))
	var urls []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		urls = append(urls, line)
	}
	return urls
}

package secedgar

import (
	"sync"
	"time"
)

// minRequestInterval keeps us under the 10 requests per second the SEC asks automated
// clients to respect. The limiter is package level so that concurrent Form 4 fetches and
// any other caller share one budget.
const minRequestInterval = 125 * time.Millisecond

type rateLimiter struct {
	mu          sync.Mutex
	lastRequest time.Time
}

var limiter = &rateLimiter{}

// wait blocks until enough time has passed since the previous request.
func (r *rateLimiter) wait() {
	r.mu.Lock()
	defer r.mu.Unlock()

	elapsed := time.Since(r.lastRequest)
	if elapsed < minRequestInterval {
		time.Sleep(minRequestInterval - elapsed)
	}

	r.lastRequest = time.Now()
}

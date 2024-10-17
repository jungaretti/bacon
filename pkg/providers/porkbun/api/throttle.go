package api

import (
	"time"
)

type Throttler struct {
	ticker *time.Ticker
}

func NewThrottler(ratePerSecond int) *Throttler {
	interval := time.Second / time.Duration(ratePerSecond)
	ticker := time.NewTicker(interval)

	return &Throttler{
		ticker: ticker,
	}
}

// Waits for the next available permit.
// Use this method to respect rate limits.
func (t *Throttler) WaitForPermit() {
	// Send a tick to 1 caller per interval
	<-t.ticker.C
}

func (t *Throttler) Stop() {
	t.ticker.Stop()
}

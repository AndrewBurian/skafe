package main

import (
	"testing"
	"time"
)

// Check the ratelimiter is operating in the right range
func TestRateLimit(t *testing.T) {

	rate := uint(100)
	span := time.Second

	in := make(chan AuditEvent, 50)
	out := make(chan AuditEvent, 50)

	ev := map[string]string{
		"key": "val",
	}

	start := time.Now()

	for i := 0; i < 20; i++ {
		in <- ev
	}
	close(in)

	RateLimit(in, out, rate, span)

	for _ = range out {
		// just get
	}

	runTime := time.Since(start)

	t.Logf("Took %d Nanoseconds", runTime.Nanoseconds())

	if runTime.Nanoseconds() < 200000000*int64(time.Nanosecond) {
		t.Errorf("Rate Limiter ran too fast")
	}

	if runTime.Nanoseconds() > 250000000*int64(time.Nanosecond) {
		t.Errorf("Rate Limiter ran too slow")
	}
}

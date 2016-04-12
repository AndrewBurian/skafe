package main

import (
	"time"
)

func RateLimit(in <-chan AuditEvent, out chan<- AuditEvent, rate uint, span time.Duration) {

	// calculate the number of seconds to delay between messages
	var delayTime float32 = float32(span.Seconds()) / float32(rate)

	// upcast to milis and truncate
	var numMilis int64 = int64(delayTime * 1000)

	delay := time.Millisecond * time.Duration(numMilis)

	// create and start the ticker
	ticker := time.NewTicker(delay)

	// so long as there's events
	for ev := range in {

		// forawrd an event
		out <- ev

		// wait for the rate limit tick
		<-ticker.C
	}

	close(out)
	ticker.Stop()

}

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

	for ev := range in {
		out <- ev

		time.Sleep(delay)
	}

	close(out)

}

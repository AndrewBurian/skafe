package main

import (
	"io"
	"sync"
)

// An thread-safe io.Writer
type SyncWriter struct {
	lock sync.Mutex
	w    io.Writer
}

func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{
		w: w,
	}
}

func (s *SyncWriter) Write(buf []byte) (int, error) {
	s.lock.Lock()
	n, err := s.w.Write(buf)
	s.lock.Unlock()
	return n, err
}

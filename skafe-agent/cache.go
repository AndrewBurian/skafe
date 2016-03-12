package main

import (
	"container/list"
	"database/sql"
)

type FileCache struct {
	db    *sql.DB
	count int
}

func Cache(in <-chan AuditEvent, out chan<- AuditEvent, numCached int, db *sql.DB) {

	// init the database cache
	cache := NewCache(db)

	// create the in-memory queue for events
	queue := list.New()

	// load up the queue from disk
	for i := 0; i < numCached && cache.HasNext(); i++ {
		queue.PushBack(cache.Pop())
	}

	// Main loop
	for {

		if queue.Len() > 0 {
			// Data is waiting to be sent
			select {
			case out <- queue.Remove(queue.Front()).(AuditEvent):
				// event sent
				if cache.HasNext() {
					queue.PushBack(cache.Pop())
				}

			case ev, ok := <-in:
				// event received
				if !ok {
					break
				}
				if queue.Len() < numCached {
					queue.PushBack(ev)
				} else {
					cache.PushBack(ev)
				}
			}
		} else {
			// No data to send, only receive
			ev, ok := <-in
			if !ok {
				break
			}
			queue.PushBack(ev)
		}
	}

	// inbound connection closed, shutting down
	for queue.Len() > 0 {
		cache.PushBack(queue.Remove(queue.Front()).(AuditEvent))
	}

}

func NewCache(db *sql.DB) *FileCache {

	c := &FileCache{
		db:    db,
		count: 0,
	}

	c.InitCount()

	return c
}

func (c *FileCache) InitCount() {
	//TODO
}

func (c *FileCache) PushBack(ev AuditEvent) {
	//TODO
}

func (c *FileCache) Pop() AuditEvent {
	//TODO
	return nil
}

func (c *FileCache) HasNext() bool {
	return c.count > 0
}

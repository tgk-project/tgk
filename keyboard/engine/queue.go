// Package engine owns the Primary keyboard state and its bounded event queue.
package engine

import (
	"errors"
	"sync"

	"github.com/tgk-project/tgk/keyboard/event"
)

var ErrInvalidQueueCapacity = errors.New("event queue capacity must be positive")

// Queue is a fixed-capacity ring buffer. Producers only enqueue events; the
// engine event loop is the sole consumer and owner of keyboard state.
type Queue struct {
	mu            sync.Mutex
	events        []event.KeyEvent
	head          int
	tail          int
	count         int
	overflowed    bool
	overflowCount uint32
}

// NewQueue preallocates a bounded queue and never grows it on the event path.
func NewQueue(capacity int) (*Queue, error) {
	if capacity <= 0 {
		return nil, ErrInvalidQueueCapacity
	}
	return &Queue{events: make([]event.KeyEvent, capacity)}, nil
}

// Push adds an event without invoking the keymap engine. It returns false when
// the event is rejected because the queue is full.
func (q *Queue) Push(value event.KeyEvent) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == len(q.events) {
		q.overflowed = true
		q.overflowCount++
		return false
	}
	q.events[q.tail] = value
	q.tail = (q.tail + 1) % len(q.events)
	q.count++
	return true
}

func (q *Queue) pop() (event.KeyEvent, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		return event.KeyEvent{}, false
	}
	value := q.events[q.head]
	q.head = (q.head + 1) % len(q.events)
	q.count--
	return value, true
}

// recoverOverflow drops queued input that may no longer have matching release
// events. The owner loop must release state before processing future events.
func (q *Queue) recoverOverflow() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.overflowed {
		return false
	}
	q.head = 0
	q.tail = 0
	q.count = 0
	q.overflowed = false
	return true
}

// QueueStats exposes loss detection without binding the core to a logger.
type QueueStats struct {
	OverflowCount uint32
}

// Stats returns a consistent queue snapshot.
func (q *Queue) Stats() QueueStats {
	q.mu.Lock()
	defer q.mu.Unlock()
	return QueueStats{OverflowCount: q.overflowCount}
}

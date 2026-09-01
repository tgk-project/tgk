package usb

import "sync"

type rawReport [RawReportSize]byte

// rawQueue is a fixed-capacity callback-to-event-loop hand-off. USB callback
// work stays bounded and allocation-free after Device construction.
type rawQueue struct {
	mu       sync.Mutex
	reports  []rawReport
	head     int
	tail     int
	count    int
	overflow uint32
}

func (q *rawQueue) push(value []byte) bool {
	if len(value) != RawReportSize {
		return false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == len(q.reports) {
		q.overflow++
		return false
	}
	copy(q.reports[q.tail][:], value)
	q.tail = (q.tail + 1) % len(q.reports)
	q.count++
	return true
}

func (q *rawQueue) pop(dst []byte) (int, bool) {
	if len(dst) < RawReportSize {
		return 0, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		return 0, false
	}
	copy(dst[:RawReportSize], q.reports[q.head][:])
	q.reports[q.head] = rawReport{}
	q.head = (q.head + 1) % len(q.reports)
	q.count--
	return RawReportSize, true
}

func (q *rawQueue) overflowCount() uint32 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.overflow
}

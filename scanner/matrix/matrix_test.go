package matrix_test

import (
	"errors"
	"testing"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/scanner/matrix"
)

type outputPin struct{ level bool }

func (p *outputPin) Set(high bool) { p.level = high }

type inputPin struct{ read func() bool }

func (p inputPin) Get() bool { return p.read() }

func TestMatrixDebounceAndPositionMapping(t *testing.T) {
	row := &outputPin{}
	pressed := true
	column := inputPin{read: func() bool {
		if !row.level && pressed { // ActiveLow row with a pull-up input.
			return false
		}
		return true
	}}
	scanner, err := matrix.New(matrix.Config{
		Rows:      []matrix.OutputPin{row},
		Columns:   []matrix.InputPin{column},
		Positions: []uint16{42},
		Source:    3,
		ActiveLow: true,
		Debounce:  2,
	})
	if err != nil {
		t.Fatal(err)
	}

	buffer := make([]event.KeyEvent, 0, 1)
	if events, err := scanner.Scan(10, buffer); err != nil || len(events) != 0 {
		t.Fatalf("first Scan() = %#v, %v; want no event", events, err)
	}
	events, err := scanner.Scan(12, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantDown := event.KeyEvent{Source: 3, Position: 42, Pressed: true, Timestamp: 12}
	if len(events) != 1 || events[0] != wantDown {
		t.Fatalf("press event = %#v, want %#v", events, wantDown)
	}

	pressed = false
	if events, err := scanner.Scan(13, buffer); err != nil || len(events) != 0 {
		t.Fatalf("release debounce start = %#v, %v; want no event", events, err)
	}
	events, err = scanner.Scan(15, buffer)
	if err != nil {
		t.Fatal(err)
	}
	wantUp := event.KeyEvent{Source: 3, Position: 42, Pressed: false, Timestamp: 15}
	if len(events) != 1 || events[0] != wantUp {
		t.Fatalf("release event = %#v, want %#v", events, wantUp)
	}
	if !row.level {
		t.Fatal("row was left active after Scan")
	}
}

func TestMatrixRetainsTransitionWhenBufferIsFull(t *testing.T) {
	row := &outputPin{}
	column := inputPin{read: func() bool { return row.level }}
	scanner, err := matrix.New(matrix.Config{
		Rows:      []matrix.OutputPin{row},
		Columns:   []matrix.InputPin{column},
		ActiveLow: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := scanner.Scan(1, nil); !errors.Is(err, matrix.ErrEventBufferFull) {
		t.Fatalf("Scan(nil) error = %v, want ErrEventBufferFull", err)
	}
	events, err := scanner.Scan(2, make([]event.KeyEvent, 0, 1))
	if err != nil || len(events) != 1 || !events[0].Pressed {
		t.Fatalf("Scan(after full buffer) = %#v, %v; want retained press", events, err)
	}
}

func TestMatrixScanEventPathDoesNotAllocate(t *testing.T) {
	row := &outputPin{}
	pressed := false
	column := inputPin{read: func() bool {
		if !row.level && pressed {
			return false
		}
		return true
	}}
	scanner, err := matrix.New(matrix.Config{
		Rows:      []matrix.OutputPin{row},
		Columns:   []matrix.InputPin{column},
		ActiveLow: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	var backing [1]event.KeyEvent
	buffer := backing[:0]
	if _, err := scanner.Scan(0, buffer); err != nil {
		t.Fatal(err)
	}
	now := uint32(1)
	allocations := testing.AllocsPerRun(1_000, func() {
		pressed = true
		if _, err := scanner.Scan(now, buffer); err != nil {
			t.Fatal(err)
		}
		now++
		pressed = false
		if _, err := scanner.Scan(now, buffer); err != nil {
			t.Fatal(err)
		}
		now++
	})
	if allocations != 0 {
		t.Fatalf("Scan() allocations = %f, want 0", allocations)
	}
}

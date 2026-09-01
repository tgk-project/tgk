// Package split implements the versioned fixed-size position-event wire
// protocol shared by a Secondary transport adapter and the Primary event loop.
package split

import (
	"encoding/binary"
	"errors"

	"github.com/tgk-project/tgk/keyboard/event"
)

const (
	// Version is incremented for any incompatible wire-format change. A Primary
	// rejects unknown versions rather than guessing how to interpret input.
	Version uint8 = 1

	// PacketSize is fixed so a transport can use bounded buffers and fits the
	// default 20-byte BLE ATT notification payload. Larger snapshots use a
	// sequence of packets ending with SnapshotLast.
	PacketSize = 18
	MaxEvents  = 4
)

const (
	flagSnapshot = 1 << iota
	flagSnapshotLast
)

var (
	ErrPacketSize     = errors.New("split packet has invalid size")
	ErrVersion        = errors.New("split packet has unsupported version")
	ErrFlags          = errors.New("split packet has unknown flags")
	ErrEventCount     = errors.New("split packet has invalid event count")
	ErrSnapshotEvent  = errors.New("split snapshot contains a release event")
	ErrPressedValue   = errors.New("split event has invalid pressed value")
	ErrTrailingData   = errors.New("split packet has non-zero unused data")
	ErrSourceMismatch = errors.New("split packet source does not match session")
	ErrSequenceGap    = errors.New("split packet sequence gap")
	ErrResyncRequired = errors.New("split session requires a snapshot")
)

// PositionEvent is the compact representation carried on the wire. It is
// deliberately a physical position, never a keycode or keymap state.
type PositionEvent struct {
	Position uint16
	Pressed  bool
}

// Packet is the decoded fixed-size wire representation. Snapshot packets list
// every currently pressed position; normal packets list press/release edges.
type Packet struct {
	Source   event.Source
	Sequence uint16
	Snapshot bool
	// SnapshotLast ends a snapshot sequence. It is valid only when Snapshot is
	// set, and lets a held-key snapshot exceed MaxEvents without variable-size
	// packets or an MTU assumption.
	SnapshotLast bool
	Count        uint8
	Events       [MaxEvents]PositionEvent
}

// Encode serializes one packet. The version-one layout is:
// version, flags, source, count, sequence (little endian), then four events of
// position (little endian) plus pressed byte. Unused event bytes are zero.
func Encode(value Packet) ([PacketSize]byte, error) {
	if err := validate(value); err != nil {
		return [PacketSize]byte{}, err
	}

	var raw [PacketSize]byte
	raw[0] = Version
	if value.Snapshot {
		raw[1] = flagSnapshot
		if value.SnapshotLast {
			raw[1] |= flagSnapshotLast
		}
	}
	raw[2] = byte(value.Source)
	raw[3] = value.Count
	binary.LittleEndian.PutUint16(raw[4:6], value.Sequence)
	for index := 0; index < int(value.Count); index++ {
		offset := 6 + index*3
		binary.LittleEndian.PutUint16(raw[offset:offset+2], value.Events[index].Position)
		if value.Events[index].Pressed {
			raw[offset+2] = 1
		}
	}
	return raw, nil
}

// Decode validates and decodes exactly one version-one packet.
func Decode(raw []byte) (Packet, error) {
	if len(raw) != PacketSize {
		return Packet{}, ErrPacketSize
	}
	if raw[0] != Version {
		return Packet{}, ErrVersion
	}
	if raw[1]&^byte(flagSnapshot|flagSnapshotLast) != 0 {
		return Packet{}, ErrFlags
	}
	value := Packet{
		Source:       event.Source(raw[2]),
		Sequence:     binary.LittleEndian.Uint16(raw[4:6]),
		Snapshot:     raw[1]&flagSnapshot != 0,
		SnapshotLast: raw[1]&flagSnapshotLast != 0,
		Count:        raw[3],
	}
	if value.Count > MaxEvents || (!value.Snapshot && value.Count == 0) {
		return Packet{}, ErrEventCount
	}
	for index := 0; index < MaxEvents; index++ {
		offset := 6 + index*3
		if index >= int(value.Count) {
			if raw[offset] != 0 || raw[offset+1] != 0 || raw[offset+2] != 0 {
				return Packet{}, ErrTrailingData
			}
			continue
		}
		if raw[offset+2] > 1 {
			return Packet{}, ErrPressedValue
		}
		value.Events[index] = PositionEvent{
			Position: binary.LittleEndian.Uint16(raw[offset : offset+2]),
			Pressed:  raw[offset+2] == 1,
		}
	}
	if err := validate(value); err != nil {
		return Packet{}, err
	}
	return value, nil
}

func validate(value Packet) error {
	if value.Count > MaxEvents || (!value.Snapshot && value.Count == 0) {
		return ErrEventCount
	}
	if value.SnapshotLast && !value.Snapshot {
		return ErrFlags
	}
	if value.Snapshot && value.Count == 0 && !value.SnapshotLast {
		return ErrEventCount
	}
	if value.Snapshot {
		for index := 0; index < int(value.Count); index++ {
			if !value.Events[index].Pressed {
				return ErrSnapshotEvent
			}
		}
	}
	return nil
}

// Result tells the Primary event loop what must happen after processing a
// packet. ReleaseSource must happen before Events are enqueued for a snapshot.
type Result struct {
	Events          [MaxEvents]event.KeyEvent
	Count           int
	ReleaseSource   bool
	RequestSnapshot bool
}

// Session keeps connection-local sequence and resynchronization state. It has
// no BLE dependency and must be owned by the Primary event loop.
type Session struct {
	source       event.Source
	nextSequence uint16
	synchronized bool
	snapshotting bool
}

// NewSession starts disconnected: a snapshot is required before normal edge
// packets can affect Primary-owned keyboard state.
func NewSession(source event.Source) *Session {
	return &Session{source: source}
}

// Process accepts one decoded packet. A sequence gap or a normal packet before
// a snapshot releases the source and requires transport-level snapshot request.
func (s *Session) Process(packet Packet) (Result, error) {
	if packet.Source != s.source {
		return Result{}, ErrSourceMismatch
	}
	if packet.Snapshot {
		if s.snapshotting && packet.Sequence != s.nextSequence {
			s.snapshotting = false
			s.synchronized = false
			return Result{ReleaseSource: true, RequestSnapshot: true}, ErrSequenceGap
		}
		release := !s.snapshotting
		s.nextSequence = packet.Sequence + 1
		s.snapshotting = !packet.SnapshotLast
		s.synchronized = packet.SnapshotLast
		return s.events(packet, release), nil
	}
	if !s.synchronized || s.snapshotting {
		s.snapshotting = false
		return Result{ReleaseSource: true, RequestSnapshot: true}, ErrResyncRequired
	}
	if packet.Sequence != s.nextSequence {
		s.synchronized = false
		return Result{ReleaseSource: true, RequestSnapshot: true}, ErrSequenceGap
	}
	s.nextSequence++
	return s.events(packet, false), nil
}

func (s *Session) events(packet Packet, release bool) Result {
	result := Result{Count: int(packet.Count), ReleaseSource: release}
	for index := 0; index < result.Count; index++ {
		item := packet.Events[index]
		result.Events[index] = event.KeyEvent{
			Source:   packet.Source,
			Position: item.Position,
			Pressed:  item.Pressed,
			Sequence: packet.Sequence,
		}
	}
	return result
}

// Disconnect invalidates the connection state. The Primary must release this
// source even when it did not observe the final key-up notification.
func (s *Session) Disconnect() Result {
	s.synchronized = false
	s.snapshotting = false
	return Result{ReleaseSource: true, RequestSnapshot: true}
}

package split_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/keyboard/transport/split"
)

func packet(source event.Source, sequence uint16, snapshot bool, events ...split.PositionEvent) split.Packet {
	value := split.Packet{Source: source, Sequence: sequence, Snapshot: snapshot, SnapshotLast: snapshot, Count: uint8(len(events))}
	copy(value.Events[:], events)
	return value
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	want := packet(3, 42, false,
		split.PositionEvent{Position: 17, Pressed: true},
		split.PositionEvent{Position: 18, Pressed: false},
	)
	raw, err := split.Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := split.Decode(raw[:])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Decode(Encode()) = %#v, want %#v", got, want)
	}
}

func TestDecodeRejectsMalformedPackets(t *testing.T) {
	valid, err := split.Encode(packet(1, 1, true))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		raw  []byte
		want error
	}{
		{"truncated", valid[:split.PacketSize-1], split.ErrPacketSize},
		{"unknown version", func() []byte { raw := valid; raw[0] = 2; return raw[:] }(), split.ErrVersion},
		{"unknown flags", func() []byte { raw := valid; raw[1] = 0x80; return raw[:] }(), split.ErrFlags},
		{"normal empty", func() []byte { raw := valid; raw[1] = 0; return raw[:] }(), split.ErrEventCount},
		{"snapshot release", func() []byte { raw := valid; raw[3] = 1; raw[8] = 0; return raw[:] }(), split.ErrSnapshotEvent},
		{"invalid pressed", func() []byte { raw := valid; raw[3] = 1; raw[8] = 2; return raw[:] }(), split.ErrPressedValue},
		{"trailing bytes", func() []byte { raw := valid; raw[6] = 1; return raw[:] }(), split.ErrTrailingData},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := split.Decode(test.raw)
			if !errors.Is(err, test.want) {
				t.Fatalf("Decode() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestSessionAcceptsMultiPacketSnapshot(t *testing.T) {
	session := split.NewSession(2)
	first := packet(2, 40, true, split.PositionEvent{Position: 1, Pressed: true})
	first.SnapshotLast = false
	result, err := session.Process(first)
	if err != nil || !result.ReleaseSource || result.Count != 1 {
		t.Fatalf("first snapshot chunk = %#v, %v", result, err)
	}
	last := packet(2, 41, true, split.PositionEvent{Position: 2, Pressed: true})
	result, err = session.Process(last)
	if err != nil || result.ReleaseSource || result.Count != 1 {
		t.Fatalf("last snapshot chunk = %#v, %v", result, err)
	}
	if _, err = session.Process(packet(2, 42, false, split.PositionEvent{Position: 2, Pressed: false})); err != nil {
		t.Fatalf("edge after complete snapshot = %v", err)
	}
}

func TestSessionRequiresSnapshotAndRecoversFromGap(t *testing.T) {
	session := split.NewSession(2)
	result, err := session.Process(packet(2, 10, false, split.PositionEvent{Position: 1, Pressed: true}))
	if !errors.Is(err, split.ErrResyncRequired) || !result.ReleaseSource || !result.RequestSnapshot {
		t.Fatalf("first edge = %#v, %v; want release and snapshot request", result, err)
	}

	result, err = session.Process(packet(2, 10, true, split.PositionEvent{Position: 1, Pressed: true}))
	if err != nil || !result.ReleaseSource || result.Count != 1 || !result.Events[0].Pressed {
		t.Fatalf("snapshot = %#v, %v", result, err)
	}
	result, err = session.Process(packet(2, 11, false, split.PositionEvent{Position: 1, Pressed: false}))
	if err != nil || result.ReleaseSource || result.Count != 1 {
		t.Fatalf("next edge = %#v, %v", result, err)
	}

	result, err = session.Process(packet(2, 13, false, split.PositionEvent{Position: 2, Pressed: true}))
	if !errors.Is(err, split.ErrSequenceGap) || !result.ReleaseSource || !result.RequestSnapshot {
		t.Fatalf("gap = %#v, %v; want release and snapshot request", result, err)
	}
	if result, err = session.Process(packet(2, 14, false, split.PositionEvent{Position: 2, Pressed: true})); !errors.Is(err, split.ErrResyncRequired) || !result.ReleaseSource {
		t.Fatalf("edge after gap = %#v, %v; want resync required", result, err)
	}
}

func TestSessionDisconnectAndSourceIsolation(t *testing.T) {
	session := split.NewSession(4)
	if _, err := session.Process(packet(5, 1, true)); !errors.Is(err, split.ErrSourceMismatch) {
		t.Fatalf("source mismatch error = %v, want ErrSourceMismatch", err)
	}
	if result := session.Disconnect(); !result.ReleaseSource || !result.RequestSnapshot {
		t.Fatalf("Disconnect() = %#v, want release and snapshot request", result)
	}
}

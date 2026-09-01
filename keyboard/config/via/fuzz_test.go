package via_test

import "testing"

func FuzzHandleMalformedPackets(f *testing.F) {
	f.Add([]byte{0x01})
	f.Add(make([]byte, 32))
	f.Add([]byte{0xff, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, packet []byte) {
		adapter, _ := newAdapter(t)
		_, _ = adapter.Handle(packet, 0)
	})
}

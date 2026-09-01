package storage_test

import (
	"errors"
	"io"
	"testing"

	"github.com/tgk-project/tgk/keyboard/config/storage"
)

func TestFlashSlotsUsesReservedConsecutiveBlocks(t *testing.T) {
	device := &fakeFlash{bytes: make([]byte, 64), eraseSize: 16}
	slots, err := storage.NewFlashSlots(device, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := slots.EraseSlot(0); err != nil {
		t.Fatal(err)
	}
	if got := device.erasedBlock; got != 2 {
		t.Fatalf("erased block = %d, want 2", got)
	}
	if err := slots.WriteSlot(1, []byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	if got := device.bytes[48:51]; string(got) != string([]byte{1, 2, 3}) {
		t.Fatalf("slot write = %v, want [1 2 3]", got)
	}
}

func TestFlashSlotsRejectsUnavailableBlocks(t *testing.T) {
	_, err := storage.NewFlashSlots(&fakeFlash{bytes: make([]byte, 32), eraseSize: 16}, 1)
	if !errors.Is(err, storage.ErrFlashLayout) {
		t.Fatalf("NewFlashSlots() error = %v, want ErrFlashLayout", err)
	}
}

type fakeFlash struct {
	bytes       []byte
	eraseSize   int64
	erasedBlock int64
}

func (d *fakeFlash) ReadAt(dst []byte, offset int64) (int, error) {
	if offset < 0 || int(offset)+len(dst) > len(d.bytes) {
		return 0, io.EOF
	}
	return copy(dst, d.bytes[offset:]), nil
}

func (d *fakeFlash) WriteAt(src []byte, offset int64) (int, error) {
	if offset < 0 || int(offset)+len(src) > len(d.bytes) {
		return 0, io.ErrShortWrite
	}
	return copy(d.bytes[offset:], src), nil
}

func (d *fakeFlash) Size() int64 { return int64(len(d.bytes)) }

func (d *fakeFlash) EraseBlockSize() int64 { return d.eraseSize }

func (d *fakeFlash) EraseBlocks(start, length int64) error {
	if start < 0 || length != 1 || start*d.eraseSize+d.eraseSize > int64(len(d.bytes)) {
		return io.ErrUnexpectedEOF
	}
	d.erasedBlock = start
	for index := start * d.eraseSize; index < (start+1)*d.eraseSize; index++ {
		d.bytes[index] = 0xff
	}
	return nil
}

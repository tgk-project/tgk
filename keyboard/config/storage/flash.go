package storage

import (
	"errors"
	"io"
)

var ErrFlashLayout = errors.New("flash cannot provide two config erase blocks")

// FlashDevice is satisfied by TinyGo's machine.BlockDevice. Keeping this
// local interface avoids importing machine into the transport-independent
// storage core.
type FlashDevice interface {
	io.ReaderAt
	io.WriterAt
	Size() int64
	EraseBlockSize() int64
	EraseBlocks(start, length int64) error
}

// FlashSlots adapts two consecutive erase blocks to SlotDevice. On nRF52840,
// construct it with machine.Flash and two board-reserved block indexes.
type FlashSlots struct {
	device     FlashDevice
	firstBlock int64
	slotSize   int
}

// NewFlashSlots validates that two full erase blocks are available. The caller
// chooses firstBlock from the keyboard's board composition, not from core.
func NewFlashSlots(device FlashDevice, firstBlock int64) (*FlashSlots, error) {
	if device == nil || firstBlock < 0 || device.EraseBlockSize() <= 0 {
		return nil, ErrFlashLayout
	}
	size := device.EraseBlockSize()
	if firstBlock+2 > device.Size()/size || size > int64(^uint(0)>>1) {
		return nil, ErrFlashLayout
	}
	return &FlashSlots{device: device, firstBlock: firstBlock, slotSize: int(size)}, nil
}

func (d *FlashSlots) SlotSize() int { return d.slotSize }

func (d *FlashSlots) ReadSlot(slot int, dst []byte) error {
	if slot < 0 || slot > 1 || len(dst) != d.slotSize {
		return ErrInvalidSlot
	}
	n, err := d.device.ReadAt(dst, (d.firstBlock+int64(slot))*int64(d.slotSize))
	if err != nil {
		return err
	}
	if n != len(dst) {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (d *FlashSlots) EraseSlot(slot int) error {
	if slot < 0 || slot > 1 {
		return ErrInvalidSlot
	}
	return d.device.EraseBlocks(d.firstBlock+int64(slot), 1)
}

func (d *FlashSlots) WriteSlot(slot int, src []byte) error {
	if slot < 0 || slot > 1 || len(src) > d.slotSize {
		return ErrInvalidSlot
	}
	n, err := d.device.WriteAt(src, (d.firstBlock+int64(slot))*int64(d.slotSize))
	if err != nil {
		return err
	}
	if n < len(src) {
		return io.ErrShortWrite
	}
	return nil
}

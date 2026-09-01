package storage

import "errors"

var (
	ErrInvalidSlot = errors.New("invalid config storage slot")
	ErrNotErased   = errors.New("config storage write requires an erased slot")
)

// MemoryDevice is an in-memory two-slot device for host-side tests and local
// composition. It models NOR Flash's erase-before-write rule.
type MemoryDevice struct {
	slots [2][]byte
}

// NewMemoryDevice creates two erased slots of slotSize bytes.
func NewMemoryDevice(slotSize int) (*MemoryDevice, error) {
	if slotSize < HeaderSize {
		return nil, ErrSlotTooSmall
	}
	device := &MemoryDevice{}
	for slot := range device.slots {
		device.slots[slot] = make([]byte, slotSize)
		for index := range device.slots[slot] {
			device.slots[slot][index] = 0xff
		}
	}
	return device, nil
}

func (d *MemoryDevice) SlotSize() int { return len(d.slots[0]) }

func (d *MemoryDevice) ReadSlot(slot int, dst []byte) error {
	if !d.validSlot(slot) || len(dst) != len(d.slots[slot]) {
		return ErrInvalidSlot
	}
	copy(dst, d.slots[slot])
	return nil
}

func (d *MemoryDevice) EraseSlot(slot int) error {
	if !d.validSlot(slot) {
		return ErrInvalidSlot
	}
	for index := range d.slots[slot] {
		d.slots[slot][index] = 0xff
	}
	return nil
}

func (d *MemoryDevice) WriteSlot(slot int, src []byte) error {
	if !d.validSlot(slot) || len(src) > len(d.slots[slot]) {
		return ErrInvalidSlot
	}
	for index, value := range src {
		if d.slots[slot][index] != 0xff {
			return ErrNotErased
		}
		d.slots[slot][index] = value
	}
	return nil
}

// Corrupt flips one byte without erasing. It is deliberately exposed so
// host-side tests can model a damaged or torn slot.
func (d *MemoryDevice) Corrupt(slot, offset int, value byte) error {
	if !d.validSlot(slot) || offset < 0 || offset >= len(d.slots[slot]) {
		return ErrInvalidSlot
	}
	d.slots[slot][offset] = value
	return nil
}

func (d *MemoryDevice) validSlot(slot int) bool { return slot >= 0 && slot < len(d.slots) }

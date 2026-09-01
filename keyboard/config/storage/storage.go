// Package storage provides crash-safe, transport-independent User Config
// persistence. It deliberately knows nothing about a board, Flash controller,
// or keyboard event loop.
package storage

import (
	"encoding/binary"
	"errors"
	"hash/crc32"

	"github.com/tgk-project/tgk/keyboard/config"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

const (
	// HeaderSize is the fixed, little-endian on-storage header size. Two reserved
	// bytes after SchemaVersion keep the 32-bit fields naturally aligned.
	HeaderSize  = 24
	bindingSize = 10
)

var (
	ErrInvalidOptions = errors.New("invalid config storage options")
	ErrSlotTooSmall   = errors.New("config storage slot is too small")
	ErrSnapshotLarge  = errors.New("config snapshot does not fit in a slot")
	ErrInvalidBinding = errors.New("config snapshot contains an invalid binding")
	ErrWriteFailed    = errors.New("config storage write failed")
	// ErrMigrationRequired prevents a firmware with a different schema, layout,
	// or format identity from erasing an otherwise valid User Config. A caller
	// must run a defined migration or explicitly Reset before committing.
	ErrMigrationRequired = errors.New("config storage migration or factory reset required")
)

// SlotDevice owns exactly two independently erasable storage slots. A
// platform adapter can wrap TinyGo machine.Flash; tests use MemoryDevice.
// Writes must only change erased bytes, as required by typical NOR Flash.
type SlotDevice interface {
	SlotSize() int
	ReadSlot(slot int, dst []byte) error
	EraseSlot(slot int) error
	WriteSlot(slot int, src []byte) error
}

// Options binds a User Config store to one immutable Factory Config identity.
// Factory is held independently in firmware and is never written to User
// Config slots.
type Options struct {
	Magic         uint32
	SchemaVersion uint16
	LayoutID      uint32
	Factory       config.Snapshot
	DebounceTicks uint32
}

// Source tells callers whether the current configuration came from Factory
// Config or one of the two User Config slots.
type Source uint8

const (
	SourceFactory Source = iota
	SourceSlotA
	SourceSlotB
)

// Store keeps the mutable User Config in RAM. Update is intentionally
// non-blocking; the composition event loop decides when to call Commit or
// CommitDue outside scan, BLE, and USB callbacks.
type Store struct {
	device            SlotDevice
	options           Options
	current           config.Snapshot
	source            Source
	dirty             bool
	changed           uint32
	hasValue          bool
	migrationRequired bool
}

// New constructs a Store and leaves the active value at Factory Config until
// Load is called. Factory bindings are copied, so they cannot be mutated by a
// caller after construction.
func New(device SlotDevice, options Options) (*Store, error) {
	if device == nil || options.Magic == 0 || options.SchemaVersion == 0 || options.LayoutID == 0 {
		return nil, ErrInvalidOptions
	}
	if device.SlotSize() < HeaderSize {
		return nil, ErrSlotTooSmall
	}
	if err := validateBindings(options.Factory.Bindings); err != nil {
		return nil, err
	}
	if HeaderSize+len(options.Factory.Bindings)*bindingSize > device.SlotSize() {
		return nil, ErrSnapshotLarge
	}
	store := &Store{device: device, options: options, source: SourceFactory}
	store.current = store.factory()
	return store, nil
}

// Load selects the newest valid, compatible User Config slot. A missing,
// corrupt, torn, or incompatible User Config safely falls back to Factory
// Config instead of being overwritten.
func (s *Store) Load() (config.Snapshot, error) {
	a, validA := s.readSlot(0)
	b, validB := s.readSlot(1)

	s.source = SourceFactory
	s.current = s.factory()
	s.migrationRequired = (validA && !s.matchesIdentity(a)) || (validB && !s.matchesIdentity(b))
	validA = validA && s.matchesIdentity(a)
	validB = validB && s.matchesIdentity(b)
	if validA && (!validB || newer(a.Header.Generation, b.Header.Generation)) {
		s.current, s.source = a, SourceSlotA
	} else if validB {
		s.current, s.source = b, SourceSlotB
	}
	s.dirty = false
	s.hasValue = true
	return clone(s.current), nil
}

// Update replaces the RAM value immediately and marks it dirty. It does not
// perform storage I/O. Generation, identity fields, payload length, and CRC
// are assigned at Commit time rather than trusted from the caller.
func (s *Store) Update(value config.Snapshot, now uint32) error {
	if err := validateBindings(value.Bindings); err != nil {
		return err
	}
	if HeaderSize+len(value.Bindings)*bindingSize > s.device.SlotSize() {
		return ErrSnapshotLarge
	}
	// Keep the RAM value structurally complete as well: Config Service can
	// observe the new keymap immediately, while Commit later advances only its
	// generation number.
	s.current = s.withHeader(value.Bindings, s.current.Header.Generation)
	s.dirty = true
	s.changed = now
	s.hasValue = true
	return nil
}

// Commit writes a new generation to the inactive A/B slot. The existing valid
// slot remains untouched until the new data has been erased, written, and can
// be read back with a valid CRC.
func (s *Store) Commit() error {
	if !s.dirty {
		return nil
	}
	if s.migrationRequired || s.hasIncompatibleSlot() {
		s.migrationRequired = true
		return ErrMigrationRequired
	}
	if !s.hasValue {
		s.current = s.factory()
		s.hasValue = true
	}

	currentGeneration := s.current.Header.Generation
	if s.source == SourceFactory {
		currentGeneration = s.highestGeneration()
	}
	nextSlot := 0
	if s.source == SourceSlotA {
		nextSlot = 1
	} else if s.source == SourceSlotB {
		nextSlot = 0
	}
	snapshot := s.withHeader(s.current.Bindings, currentGeneration+1)
	encoded := encode(snapshot)
	if err := s.device.EraseSlot(nextSlot); err != nil {
		return err
	}
	if err := s.device.WriteSlot(nextSlot, encoded); err != nil {
		return errors.Join(ErrWriteFailed, err)
	}
	verified, valid := s.readSlot(nextSlot)
	if !valid || verified.Header.Generation != snapshot.Header.Generation {
		return ErrWriteFailed
	}
	s.current = verified
	if nextSlot == 0 {
		s.source = SourceSlotA
	} else {
		s.source = SourceSlotB
	}
	s.dirty = false
	return nil
}

// CommitDue commits a staged configuration only after DebounceTicks has
// elapsed. A zero debounce commits immediately. Tick arithmetic is safe across
// a uint32 timer wrap for intervals shorter than 2^31 ticks.
func (s *Store) CommitDue(now uint32) (bool, error) {
	if !s.dirty || (s.options.DebounceTicks != 0 && uint32(now-s.changed) < s.options.DebounceTicks) {
		return false, nil
	}
	return true, s.Commit()
}

// Save implements config.ConfigStore as an explicit, synchronous commit.
// Runtime configuration paths should prefer Update followed by CommitDue.
func (s *Store) Save(value config.Snapshot) error {
	if err := s.Update(value, 0); err != nil {
		return err
	}
	return s.Commit()
}

// Reset erases only User Config and restores the in-memory Factory Config.
// It intentionally has no access to bond data.
func (s *Store) Reset() error {
	if err := s.device.EraseSlot(0); err != nil {
		return err
	}
	if err := s.device.EraseSlot(1); err != nil {
		return err
	}
	s.current = s.factory()
	s.source = SourceFactory
	s.dirty = false
	s.hasValue = true
	s.migrationRequired = false
	return nil
}

// Current returns a copy of the immediately active RAM configuration.
func (s *Store) Current() config.Snapshot { return clone(s.current) }

// Source reports where Current was recovered from.
func (s *Store) Source() Source { return s.source }

// Dirty reports whether RAM changes still need an explicit or debounced commit.
func (s *Store) Dirty() bool { return s.dirty }

// MigrationRequired reports that a valid User Config exists but does not
// match this store's identity. The active RAM configuration remains Factory
// Config, and Commit will return ErrMigrationRequired until Reset or a
// future migration path handles the old data explicitly.
func (s *Store) MigrationRequired() bool { return s.migrationRequired }

func (s *Store) factory() config.Snapshot {
	return s.withHeader(s.options.Factory.Bindings, 0)
}

func (s *Store) withHeader(bindings []keymap.Binding, generation uint32) config.Snapshot {
	snapshot := config.Snapshot{
		Header: config.Header{
			Magic:         s.options.Magic,
			SchemaVersion: s.options.SchemaVersion,
			LayoutID:      s.options.LayoutID,
			Generation:    generation,
			PayloadLength: uint32(len(bindings) * bindingSize),
		},
		Bindings: cloneBindings(bindings),
	}
	snapshot.Header.CRC32 = checksum(snapshot.Header, snapshot.Bindings)
	return snapshot
}

func (s *Store) highestGeneration() uint32 {
	a, validA := s.readSlot(0)
	b, validB := s.readSlot(1)
	if validA && (!validB || newer(a.Header.Generation, b.Header.Generation)) {
		return a.Header.Generation
	}
	if validB {
		return b.Header.Generation
	}
	return 0
}

func (s *Store) hasIncompatibleSlot() bool {
	for slot := 0; slot < 2; slot++ {
		snapshot, valid := s.readSlot(slot)
		if valid && !s.matchesIdentity(snapshot) {
			return true
		}
	}
	return false
}

func (s *Store) matchesIdentity(snapshot config.Snapshot) bool {
	return snapshot.Header.Magic == s.options.Magic &&
		snapshot.Header.SchemaVersion == s.options.SchemaVersion &&
		snapshot.Header.LayoutID == s.options.LayoutID
}

func (s *Store) readSlot(slot int) (config.Snapshot, bool) {
	buffer := make([]byte, s.device.SlotSize())
	if err := s.device.ReadSlot(slot, buffer); err != nil {
		return config.Snapshot{}, false
	}
	snapshot, err := decode(buffer)
	if err != nil {
		return config.Snapshot{}, false
	}
	return snapshot, true
}

func encode(snapshot config.Snapshot) []byte {
	buffer := make([]byte, HeaderSize+len(snapshot.Bindings)*bindingSize)
	binary.LittleEndian.PutUint32(buffer[0:4], snapshot.Header.Magic)
	binary.LittleEndian.PutUint16(buffer[4:6], snapshot.Header.SchemaVersion)
	binary.LittleEndian.PutUint32(buffer[8:12], snapshot.Header.LayoutID)
	binary.LittleEndian.PutUint32(buffer[12:16], snapshot.Header.Generation)
	binary.LittleEndian.PutUint32(buffer[16:20], snapshot.Header.PayloadLength)
	binary.LittleEndian.PutUint32(buffer[20:24], snapshot.Header.CRC32)
	for index, binding := range snapshot.Bindings {
		offset := HeaderSize + index*bindingSize
		binary.LittleEndian.PutUint16(buffer[offset:offset+2], uint16(binding.Behavior))
		binary.LittleEndian.PutUint32(buffer[offset+2:offset+6], binding.Param1)
		binary.LittleEndian.PutUint32(buffer[offset+6:offset+10], binding.Param2)
	}
	return buffer
}

func decode(buffer []byte) (config.Snapshot, error) {
	if len(buffer) < HeaderSize {
		return config.Snapshot{}, ErrSlotTooSmall
	}
	header := config.Header{
		Magic:         binary.LittleEndian.Uint32(buffer[0:4]),
		SchemaVersion: binary.LittleEndian.Uint16(buffer[4:6]),
		LayoutID:      binary.LittleEndian.Uint32(buffer[8:12]),
		Generation:    binary.LittleEndian.Uint32(buffer[12:16]),
		PayloadLength: binary.LittleEndian.Uint32(buffer[16:20]),
		CRC32:         binary.LittleEndian.Uint32(buffer[20:24]),
	}
	if buffer[6] != 0 || buffer[7] != 0 ||
		header.PayloadLength%bindingSize != 0 || header.PayloadLength > uint32(len(buffer)-HeaderSize) {
		return config.Snapshot{}, ErrSnapshotLarge
	}
	bindings := make([]keymap.Binding, int(header.PayloadLength)/bindingSize)
	for index := range bindings {
		offset := HeaderSize + index*bindingSize
		bindings[index] = keymap.Binding{
			Behavior: keymap.BehaviorID(binary.LittleEndian.Uint16(buffer[offset : offset+2])),
			Param1:   binary.LittleEndian.Uint32(buffer[offset+2 : offset+6]),
			Param2:   binary.LittleEndian.Uint32(buffer[offset+6 : offset+10]),
		}
	}
	if err := validateBindings(bindings); err != nil || checksum(header, bindings) != header.CRC32 {
		return config.Snapshot{}, ErrInvalidBinding
	}
	return config.Snapshot{Header: header, Bindings: bindings}, nil
}

func checksum(header config.Header, bindings []keymap.Binding) uint32 {
	buffer := make([]byte, 20+len(bindings)*bindingSize)
	binary.LittleEndian.PutUint32(buffer[0:4], header.Magic)
	binary.LittleEndian.PutUint16(buffer[4:6], header.SchemaVersion)
	binary.LittleEndian.PutUint32(buffer[8:12], header.LayoutID)
	binary.LittleEndian.PutUint32(buffer[12:16], header.Generation)
	binary.LittleEndian.PutUint32(buffer[16:20], header.PayloadLength)
	for index, binding := range bindings {
		offset := 20 + index*bindingSize
		binary.LittleEndian.PutUint16(buffer[offset:offset+2], uint16(binding.Behavior))
		binary.LittleEndian.PutUint32(buffer[offset+2:offset+6], binding.Param1)
		binary.LittleEndian.PutUint32(buffer[offset+6:offset+10], binding.Param2)
	}
	return crc32.ChecksumIEEE(buffer)
}

func validateBindings(bindings []keymap.Binding) error {
	for _, binding := range bindings {
		if err := keymap.ValidateBinding(binding); err != nil {
			return ErrInvalidBinding
		}
	}
	return nil
}

func newer(a, b uint32) bool { return int32(a-b) > 0 }

func clone(value config.Snapshot) config.Snapshot {
	value.Bindings = cloneBindings(value.Bindings)
	return value
}

func cloneBindings(bindings []keymap.Binding) []keymap.Binding {
	return append([]keymap.Binding(nil), bindings...)
}

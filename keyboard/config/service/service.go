// Package service owns runtime keymap configuration independently of any
// particular Config transport or command protocol.
package service

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/config"
	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

var (
	ErrInvalidOptions = errors.New("invalid dynamic config service options")
	ErrLocked         = errors.New("configuration writes require physical unlock")
	ErrInvalidAddress = errors.New("keymap address is outside the keyboard definition")
	ErrInvalidBuffer  = errors.New("keymap buffer has an invalid size")
)

// KeymapSink is normally the Primary keymap engine. It is deliberately small
// so Config Service does not depend on engine, USB, BLE, or Flash packages.
type KeymapSink interface {
	ReplaceKeymap(keymap.Keymap) error
}

// Options is immutable composition data supplied by a keyboard definition.
// Factory is a layer-major flat binding table and never changes at runtime.
type Options struct {
	Identity    definition.Identity
	Layers      uint8
	Positions   uint16
	Factory     []keymap.Binding
	Store       config.ConfigStore
	Sink        KeymapSink
	UnlockTicks uint32
}

// Service owns the active RAM keymap and stages durable writes for an
// event-loop controlled Commit call. USB callbacks must not call it directly.
type Service struct {
	options Options
	active  []keymap.Binding
	dirty   bool
	reset   bool

	unlockedAt uint32
	unlocked   bool
}

// New loads the persisted User Config (or its Store's Factory fallback) and
// immediately applies it to the Primary keymap sink.
func New(options Options) (*Service, error) {
	if options.Store == nil || options.Sink == nil || options.Layers == 0 || options.Positions == 0 ||
		len(options.Factory) != int(options.Layers)*int(options.Positions) {
		return nil, ErrInvalidOptions
	}
	if _, err := makeKeymap(options.Factory, options.Layers, options.Positions); err != nil {
		return nil, err
	}
	loaded, err := options.Store.Load()
	if err != nil {
		return nil, err
	}
	if len(loaded.Bindings) == 0 {
		loaded.Bindings = options.Factory
	}
	if _, err := makeKeymap(loaded.Bindings, options.Layers, options.Positions); err != nil {
		return nil, err
	}
	result := &Service{options: options, active: clone(loaded.Bindings)}
	if err := result.apply(result.active); err != nil {
		return nil, err
	}
	return result, nil
}

// Identity returns the immutable keyboard identity for transport adapters and
// keyboard-definition lookup.
func (s *Service) Identity() definition.Identity { return s.options.Identity }

// LayerCount and PositionCount describe the fixed keyboard geometry.
func (s *Service) LayerCount() uint8     { return s.options.Layers }
func (s *Service) PositionCount() uint16 { return s.options.Positions }

// Binding returns one active binding by layer and dense physical position.
func (s *Service) Binding(layer uint8, position uint16) (keymap.Binding, bool) {
	index, err := s.index(layer, position)
	if err != nil {
		return keymap.Binding{}, false
	}
	return s.active[index], true
}

// Bindings returns a copy in layer-major order.
func (s *Service) Bindings() []keymap.Binding { return clone(s.active) }

// Unlock authorizes destructive commands for UnlockTicks from now. The caller
// must invoke this only after its board-specific physical gesture succeeds.
func (s *Service) Unlock(now uint32) {
	s.unlocked = true
	s.unlockedAt = now
}

// Lock immediately removes write authorization.
func (s *Service) Lock() { s.unlocked = false }

// Locked reports whether the physical unlock window is absent or expired.
func (s *Service) Locked(now uint32) bool {
	if !s.unlocked {
		return true
	}
	if s.options.UnlockTicks != 0 && uint32(now-s.unlockedAt) >= s.options.UnlockTicks {
		s.unlocked = false
		return true
	}
	return false
}

// SetBinding applies a fully validated update to RAM immediately. It does not
// access Flash; composition calls Commit later at a safe event-loop boundary.
func (s *Service) SetBinding(now uint32, layer uint8, position uint16, value keymap.Binding) error {
	if s.Locked(now) {
		return ErrLocked
	}
	index, err := s.index(layer, position)
	if err != nil {
		return err
	}
	next := clone(s.active)
	next[index] = value
	return s.replace(next)
}

// ReplaceBindings atomically replaces the complete flat keymap after physical
// unlock. It is used by transport adapters that support bulk keymap writes.
func (s *Service) ReplaceBindings(now uint32, value []keymap.Binding) error {
	if s.Locked(now) {
		return ErrLocked
	}
	if len(value) != len(s.active) {
		return ErrInvalidBuffer
	}
	return s.replace(value)
}

// FactoryReset restores the embedded Factory Config in RAM. Persistence is
// deferred until Commit, so no Flash work occurs while parsing a RAW packet.
func (s *Service) FactoryReset(now uint32) error {
	if s.Locked(now) {
		return ErrLocked
	}
	if err := s.apply(s.options.Factory); err != nil {
		return err
	}
	s.active = clone(s.options.Factory)
	s.dirty = true
	s.reset = true
	return nil
}

// Commit performs the pending persistence operation. It must be called from
// the Primary event loop, never a USB callback. A successful FactoryReset
// calls Store.Reset so only User Config is erased; bond state stays separate.
func (s *Service) Commit() error {
	if !s.dirty {
		return nil
	}
	if s.reset {
		if err := s.options.Store.Reset(); err != nil {
			return err
		}
	} else if err := s.options.Store.Save(config.Snapshot{Bindings: clone(s.active)}); err != nil {
		return err
	}
	s.dirty = false
	s.reset = false
	return nil
}

// Dirty reports whether a later Commit is required for restart persistence.
func (s *Service) Dirty() bool { return s.dirty }

func (s *Service) replace(value []keymap.Binding) error {
	if _, err := makeKeymap(value, s.options.Layers, s.options.Positions); err != nil {
		return err
	}
	if err := s.apply(value); err != nil {
		return err
	}
	s.active = clone(value)
	s.dirty = true
	s.reset = false
	return nil
}

func (s *Service) apply(value []keymap.Binding) error {
	km, err := makeKeymap(value, s.options.Layers, s.options.Positions)
	if err != nil {
		return err
	}
	return s.options.Sink.ReplaceKeymap(km)
}

func (s *Service) index(layer uint8, position uint16) (int, error) {
	if layer >= s.options.Layers || position >= s.options.Positions {
		return 0, ErrInvalidAddress
	}
	return int(layer)*int(s.options.Positions) + int(position), nil
}

func makeKeymap(bindings []keymap.Binding, layers uint8, positions uint16) (keymap.Keymap, error) {
	if len(bindings) != int(layers)*int(positions) {
		return keymap.Keymap{}, ErrInvalidBuffer
	}
	values := make([][]keymap.Binding, layers)
	for layer := range values {
		start := layer * int(positions)
		values[layer] = append([]keymap.Binding(nil), bindings[start:start+int(positions)]...)
	}
	return keymap.New(values)
}

func clone(value []keymap.Binding) []keymap.Binding {
	result := make([]keymap.Binding, len(value))
	copy(result, value)
	return result
}

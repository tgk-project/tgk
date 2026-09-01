// Package definition declares the public, hardware-independent composition
// boundary for a keyboard package.
package definition

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/keymap"
)

var (
	ErrMissingIdentity   = errors.New("keyboard definition requires a non-zero identity")
	ErrEmptyName         = errors.New("keyboard definition requires a product name")
	ErrPositionCount     = errors.New("keyboard definition position count does not match keymap")
	ErrDuplicatePosition = errors.New("keyboard definition contains a duplicate physical position")
)

// Identity is the immutable identity supplied by a keyboard package. LayoutID
// also binds a User Config to the physical layout that created it.
type Identity struct {
	VendorID    uint16
	ProductID   uint16
	ProductName string
	LayoutID    uint32
}

// Definition is owned by a keyboard package. Positions maps the keymap's
// dense index to the physical Position emitted by its scanner; it never holds
// GPIO, ADC, transport, or board details.
type Definition struct {
	Identity  Identity
	Positions []uint16
	Layers    [][]keymap.Binding
}

// FactoryConfig is a validated, immutable configuration embedded by a
// keyboard package. It is intentionally independent of mutable User Config.
type FactoryConfig struct {
	identity  Identity
	positions []uint16
	keymap    keymap.Keymap
	bindings  []keymap.Binding
}

// New validates a keyboard package's factory configuration and copies its
// data. The returned values can safely be retained by a Primary composition.
func New(value Definition) (FactoryConfig, error) {
	if value.Identity.VendorID == 0 || value.Identity.ProductID == 0 || value.Identity.LayoutID == 0 {
		return FactoryConfig{}, ErrMissingIdentity
	}
	if value.Identity.ProductName == "" {
		return FactoryConfig{}, ErrEmptyName
	}

	km, err := keymap.New(value.Layers)
	if err != nil {
		return FactoryConfig{}, err
	}
	if len(value.Positions) != int(km.PositionCount()) {
		return FactoryConfig{}, ErrPositionCount
	}

	positions := make([]uint16, len(value.Positions))
	seen := make(map[uint16]struct{}, len(positions))
	for index, position := range value.Positions {
		if _, duplicate := seen[position]; duplicate {
			return FactoryConfig{}, ErrDuplicatePosition
		}
		seen[position] = struct{}{}
		positions[index] = position
	}

	bindings := make([]keymap.Binding, 0, int(km.LayerCount())*int(km.PositionCount()))
	for layer := uint8(0); layer < km.LayerCount(); layer++ {
		for position := uint16(0); position < km.PositionCount(); position++ {
			binding, _ := km.BindingAt(layer, position)
			bindings = append(bindings, binding)
		}
	}
	return FactoryConfig{
		identity:  value.Identity,
		positions: positions,
		keymap:    km,
		bindings:  bindings,
	}, nil
}

// Identity returns a value copy of the keyboard identity.
func (f FactoryConfig) Identity() Identity { return f.identity }

// Keymap returns the validated Factory keymap. keymap.Keymap is immutable.
func (f FactoryConfig) Keymap() keymap.Keymap { return f.keymap }

// Positions returns a copy of the scanner's physical position map.
func (f FactoryConfig) Positions() []uint16 {
	result := make([]uint16, len(f.positions))
	copy(result, f.positions)
	return result
}

// Bindings returns a copy in layer-major order for a future User Config store.
// It contains stable behavior IDs and numeric parameters only.
func (f FactoryConfig) Bindings() []keymap.Binding {
	result := make([]keymap.Binding, len(f.bindings))
	copy(result, f.bindings)
	return result
}

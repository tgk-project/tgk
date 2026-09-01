// Package keymap defines the persistent, transport-independent keymap model.
package keymap

import "errors"

// BehaviorID is a stable value suitable for persistence. It is deliberately not
// a Go type name or a pointer to an implementation.
type BehaviorID uint16

const (
	BehaviorNone BehaviorID = iota
	BehaviorKey
	BehaviorModifier
	BehaviorConsumer
	BehaviorTransparent
	BehaviorMomentaryLayer // MO
	BehaviorToLayer        // TO
	BehaviorToggleLayer    // TG
)

var (
	ErrNoLayers          = errors.New("keymap must contain at least one layer")
	ErrEmptyLayer        = errors.New("keymap layers must contain at least one position")
	ErrUnevenLayer       = errors.New("keymap layers must have the same position count")
	ErrUnknownBehavior   = errors.New("keymap contains an unknown behavior")
	ErrInvalidParameter  = errors.New("keymap behavior has an invalid parameter")
	ErrInvalidLayerIndex = errors.New("keymap behavior refers to an invalid layer")
)

// Binding is a behavior identifier and its fixed parameters. Param1 and Param2
// are kept numeric so snapshots can be encoded without Go runtime metadata.
type Binding struct {
	Behavior BehaviorID
	Param1   uint32
	Param2   uint32
}

// Keymap is an immutable position-to-binding table. Construct it with New.
type Keymap struct {
	layers    [][]Binding
	positions uint16
}

// New validates and copies layers, isolating the engine from later caller
// mutation of the configuration buffer.
func New(layers [][]Binding) (Keymap, error) {
	if len(layers) == 0 {
		return Keymap{}, ErrNoLayers
	}
	if len(layers) > 1<<8-1 {
		return Keymap{}, ErrInvalidParameter
	}
	if len(layers[0]) == 0 {
		return Keymap{}, ErrEmptyLayer
	}
	if len(layers[0]) > 1<<16-1 {
		return Keymap{}, ErrInvalidParameter
	}

	positions := len(layers[0])
	for _, layer := range layers {
		if len(layer) != positions {
			return Keymap{}, ErrUnevenLayer
		}
		for _, binding := range layer {
			if err := ValidateBinding(binding); err != nil {
				return Keymap{}, err
			}
			if isLayerBehavior(binding.Behavior) && binding.Param1 >= uint32(len(layers)) {
				return Keymap{}, ErrInvalidLayerIndex
			}
		}
	}

	copyLayers := make([][]Binding, len(layers))
	for i, layer := range layers {
		copyLayers[i] = make([]Binding, len(layer))
		copy(copyLayers[i], layer)
	}
	return Keymap{layers: copyLayers, positions: uint16(positions)}, nil
}

// ValidateBinding checks whether a binding can be represented by the core.
// Layer-target bounds depend on a complete Keymap and are checked by New.
// Keeping the per-binding validation public lets persistence reject malformed
// data before it reaches the Primary-owned keymap engine.
func ValidateBinding(binding Binding) error {
	if !knownBehavior(binding.Behavior) {
		return ErrUnknownBehavior
	}
	if !validParameter(binding) {
		return ErrInvalidParameter
	}
	return nil
}

// LayerCount returns the number of available layers.
func (k Keymap) LayerCount() uint8 { return uint8(len(k.layers)) }

// PositionCount returns the number of physical positions per layer.
func (k Keymap) PositionCount() uint16 { return k.positions }

// BindingAt returns the binding for a layer and a physical position.
func (k Keymap) BindingAt(layer uint8, position uint16) (Binding, bool) {
	if int(layer) >= len(k.layers) || position >= k.positions {
		return Binding{}, false
	}
	return k.layers[layer][position], true
}

func knownBehavior(id BehaviorID) bool {
	return id >= BehaviorNone && id <= BehaviorToggleLayer
}

func isLayerBehavior(id BehaviorID) bool {
	return id == BehaviorMomentaryLayer || id == BehaviorToLayer || id == BehaviorToggleLayer
}

func validParameter(binding Binding) bool {
	switch binding.Behavior {
	case BehaviorNone, BehaviorTransparent:
		return binding.Param1 == 0 && binding.Param2 == 0
	case BehaviorKey:
		return binding.Param1 > 0 && binding.Param1 <= 0xff && binding.Param2 == 0
	case BehaviorModifier:
		return binding.Param1 > 0 && binding.Param1 <= 0xff && binding.Param2 == 0
	case BehaviorConsumer:
		return binding.Param1 > 0 && binding.Param1 <= 0xffff && binding.Param2 == 0
	case BehaviorMomentaryLayer, BehaviorToLayer, BehaviorToggleLayer:
		return binding.Param2 == 0
	default:
		return false
	}
}

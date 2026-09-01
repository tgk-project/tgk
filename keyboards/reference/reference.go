// Package reference is the minimal XIAO BLE reference keyboard definition.
// Board wiring belongs in its firmware composition, not in this package.
package reference

import (
	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

const (
	VendorID  uint16 = 0x1209
	ProductID uint16 = 0x0001
	LayoutID  uint32 = 0x54474b01
)

// FactoryConfig returns the immutable default for the one-key reference build.
func FactoryConfig() (definition.FactoryConfig, error) {
	return definition.New(definition.Definition{
		Identity: definition.Identity{
			VendorID: VendorID, ProductID: ProductID, ProductName: "TGK XIAO BLE Reference", LayoutID: LayoutID,
		},
		Positions: []uint16{0},
		Layers: [][]keymap.Binding{{
			{Behavior: keymap.BehaviorKey, Param1: 4}, // HID Keyboard A
		}},
	})
}

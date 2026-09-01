// Package template is a second, minimal keyboard package that demonstrates
// that a new layout composes TGK through keyboard/definition alone.
package template

import (
	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

// FactoryConfig returns a two-position example for a separate keyboard.
func FactoryConfig() (definition.FactoryConfig, error) {
	return definition.New(definition.Definition{
		Identity: definition.Identity{
			VendorID: 0x1209, ProductID: 0x0002, ProductName: "TGK Keyboard Template", LayoutID: 0x54474b02,
		},
		Positions: []uint16{10, 20},
		Layers: [][]keymap.Binding{{
			{Behavior: keymap.BehaviorKey, Param1: 4},
			{Behavior: keymap.BehaviorKey, Param1: 5},
		}},
	})
}

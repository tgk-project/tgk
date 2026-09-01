package definition_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

func validDefinition() definition.Definition {
	return definition.Definition{
		Identity: definition.Identity{
			VendorID: 0x1209, ProductID: 0x0001, ProductName: "Test keyboard", LayoutID: 0x01,
		},
		Positions: []uint16{12, 4},
		Layers: [][]keymap.Binding{{
			{Behavior: keymap.BehaviorKey, Param1: 4},
			{Behavior: keymap.BehaviorKey, Param1: 5},
		}},
	}
}

func TestNewBuildsIndependentFactoryConfig(t *testing.T) {
	value := validDefinition()
	factory, err := definition.New(value)
	if err != nil {
		t.Fatal(err)
	}
	value.Positions[0] = 99
	value.Layers[0][0].Param1 = 9

	if got, want := factory.Positions(), []uint16{12, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Positions() = %v, want %v", got, want)
	}
	binding, ok := factory.Keymap().BindingAt(0, 0)
	if !ok || binding.Param1 != 4 {
		t.Fatalf("Factory Keymap() binding = %+v, %t; want usage 4", binding, ok)
	}
	bindings := factory.Bindings()
	bindings[0].Param1 = 10
	if got := factory.Bindings()[0].Param1; got != 4 {
		t.Fatalf("Bindings() exposed mutable Factory Config: %d", got)
	}
}

func TestNewRejectsInvalidComposition(t *testing.T) {
	tests := []struct {
		name string
		edit func(*definition.Definition)
		want error
	}{
		{
			name: "missing layout identity",
			edit: func(value *definition.Definition) { value.Identity.LayoutID = 0 },
			want: definition.ErrMissingIdentity,
		},
		{
			name: "position count",
			edit: func(value *definition.Definition) { value.Positions = value.Positions[:1] },
			want: definition.ErrPositionCount,
		},
		{
			name: "duplicate physical position",
			edit: func(value *definition.Definition) { value.Positions[1] = value.Positions[0] },
			want: definition.ErrDuplicatePosition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := validDefinition()
			tt.edit(&value)
			_, err := definition.New(value)
			if !errors.Is(err, tt.want) {
				t.Fatalf("New() error = %v, want %v", err, tt.want)
			}
		})
	}
}

package keymap_test

import (
	"errors"
	"testing"

	"github.com/tgk-project/tgk/keyboard/keymap"
)

func TestNewRejectsInvalidTables(t *testing.T) {
	tests := []struct {
		name   string
		layers [][]keymap.Binding
		want   error
	}{
		{name: "no layers", want: keymap.ErrNoLayers},
		{name: "empty layer", layers: [][]keymap.Binding{{}}, want: keymap.ErrEmptyLayer},
		{
			name: "uneven layers",
			layers: [][]keymap.Binding{
				{{Behavior: keymap.BehaviorKey, Param1: 4}},
				{{Behavior: keymap.BehaviorKey, Param1: 4}, {Behavior: keymap.BehaviorKey, Param1: 5}},
			},
			want: keymap.ErrUnevenLayer,
		},
		{
			name:   "unknown behavior",
			layers: [][]keymap.Binding{{{Behavior: 99}}},
			want:   keymap.ErrUnknownBehavior,
		},
		{
			name:   "unknown layer",
			layers: [][]keymap.Binding{{{Behavior: keymap.BehaviorMomentaryLayer, Param1: 1}}},
			want:   keymap.ErrInvalidLayerIndex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := keymap.New(tt.layers)
			if !errors.Is(err, tt.want) {
				t.Fatalf("New() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewCopiesInput(t *testing.T) {
	layers := [][]keymap.Binding{{{Behavior: keymap.BehaviorKey, Param1: 4}}}
	km, err := keymap.New(layers)
	if err != nil {
		t.Fatal(err)
	}
	layers[0][0] = keymap.Binding{Behavior: keymap.BehaviorKey, Param1: 5}
	binding, ok := km.BindingAt(0, 0)
	if !ok || binding.Param1 != 4 {
		t.Fatalf("BindingAt() = %+v, %t; keymap should be immutable", binding, ok)
	}
}

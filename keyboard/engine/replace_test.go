package engine_test

import (
	"testing"

	"github.com/tgk-project/tgk/keyboard/keymap"
)

func TestReplaceKeymapReleasesPressedStateBeforeApplyingNewBinding(t *testing.T) {
	host := &fakeHost{ready: true}
	value := newEngine(t, host)
	runEvents(t, value, keyEvent(1, true))

	layers := make([][]keymap.Binding, 3)
	for layer := range layers {
		layers[layer] = make([]keymap.Binding, 6)
		for position := range layers[layer] {
			layers[layer][position] = keymap.Binding{Behavior: keymap.BehaviorNone}
		}
	}
	layers[0][1] = keymap.Binding{Behavior: keymap.BehaviorKey, Param1: keyC}
	replacement, err := keymap.New(layers)
	if err != nil {
		t.Fatal(err)
	}
	if err := value.ReplaceKeymap(replacement); err != nil {
		t.Fatal(err)
	}
	if host.releases != 1 {
		t.Fatalf("ReleaseAll calls = %d, want 1", host.releases)
	}
	runEvents(t, value, keyEvent(1, true))
	if got := value.Report().Keys[0]; got != keyC {
		t.Fatalf("replacement key usage = %d, want %d", got, keyC)
	}
}

package service_test

import (
	"errors"
	"testing"

	"github.com/tgk-project/tgk/keyboard/config"
	"github.com/tgk-project/tgk/keyboard/config/service"
	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

type memoryStore struct {
	value  config.Snapshot
	saves  int
	resets int
}

func (s *memoryStore) Load() (config.Snapshot, error) { return s.value, nil }
func (s *memoryStore) Save(value config.Snapshot) error {
	s.value = value
	s.saves++
	return nil
}
func (s *memoryStore) Reset() error {
	s.value = config.Snapshot{}
	s.resets++
	return nil
}

type sink struct{ value keymap.Keymap }

func (s *sink) ReplaceKeymap(value keymap.Keymap) error { s.value = value; return nil }

func newService(t *testing.T, stored []keymap.Binding) (*service.Service, *memoryStore, *sink) {
	t.Helper()
	store := &memoryStore{value: config.Snapshot{Bindings: stored}}
	target := &sink{}
	value, err := service.New(service.Options{
		Identity: definition.Identity{VendorID: 0x1209, ProductID: 1, ProductName: "test", LayoutID: 1},
		Layers:   1, Positions: 2, UnlockTicks: 10,
		Factory: []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: 4}, {Behavior: keymap.BehaviorKey, Param1: 5}},
		Store:   store, Sink: target,
	})
	if err != nil {
		t.Fatal(err)
	}
	return value, store, target
}

func TestWriteRequiresPhysicalUnlockAndAppliesImmediately(t *testing.T) {
	value, store, target := newService(t, nil)
	if err := value.SetBinding(100, 0, 0, keymap.Binding{Behavior: keymap.BehaviorKey, Param1: 6}); !errors.Is(err, service.ErrLocked) {
		t.Fatalf("locked SetBinding error = %v", err)
	}
	value.Unlock(100)
	if err := value.SetBinding(109, 0, 0, keymap.Binding{Behavior: keymap.BehaviorKey, Param1: 6}); err != nil {
		t.Fatal(err)
	}
	if got, _ := value.Binding(0, 0); got.Param1 != 6 {
		t.Fatalf("RAM binding = %#v, want usage 6", got)
	}
	if got, _ := target.value.BindingAt(0, 0); got.Param1 != 6 {
		t.Fatalf("sink binding = %#v, want usage 6", got)
	}
	if store.saves != 0 || !value.Dirty() {
		t.Fatal("SetBinding persisted instead of staging an immediate RAM update")
	}
	if err := value.Commit(); err != nil {
		t.Fatal(err)
	}
	if store.saves != 1 || value.Dirty() {
		t.Fatalf("Commit saves=%d dirty=%t", store.saves, value.Dirty())
	}
	restartedSink := &sink{}
	restarted, err := service.New(service.Options{
		Identity: definition.Identity{VendorID: 0x1209, ProductID: 1, ProductName: "test", LayoutID: 1},
		Layers:   1, Positions: 2, UnlockTicks: 10,
		Factory: []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: 4}, {Behavior: keymap.BehaviorKey, Param1: 5}},
		Store:   store, Sink: restartedSink,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := restarted.Binding(0, 0); got.Param1 != 6 {
		t.Fatalf("committed restart binding = %#v, want usage 6", got)
	}
	if err := value.SetBinding(110, 0, 1, keymap.Binding{Behavior: keymap.BehaviorKey, Param1: 7}); !errors.Is(err, service.ErrLocked) {
		t.Fatalf("expired unlock error = %v", err)
	}
}

func TestFactoryResetIsDeferredAndDoesNotUseSave(t *testing.T) {
	value, store, _ := newService(t, nil)
	value.Unlock(1)
	if err := value.SetBinding(1, 0, 0, keymap.Binding{Behavior: keymap.BehaviorKey, Param1: 6}); err != nil {
		t.Fatal(err)
	}
	if err := value.FactoryReset(2); err != nil {
		t.Fatal(err)
	}
	if got, _ := value.Binding(0, 0); got.Param1 != 4 {
		t.Fatalf("FactoryReset RAM binding = %#v", got)
	}
	if store.resets != 0 {
		t.Fatal("FactoryReset performed persistence before Commit")
	}
	if err := value.Commit(); err != nil {
		t.Fatal(err)
	}
	if store.resets != 1 || store.saves != 0 {
		t.Fatalf("reset=%d saves=%d", store.resets, store.saves)
	}
}

func TestRejectsWrongGeometryAndInvalidLayerTarget(t *testing.T) {
	value, _, _ := newService(t, nil)
	value.Unlock(0)
	if err := value.ReplaceBindings(0, []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: 4}}); !errors.Is(err, service.ErrInvalidBuffer) {
		t.Fatalf("wrong geometry error = %v", err)
	}
	if err := value.SetBinding(0, 0, 0, keymap.Binding{Behavior: keymap.BehaviorMomentaryLayer, Param1: 1}); !errors.Is(err, keymap.ErrInvalidLayerIndex) {
		t.Fatalf("invalid layer target error = %v", err)
	}
}

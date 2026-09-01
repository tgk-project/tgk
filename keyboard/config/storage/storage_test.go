package storage_test

import (
	"errors"
	"testing"

	"github.com/tgk-project/tgk/keyboard/config"
	"github.com/tgk-project/tgk/keyboard/config/storage"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

const (
	testMagic    = 0x54474b32
	testSchema   = 1
	testLayoutID = 0x99aa55cc
)

func factory() config.Snapshot {
	return config.Snapshot{Bindings: []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: 4}}}
}

func snapshot(usage uint32) config.Snapshot {
	return config.Snapshot{Bindings: []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: usage}}}
}

func newStore(t *testing.T, device storage.SlotDevice, schema uint16, layout uint32, debounce uint32) *storage.Store {
	return newStoreWithMagic(t, device, testMagic, schema, layout, debounce)
}

func newStoreWithMagic(t *testing.T, device storage.SlotDevice, magic uint32, schema uint16, layout uint32, debounce uint32) *storage.Store {
	t.Helper()
	value, err := storage.New(device, storage.Options{
		Magic:         magic,
		SchemaVersion: schema,
		LayoutID:      layout,
		Factory:       factory(),
		DebounceTicks: debounce,
	})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func newMemoryStore(t *testing.T) (*storage.Store, *storage.MemoryDevice) {
	t.Helper()
	device, err := storage.NewMemoryDevice(128)
	if err != nil {
		t.Fatal(err)
	}
	return newStore(t, device, testSchema, testLayoutID, 0), device
}

func assertUsage(t *testing.T, got config.Snapshot, want uint32) {
	t.Helper()
	if len(got.Bindings) != 1 || got.Bindings[0].Param1 != want {
		t.Fatalf("bindings = %#v, want one key usage %d", got.Bindings, want)
	}
}

func TestLoadUsesFactoryWhenUserSlotsAreEmpty(t *testing.T) {
	store, _ := newMemoryStore(t)
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 4)
	if store.Source() != storage.SourceFactory {
		t.Fatalf("Source() = %v, want Factory", store.Source())
	}
}

func TestNewestCorruptSlotFallsBackToPreviousGeneration(t *testing.T) {
	store, device := newMemoryStore(t)
	if err := store.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(snapshot(6)); err != nil {
		t.Fatal(err)
	}
	// Keep the decoded binding structurally valid while changing its key usage,
	// so this assertion specifically exercises CRC rejection rather than the
	// independent binding validation path.
	if err := device.Corrupt(1, storage.HeaderSize+2, 7); err != nil {
		t.Fatal(err)
	}

	restarted := newStore(t, device, testSchema, testLayoutID, 0)
	got, err := restarted.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 5)
	if restarted.Source() != storage.SourceSlotA {
		t.Fatalf("Source() = %v, want SlotA", restarted.Source())
	}
}

func TestInvalidSchemaAndLayoutSafelyUseFactory(t *testing.T) {
	store, device := newMemoryStore(t)
	if err := store.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name   string
		schema uint16
		layout uint32
	}{
		{name: "schema", schema: testSchema + 1, layout: testLayoutID},
		{name: "layout", schema: testSchema, layout: testLayoutID + 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			restarted := newStore(t, device, tc.schema, tc.layout, 0)
			got, err := restarted.Load()
			if err != nil {
				t.Fatal(err)
			}
			assertUsage(t, got, 4)
			if restarted.Source() != storage.SourceFactory {
				t.Fatalf("Source() = %v, want Factory", restarted.Source())
			}
			if !restarted.MigrationRequired() {
				t.Fatal("MigrationRequired() = false, want true")
			}
			if err := restarted.Save(snapshot(6)); !errors.Is(err, storage.ErrMigrationRequired) {
				t.Fatalf("Save() error = %v, want ErrMigrationRequired", err)
			}

			compatible := newStore(t, device, testSchema, testLayoutID, 0)
			got, err = compatible.Load()
			if err != nil {
				t.Fatal(err)
			}
			assertUsage(t, got, 5)
			if compatible.Source() != storage.SourceSlotA {
				t.Fatalf("compatible Source() = %v, want SlotA", compatible.Source())
			}
		})
	}
}

func TestIncompatibleLoadRequiresExplicitResetBeforeCommit(t *testing.T) {
	store, device := newMemoryStore(t)
	if err := store.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}

	updatedFirmware := newStore(t, device, testSchema+1, testLayoutID, 0)
	got, err := updatedFirmware.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 4)
	if updatedFirmware.Source() != storage.SourceFactory {
		t.Fatalf("incompatible Source() = %v, want Factory", updatedFirmware.Source())
	}
	if !updatedFirmware.MigrationRequired() {
		t.Fatal("MigrationRequired() = false, want true")
	}
	if err := updatedFirmware.Save(snapshot(6)); !errors.Is(err, storage.ErrMigrationRequired) {
		t.Fatalf("Save() error = %v, want ErrMigrationRequired", err)
	}

	previousFirmware := newStore(t, device, testSchema, testLayoutID, 0)
	got, err = previousFirmware.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 5)
	if previousFirmware.Source() != storage.SourceSlotA {
		t.Fatalf("compatible Source() = %v, want SlotA", previousFirmware.Source())
	}

	if err := updatedFirmware.Reset(); err != nil {
		t.Fatal(err)
	}
	if updatedFirmware.MigrationRequired() {
		t.Fatal("MigrationRequired() = true after Reset")
	}
	if err := updatedFirmware.Save(snapshot(6)); err != nil {
		t.Fatal(err)
	}
	restartedUpdatedFirmware := newStore(t, device, testSchema+1, testLayoutID, 0)
	got, err = restartedUpdatedFirmware.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 6)
}

func TestChangedFormatMagicRequiresExplicitResetBeforeCommit(t *testing.T) {
	_, device := newMemoryStore(t)
	previousFirmware := newStoreWithMagic(t, device, testMagic-1, testSchema, testLayoutID, 0)
	if err := previousFirmware.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}

	updatedFirmware := newStore(t, device, testSchema, testLayoutID, 0)
	got, err := updatedFirmware.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 4)
	if !updatedFirmware.MigrationRequired() {
		t.Fatal("MigrationRequired() = false, want true")
	}
	if err := updatedFirmware.Save(snapshot(6)); !errors.Is(err, storage.ErrMigrationRequired) {
		t.Fatalf("Save() error = %v, want ErrMigrationRequired", err)
	}

	stillPrevious := newStoreWithMagic(t, device, testMagic-1, testSchema, testLayoutID, 0)
	got, err = stillPrevious.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 5)
}

func TestCorruptReservedHeaderByteFallsBackToFactory(t *testing.T) {
	store, device := newMemoryStore(t)
	if err := store.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}
	if err := device.Corrupt(0, 6, 0x01); err != nil {
		t.Fatal(err)
	}

	restarted := newStore(t, device, testSchema, testLayoutID, 0)
	got, err := restarted.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 4)
}

func TestInterruptedWritePreservesPreviousSlot(t *testing.T) {
	base, err := storage.NewMemoryDevice(128)
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t, base, testSchema, testLayoutID, 0)
	if err := store.Save(snapshot(5)); err != nil {
		t.Fatal(err)
	}

	faulty := &tornDevice{SlotDevice: base, failAfter: 8}
	store = newStore(t, faulty, testSchema, testLayoutID, 0)
	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	if err := store.Update(snapshot(6), 0); err != nil {
		t.Fatal(err)
	}
	if err := store.Commit(); !errors.Is(err, storage.ErrWriteFailed) {
		t.Fatalf("Commit() error = %v, want ErrWriteFailed", err)
	}

	restarted := newStore(t, base, testSchema, testLayoutID, 0)
	got, err := restarted.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 5)
}

func TestUpdateIsImmediateAndDebouncedCommitIsExplicit(t *testing.T) {
	device, err := storage.NewMemoryDevice(128)
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t, device, testSchema, testLayoutID, 10)
	if err := store.Update(snapshot(7), 100); err != nil {
		t.Fatal(err)
	}
	assertUsage(t, store.Current(), 7)
	if !store.Dirty() {
		t.Fatal("Dirty() = false, want true")
	}
	if committed, err := store.CommitDue(109); err != nil || committed {
		t.Fatalf("CommitDue(before deadline) = %t, %v", committed, err)
	}
	if committed, err := store.CommitDue(110); err != nil || !committed {
		t.Fatalf("CommitDue(deadline) = %t, %v", committed, err)
	}
	if store.Dirty() {
		t.Fatal("Dirty() = true after commit")
	}

	restarted := newStore(t, device, testSchema, testLayoutID, 10)
	got, err := restarted.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 7)
}

func TestResetErasesOnlyUserConfigAndRestoresFactory(t *testing.T) {
	store, device := newMemoryStore(t)
	if err := store.Save(snapshot(8)); err != nil {
		t.Fatal(err)
	}
	if err := store.Reset(); err != nil {
		t.Fatal(err)
	}
	assertUsage(t, store.Current(), 4)

	restarted := newStore(t, device, testSchema, testLayoutID, 0)
	got, err := restarted.Load()
	if err != nil {
		t.Fatal(err)
	}
	assertUsage(t, got, 4)
}

func TestUpdateRejectsMalformedBinding(t *testing.T) {
	store, _ := newMemoryStore(t)
	for _, malformed := range []config.Snapshot{
		{Bindings: []keymap.Binding{{Behavior: keymap.BehaviorKey}}},
		{Bindings: []keymap.Binding{{Behavior: keymap.BehaviorID(99)}}},
	} {
		if err := store.Update(malformed, 0); !errors.Is(err, storage.ErrInvalidBinding) {
			t.Fatalf("Update(%#v) error = %v, want ErrInvalidBinding", malformed, err)
		}
	}
}

func TestNewRejectsFactoryThatDoesNotFitSlot(t *testing.T) {
	device, err := storage.NewMemoryDevice(storage.HeaderSize + 10)
	if err != nil {
		t.Fatal(err)
	}
	_, err = storage.New(device, storage.Options{
		Magic:         testMagic,
		SchemaVersion: testSchema,
		LayoutID:      testLayoutID,
		Factory: config.Snapshot{Bindings: []keymap.Binding{
			{Behavior: keymap.BehaviorKey, Param1: 4},
			{Behavior: keymap.BehaviorKey, Param1: 5},
		}},
	})
	if !errors.Is(err, storage.ErrSnapshotLarge) {
		t.Fatalf("New() error = %v, want ErrSnapshotLarge", err)
	}
}

type tornDevice struct {
	storage.SlotDevice
	failAfter int
}

func (d *tornDevice) WriteSlot(slot int, src []byte) error {
	if len(src) > d.failAfter {
		if err := d.SlotDevice.WriteSlot(slot, src[:d.failAfter]); err != nil {
			return err
		}
		return errors.New("simulated power loss")
	}
	return d.SlotDevice.WriteSlot(slot, src)
}

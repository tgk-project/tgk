package via_test

import (
	"errors"
	"testing"

	"github.com/tgk-project/tgk/keyboard/config"
	"github.com/tgk-project/tgk/keyboard/config/service"
	"github.com/tgk-project/tgk/keyboard/config/via"
	"github.com/tgk-project/tgk/keyboard/definition"
	"github.com/tgk-project/tgk/keyboard/keymap"
)

type memoryStore struct{ value config.Snapshot }

func (s *memoryStore) Load() (config.Snapshot, error)   { return s.value, nil }
func (s *memoryStore) Save(value config.Snapshot) error { s.value = value; return nil }
func (s *memoryStore) Reset() error                     { s.value = config.Snapshot{}; return nil }

type sink struct{ keymap.Keymap }

func (s *sink) ReplaceKeymap(value keymap.Keymap) error { s.Keymap = value; return nil }

type rawTransport struct {
	requests  [][]byte
	responses [][]byte
}

func (r *rawTransport) ReceiveRaw(dst []byte) (int, bool) {
	if len(r.requests) == 0 {
		return 0, false
	}
	value := r.requests[0]
	r.requests = r.requests[1:]
	copy(dst, value)
	return len(value), true
}

func (r *rawTransport) SendRaw(value []byte) error {
	r.responses = append(r.responses, append([]byte(nil), value...))
	return nil
}

func newAdapter(t *testing.T) (*via.Adapter, *service.Service) {
	t.Helper()
	store := &memoryStore{}
	target := &sink{}
	configService, err := service.New(service.Options{
		Identity: definition.Identity{VendorID: 0x1209, ProductID: 1, ProductName: "test", LayoutID: 1},
		Layers:   1, Positions: 2, UnlockTicks: 10,
		Factory: []keymap.Binding{{Behavior: keymap.BehaviorKey, Param1: 4}, {Behavior: keymap.BehaviorModifier, Param1: 1}},
		Store:   store, Sink: target,
	})
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := via.New(configService, via.Matrix{Rows: 1, Columns: 2})
	if err != nil {
		t.Fatal(err)
	}
	return adapter, configService
}

func packet(command byte) []byte {
	value := make([]byte, via.ReportSize)
	value[0] = command
	return value
}

func TestProtocolReadAndKeymapBuffer(t *testing.T) {
	adapter, _ := newAdapter(t)
	response, err := adapter.Handle(packet(via.GetProtocolVersion), 0)
	if err != nil || response[1] != 0 || response[2] != via.ProtocolVersion {
		t.Fatalf("protocol response=%#v error=%v", response, err)
	}
	response, err = adapter.Handle(packet(via.GetLayerCount), 0)
	if err != nil || response[1] != 1 {
		t.Fatalf("layer response=%#v error=%v", response, err)
	}
	request := packet(via.GetKeymapBuffer)
	request[3] = 4
	response, err = adapter.Handle(request, 0)
	if err != nil || response[4] != 0 || response[5] != 4 || response[6] != 0 || response[7] != 0xe0 {
		t.Fatalf("buffer response=%#v error=%v", response, err)
	}
}

func TestWritesAreLockedThenImmediateAndCommitSeparately(t *testing.T) {
	adapter, configService := newAdapter(t)
	request := packet(via.SetKeycode)
	request[1], request[2], request[3] = 0, 0, 1
	request[5] = 5
	if _, err := adapter.Handle(request, 100); !errors.Is(err, service.ErrLocked) {
		t.Fatalf("locked write error = %v", err)
	}
	configService.Unlock(100)
	response, err := adapter.Handle(request, 100)
	if err != nil || response[5] != 5 || !configService.Dirty() {
		t.Fatalf("write response=%#v error=%v dirty=%t", response, err, configService.Dirty())
	}
	if got, _ := configService.Binding(0, 1); got.Param1 != 5 {
		t.Fatalf("updated binding=%#v", got)
	}
	if err := adapter.Commit(); err != nil || configService.Dirty() {
		t.Fatalf("commit error=%v dirty=%t", err, configService.Dirty())
	}
}

func TestBulkValidationAndMalformedInputNeverPanic(t *testing.T) {
	adapter, configService := newAdapter(t)
	configService.Unlock(0)
	request := packet(via.SetKeymapBuffer)
	request[3] = 3 // a keycode buffer must have complete 16-bit keycodes
	if _, err := adapter.Handle(request, 0); !errors.Is(err, service.ErrInvalidBuffer) {
		t.Fatalf("odd bulk length error = %v", err)
	}
	if _, err := adapter.Handle([]byte{via.GetKeycode}, 0); !errors.Is(err, via.ErrInvalidPacket) {
		t.Fatalf("short packet error = %v", err)
	}
	for command := 0; command < 256; command++ {
		request := packet(byte(command))
		_, _ = adapter.Handle(request, 0)
	}
}

func TestRunnerHandlesAtMostOneQueuedPacketPerEventLoopTurn(t *testing.T) {
	adapter, _ := newAdapter(t)
	raw := &rawTransport{requests: [][]byte{packet(via.GetProtocolVersion), packet(via.GetLayerCount)}}
	runner, err := via.NewRunner(adapter, raw)
	if err != nil {
		t.Fatal(err)
	}
	for want := 1; want <= 2; want++ {
		processed, err := runner.ProcessNext(0)
		if err != nil || !processed || len(raw.responses) != want {
			t.Fatalf("turn %d processed=%t error=%v responses=%d", want, processed, err, len(raw.responses))
		}
	}
	if processed, err := runner.ProcessNext(0); err != nil || processed {
		t.Fatalf("empty ProcessNext = %t, %v", processed, err)
	}
}

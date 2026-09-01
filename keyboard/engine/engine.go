package engine

import (
	"errors"

	"github.com/tgk-project/tgk/keyboard/event"
	"github.com/tgk-project/tgk/keyboard/keymap"
	"github.com/tgk-project/tgk/keyboard/report"
	"github.com/tgk-project/tgk/keyboard/transport"
)

var (
	ErrInvalidPressedCapacity = errors.New("maximum pressed positions must be positive")
	ErrPressedCapacity        = errors.New("pressed position capacity exceeded")
	ErrQueueOverflow          = errors.New("event queue overflow; keyboard state released")
	ErrInvalidEventBatch      = errors.New("event source returned more events than the caller buffer")
	ErrIncompatibleKeymap     = errors.New("replacement keymap has incompatible dimensions")
)

// Options bounds all engine-owned memory. MaxPressed must cover the maximum
// simultaneous physical positions from all local and split sources.
type Options struct {
	QueueCapacity int
	MaxPressed    int
}

// Engine is owned by one Primary event loop. Callbacks may only use Enqueue;
// Run, Reset, and Report must be called by the owning loop.
type Engine struct {
	keymap keymap.Keymap
	host   transport.HostTransport
	queue  *Queue

	pressed      []pressedPosition
	pressedCount int

	baseLayer uint8
	toggled   []bool
	momentary []uint8

	modifierCounts [8]uint8
	keyCounts      [256]uint8
	consumers      [6]consumerSlot
}

type physicalPosition struct {
	source   event.Source
	position uint16
}

type pressedPosition struct {
	physical physicalPosition
	binding  keymap.Binding
}

type consumerSlot struct {
	usage uint16
	count uint8
}

// New creates an engine for one scanner source. Split keyboards should use
// NewWithOptions and provide a MaxPressed value that covers both halves.
func New(km keymap.Keymap, host transport.HostTransport, queueCapacity int) (*Engine, error) {
	return NewWithOptions(km, host, Options{
		QueueCapacity: queueCapacity,
		MaxPressed:    int(km.PositionCount()),
	})
}

// NewWithOptions preallocates state so normal press/release processing does
// not need to allocate.
func NewWithOptions(km keymap.Keymap, host transport.HostTransport, options Options) (*Engine, error) {
	if options.MaxPressed <= 0 {
		return nil, ErrInvalidPressedCapacity
	}
	queue, err := NewQueue(options.QueueCapacity)
	if err != nil {
		return nil, err
	}
	layers := int(km.LayerCount())
	if layers == 0 {
		return nil, keymap.ErrNoLayers
	}
	return &Engine{
		keymap:    km,
		host:      host,
		queue:     queue,
		pressed:   make([]pressedPosition, options.MaxPressed),
		toggled:   make([]bool, layers),
		momentary: make([]uint8, layers),
	}, nil
}

// Enqueue is safe to call from a scanner, split, USB, or BLE callback. It does
// no keymap or persistence work and does not allocate.
func (e *Engine) Enqueue(value event.KeyEvent) bool {
	return e.queue.Push(value)
}

// PollScanner receives a bounded batch from a local scanner and enqueues it.
func (e *Engine) PollScanner(scanner event.Scanner, now uint32, dst []event.KeyEvent) (int, error) {
	values, err := scanner.Scan(now, dst[:0])
	if err != nil {
		return 0, err
	}
	return e.enqueueBatch(values, cap(dst))
}

// ReceiveSplit receives a bounded batch from a split transport and enqueues it.
func (e *Engine) ReceiveSplit(split transport.SplitTransport, dst []event.KeyEvent) (int, error) {
	values, err := split.ReceiveEvents(dst[:0])
	if err != nil {
		return 0, err
	}
	return e.enqueueBatch(values, cap(dst))
}

func (e *Engine) enqueueBatch(values []event.KeyEvent, capacity int) (int, error) {
	if len(values) > capacity {
		return 0, ErrInvalidEventBatch
	}
	for i, value := range values {
		if !e.Enqueue(value) {
			return i, ErrQueueOverflow
		}
	}
	return len(values), nil
}

// Run processes up to limit queued events. A non-positive limit drains all
// currently queued events. It is the only method that mutates keyboard state.
func (e *Engine) Run(limit int) (int, error) {
	if e.queue.recoverOverflow() {
		if err := e.resetState(); err != nil {
			return 0, err
		}
		return 0, ErrQueueOverflow
	}

	processed := 0
	for limit <= 0 || processed < limit {
		value, ok := e.queue.pop()
		if !ok {
			break
		}
		if err := e.process(value); err != nil {
			return processed, err
		}
		processed++
	}
	return processed, nil
}

// Reset safely releases all keyboard state, for example when a Host or split
// transport disconnects.
func (e *Engine) Reset() error {
	return e.resetState()
}

// ReleaseSource safely releases every pressed position from one physical
// producer while preserving state owned by other scanner or split sources.
// It must be called by the Primary event-loop owner, not a transport callback.
func (e *Engine) ReleaseSource(source event.Source) error {
	e.queue.dropSource(source)

	changed := false
	for index := 0; index < e.pressedCount; {
		if e.pressed[index].physical.source != source {
			index++
			continue
		}
		binding := e.pressed[index].binding
		e.pressedCount--
		e.pressed[index] = e.pressed[e.pressedCount]
		e.pressed[e.pressedCount] = pressedPosition{}
		if e.releaseBinding(binding) {
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return e.sendReport()
}

// Report returns a value copy of the current Primary-owned Host state.
func (e *Engine) Report() report.KeyboardReport {
	return e.currentReport()
}

// ReplaceKeymap swaps the Primary-owned keymap from the event-loop owner.
// Pressed bindings are resolved at press time, so state is first released to
// prevent an old press from being released through a newly assigned binding.
// The physical position and layer dimensions stay fixed for a running
// keyboard definition.
func (e *Engine) ReplaceKeymap(value keymap.Keymap) error {
	if value.PositionCount() != e.keymap.PositionCount() || value.LayerCount() != e.keymap.LayerCount() {
		return ErrIncompatibleKeymap
	}
	if err := e.resetState(); err != nil {
		return err
	}
	e.keymap = value
	return nil
}

// QueueStats reports callback-side input loss for observability.
func (e *Engine) QueueStats() QueueStats {
	return e.queue.Stats()
}

func (e *Engine) process(value event.KeyEvent) error {
	physical := physicalPosition{source: value.Source, position: value.Position}
	if value.Pressed {
		if e.findPressed(physical) >= 0 {
			return nil
		}
		if e.pressedCount == len(e.pressed) {
			if err := e.resetState(); err != nil {
				return err
			}
			return ErrPressedCapacity
		}
		binding := e.resolve(value.Position)
		e.pressed[e.pressedCount] = pressedPosition{physical: physical, binding: binding}
		e.pressedCount++
		return e.applyPress(binding)
	}

	index := e.findPressed(physical)
	if index < 0 {
		return nil
	}
	binding := e.pressed[index].binding
	e.pressedCount--
	e.pressed[index] = e.pressed[e.pressedCount]
	e.pressed[e.pressedCount] = pressedPosition{}
	return e.applyRelease(binding)
}

func (e *Engine) findPressed(physical physicalPosition) int {
	for i := 0; i < e.pressedCount; i++ {
		if e.pressed[i].physical == physical {
			return i
		}
	}
	return -1
}

func (e *Engine) resolve(position uint16) keymap.Binding {
	for layer := int(e.keymap.LayerCount()) - 1; layer >= 0; layer-- {
		if !e.layerActive(uint8(layer)) {
			continue
		}
		binding, ok := e.keymap.BindingAt(uint8(layer), position)
		if !ok || binding.Behavior == keymap.BehaviorTransparent {
			continue
		}
		return binding
	}
	return keymap.Binding{Behavior: keymap.BehaviorNone}
}

func (e *Engine) layerActive(layer uint8) bool {
	return e.baseLayer == layer || e.toggled[layer] || e.momentary[layer] > 0
}

func (e *Engine) applyPress(binding keymap.Binding) error {
	switch binding.Behavior {
	case keymap.BehaviorKey:
		e.addKey(uint8(binding.Param1))
		return e.sendReport()
	case keymap.BehaviorModifier:
		e.addModifiers(uint8(binding.Param1))
		return e.sendReport()
	case keymap.BehaviorConsumer:
		e.addConsumer(uint16(binding.Param1))
		return e.sendReport()
	case keymap.BehaviorMomentaryLayer:
		layer := uint8(binding.Param1)
		if e.momentary[layer] < ^uint8(0) {
			e.momentary[layer]++
		}
	case keymap.BehaviorToLayer:
		e.baseLayer = uint8(binding.Param1)
		clear(e.toggled)
		clear(e.momentary)
	case keymap.BehaviorToggleLayer:
		layer := uint8(binding.Param1)
		e.toggled[layer] = !e.toggled[layer]
	}
	return nil
}

func (e *Engine) applyRelease(binding keymap.Binding) error {
	if !e.releaseBinding(binding) {
		return nil
	}
	return e.sendReport()
}

func (e *Engine) releaseBinding(binding keymap.Binding) bool {
	switch binding.Behavior {
	case keymap.BehaviorKey:
		e.removeKey(uint8(binding.Param1))
		return true
	case keymap.BehaviorModifier:
		e.removeModifiers(uint8(binding.Param1))
		return true
	case keymap.BehaviorConsumer:
		e.removeConsumer(uint16(binding.Param1))
		return true
	case keymap.BehaviorMomentaryLayer:
		layer := uint8(binding.Param1)
		if e.momentary[layer] > 0 {
			e.momentary[layer]--
		}
	}
	return false
}

func (e *Engine) addKey(usage uint8) {
	if e.keyCounts[usage] < ^uint8(0) {
		e.keyCounts[usage]++
	}
}

func (e *Engine) removeKey(usage uint8) {
	if e.keyCounts[usage] == 0 {
		return
	}
	e.keyCounts[usage]--
}

func (e *Engine) addModifiers(mask uint8) {
	for bit := uint8(0); bit < 8; bit++ {
		if mask&(1<<bit) != 0 && e.modifierCounts[bit] < ^uint8(0) {
			e.modifierCounts[bit]++
		}
	}
}

func (e *Engine) removeModifiers(mask uint8) {
	for bit := uint8(0); bit < 8; bit++ {
		if mask&(1<<bit) != 0 && e.modifierCounts[bit] > 0 {
			e.modifierCounts[bit]--
		}
	}
}

func (e *Engine) addConsumer(usage uint16) {
	for i := range e.consumers {
		if e.consumers[i].usage == usage {
			if e.consumers[i].count < ^uint8(0) {
				e.consumers[i].count++
			}
			return
		}
	}
	for i := range e.consumers {
		if e.consumers[i].count == 0 {
			e.consumers[i] = consumerSlot{usage: usage, count: 1}
			return
		}
	}
}

func (e *Engine) removeConsumer(usage uint16) {
	for i := range e.consumers {
		if e.consumers[i].usage != usage || e.consumers[i].count == 0 {
			continue
		}
		e.consumers[i].count--
		if e.consumers[i].count == 0 {
			e.consumers[i] = consumerSlot{}
		}
		return
	}
}

func (e *Engine) currentReport() report.KeyboardReport {
	var modifiers uint8
	for bit := uint8(0); bit < 8; bit++ {
		if e.modifierCounts[bit] > 0 {
			modifiers |= 1 << bit
		}
	}
	var consumer uint16
	for _, slot := range e.consumers {
		if slot.count > 0 {
			consumer = slot.usage
			break
		}
	}
	return report.KeyboardReport{Modifiers: modifiers, Keys: e.keyboardKeys(), Consumer: consumer}
}

func (e *Engine) keyboardKeys() [6]uint8 {
	var keys [6]uint8
	index := 0
	for usage := 1; usage <= 0xff; usage++ {
		if e.keyCounts[usage] == 0 {
			continue
		}
		if index == len(keys) {
			return [6]uint8{
				report.ErrorRollOver,
				report.ErrorRollOver,
				report.ErrorRollOver,
				report.ErrorRollOver,
				report.ErrorRollOver,
				report.ErrorRollOver,
			}
		}
		keys[index] = uint8(usage)
		index++
	}
	return keys
}

func (e *Engine) sendReport() error {
	if e.host == nil || !e.host.Ready() {
		return nil
	}
	return e.host.SendKeyboardReport(e.currentReport())
}

func (e *Engine) resetState() error {
	e.pressedCount = 0
	clear(e.pressed)
	e.baseLayer = 0
	clear(e.toggled)
	clear(e.momentary)
	e.modifierCounts = [8]uint8{}
	e.keyCounts = [256]uint8{}
	e.consumers = [6]consumerSlot{}
	if e.host == nil {
		return nil
	}
	return e.host.ReleaseAll()
}

// Package config declares persistence boundaries used by future configuration
// services. It contains no flash or board implementation.
package config

import "github.com/tgk-project/tgk/keyboard/keymap"

// Header identifies a serialized user configuration snapshot.
type Header struct {
	Magic         uint32
	SchemaVersion uint16
	LayoutID      uint32
	Generation    uint32
	PayloadLength uint32
	CRC32         uint32
}

// Snapshot stores serializable bindings, rather than Go pointers or behavior
// implementation names.
type Snapshot struct {
	Header   Header
	Bindings []keymap.Binding
}

// ConfigStore is the persistence boundary. Reset affects only User Config;
// callers must keep bond removal as a separate operation.
type ConfigStore interface {
	Load() (Snapshot, error)
	Save(Snapshot) error
	Reset() error
}

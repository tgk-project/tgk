//go:build xiao_ble

package xiaoble

import (
	"machine"

	"github.com/tgk-project/tgk/keyboard/config/storage"
)

// NewConfigSlots adapts the XIAO BLE nRF52840 internal Flash to two consecutive
// User Config slots. firstEraseBlock is a board-composition decision: it must
// refer to two blocks reserved outside the firmware image and Factory Config.
func NewConfigSlots(firstEraseBlock int64) (*storage.FlashSlots, error) {
	return storage.NewFlashSlots(machine.Flash, firstEraseBlock)
}

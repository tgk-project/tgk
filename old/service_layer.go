package tgk

import (
	"fmt"
)

// LayerService はレイヤーサービスのインターフェースです
type LayerService interface {
	// Init はレイヤーサービスを初期化します
	Init(config KeyboardConfig) error
	// GetCurrentLayer は現在のレイヤーを返します
	GetCurrentLayer() int
	// SetLayer は指定されたレイヤーに切り替えます
	SetLayer(layer int) error
	// GetKeymap は指定されたレイヤーとキーのキーマップを返します
	GetKeymap(layer, key int) (int, error)
	// IsLayerKey は指定されたキーがレイヤーキーかどうかを返します
	IsLayerKey(key int) bool
	// GetLayerName は現在のレイヤーの名前を返します
	GetLayerName() string
	// LayerTask はレイヤーキーが押された時の処理を行います
	LayerTask(layer int, isHold bool) error
}

// NewLayerService はレイヤーサービスを作成します
func NewLayerService(keymapRepo KeymapRepository, configRepo ConfigRepository) LayerService {
	return &layerService{
		currentLayer: KC_BASE,
		keymapRepo:   keymapRepo,
		configRepo:   configRepo,
	}
}

// layerService はレイヤーサービスの実装です
type layerService struct {
	currentLayer int
	keymapRepo   KeymapRepository
	configRepo   ConfigRepository
}

func (s *layerService) Init(config KeyboardConfig) error {
	// TODO: 実装
	return nil
}

func (s *layerService) GetCurrentLayer() int {
	return s.currentLayer
}

func (s *layerService) SetLayer(layer int) error {
	// if layer < 0 || layer >= len(s.keymapRepo.GetKeymaps()) {
	// 	return fmt.Errorf("invalid layer: %d", layer)
	// }
	// s.currentLayer = layer
	return nil
}

func (s *layerService) GetKeymap(layer, key int) (int, error) {
	// if layer < 0 || layer >= len(s.keymapRepo.GetKeymaps()) {
	// 	return 0, fmt.Errorf("invalid layer: %d", layer)
	// }
	// if key < 0 || key >= len(s.keymapRepo.GetKeymaps()[layer]) {
	// 	return 0, fmt.Errorf("invalid key: %d", key)
	// }
	// return s.keymapRepo.GetKeymap(layer, key), nil
	return 0, nil
}

func (s *layerService) IsLayerKey(key int) bool {
	switch key {
	case KC_BASE, KC_LOWER, KC_RAISE, KC_ADJUST:
		return true
	}
	return false
}

func (s *layerService) GetLayerName() string {
	switch s.currentLayer {
	case KC_BASE:
		return "BASE"
	case KC_LOWER:
		return "LOWER"
	case KC_RAISE:
		return "RAISE"
	case KC_ADJUST:
		return "ADJUST"
	}
	return ""
}

func (s *layerService) LayerTask(layer int, isHold bool) error {
	nowLayer := s.GetCurrentLayer()
	switch layer {
	case KC_LOWER:
		if isHold {
			if nowLayer == KC_RAISE {
				return s.SetLayer(KC_ADJUST)
			} else {
				return s.SetLayer(KC_LOWER)
			}
		} else {
			if nowLayer == KC_ADJUST {
				return s.SetLayer(KC_RAISE)
			} else {
				return s.SetLayer(KC_BASE)
			}
		}
	case KC_RAISE:
		if isHold {
			if nowLayer == KC_LOWER {
				return s.SetLayer(KC_ADJUST)
			} else {
				return s.SetLayer(KC_RAISE)
			}
		} else {
			if nowLayer == KC_ADJUST {
				return s.SetLayer(KC_LOWER)
			} else {
				return s.SetLayer(KC_BASE)
			}
		}
	}
	return fmt.Errorf("invalid layer task: layer=%d, isHold=%v", layer, isHold)
}

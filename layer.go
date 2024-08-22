package tgk

type Layer struct {
	nowLayer uint8
}

func NewLayer() Layer {
	return Layer{
		nowLayer: KC_BASE,
	}
}

func (l Layer) SetLayer(layer uint8) {
	l.nowLayer = layer
}

func (l Layer) GetLayer() uint8 {
	return l.nowLayer
}

func (l Layer) IsLayerKey(key uint8) bool {
	switch key {
	case KC_BASE, KC_LOWER, KC_RAISE, KC_ADJUST:
		return true
	}
	return false
}

func (l Layer) GetLayerName() string {
	switch l.nowLayer {
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

func (l Layer) LayerTask(layer uint8, isHold bool) {
	nowLayer := l.GetLayer()
	switch layer {
	case KC_LOWER:
		if isHold {
			if nowLayer == KC_RAISE {
				l.SetLayer(KC_ADJUST)
			} else {
				l.SetLayer(KC_LOWER)
			}
		} else {
			if nowLayer == KC_ADJUST {
				l.SetLayer(KC_RAISE)
			} else {
				l.SetLayer(KC_BASE)
			}
		}
	case KC_RAISE:
		if isHold {
			if nowLayer == KC_LOWER {
				l.SetLayer(KC_ADJUST)
			} else {
				l.SetLayer(KC_RAISE)
			}
		} else {
			if nowLayer == KC_ADJUST {
				l.SetLayer(KC_LOWER)
			} else {
				l.SetLayer(KC_BASE)
			}
		}
	}
}

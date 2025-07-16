package keyscan

import (
	"fmt"
	"machine"
)

// 循環参照を避けるため、tgk.KeyboardConfigを使用します

// KeyScanMX はMXキースキャンの実装です。
type KeyScanMX struct {
	rowPins    []machine.Pin
	colPins    []machine.Pin
	nowPushing [][]bool
	nowRelease [][]bool
	config     KeyboardConfig
}

// NewMXKeyScan はMXキースキャンを作成します
func NewMXKeyScan() *KeyScanMX {
	return &KeyScanMX{
		rowPins: []machine.Pin{machine.D12, machine.D11, machine.D10, machine.D9, machine.D6},
		colPins: []machine.Pin{machine.D5, machine.D4, machine.D3, machine.D2, machine.D1, machine.D0},
	}
}

func (s *KeyScanMX) Init(config KeyboardConfig) error {
	s.config = config

	// マトリックスの状態を初期化
	s.nowPushing = make([][]bool, config.Matrix.Rows)
	s.nowRelease = make([][]bool, config.Matrix.Rows)
	for i := range s.nowPushing {
		s.nowPushing[i] = make([]bool, config.Matrix.Cols)
		s.nowRelease[i] = make([]bool, config.Matrix.Cols)
	}

	// TODO: 実際のピン設定を実装
	// diodeDirection := config.KeyScanExtraConfigs.DiodeDirection
	// switch diodeDirection {
	// case "col2row":
	// 	// rowは出力、colは入力プルアップ
	// case "row2col":
	// 	// colは出力、rowは入力プルアップ
	// }

	return nil
}

func (s *KeyScanMX) Scan() bool {
	var isMatrixUpdate bool

	// TODO: 実際のキーマトリックススキャンを実装
	// 現在はダミー実装
	// diodeDirection := s.config.KeyScanExtraConfigs.DiodeDirection
	// switch diodeDirection {
	// case "col2row":
	// 	for row := 0; row < s.config.Matrix.Rows; row++ {
	// 		rowPin := s.rowPins[row]
	// 		rowPin.Low()
	// 		for col := 0; col < s.config.Matrix.Cols; col++ {
	// 			colPin := s.colPins[col]
	// 			if !colPin.Get() {
	// 				if !s.nowPushing[row][col] {
	// 					// keydown
	// 					isMatrixUpdate = true
	// 					s.nowRelease[row][col] = false
	// 				}
	// 				s.nowPushing[row][col] = true
	// 			} else {
	// 				if s.nowPushing[row][col] {
	// 					// keyup
	// 					isMatrixUpdate = true
	// 					s.nowRelease[row][col] = true
	// 				}
	// 				s.nowPushing[row][col] = false
	// 			}
	// 		}
	// 		rowPin.High()
	// 	}
	// case "row2col":
	// 	// 同様の実装
	// }

	return isMatrixUpdate
}

func (s *KeyScanMX) GetNowPushing() [][]bool {
	return s.nowPushing
}

func (s *KeyScanMX) GetNowRelease() [][]bool {
	return s.nowRelease
}

func (s *KeyScanMX) Print() string {

	paper := ""

	for row := 0; row < len(s.rowPins); row++ {
		for col := 0; col < len(s.colPins); col++ {
			paper += fmt.Sprintf("%d ", s.nowPushing[row][col])
		}
		paper += "\n"
	}

	return paper
}

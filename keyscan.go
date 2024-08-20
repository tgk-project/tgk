package tgk

import (
	"context"
	"machine"
	"strings"
)

/*
 * KeyScan interface
 */
type Keyscan interface {
	// 初期化
	Init(ctx context.Context, keyboard *Keyboard)
	// スキャン
	Scan(ctx context.Context, keyboard *Keyboard) bool
	// デバッグ用
	Print(ctx context.Context, keyboard *Keyboard) string
	// スキャン後の結果を取得する
	GetEventsAfterScan(keyboard *Keyboard) []KeyscanEvent
}

type KeyscanEvent struct {
	Row      uint8
	Col      uint8
	HoldFlag bool
}

/*
 * Default keyscan mx
 * 一般的なMXスイッチのキーマトリクスをスキャンする
 */
type MX struct {
	HoldMatrix [][]bool // 今押してる状態を保持する
	ScanMatrix [][]bool // スキャンのタイミングで押下しているかを保持する
}

func NewMXMatrix(row int, col int) *MX {
	// HoldMatrixとScanMatrixをn x mサイズのスライスで初期化
	holdMatrix := make([][]bool, row)
	scanMatrix := make([][]bool, row)
	for i := range holdMatrix {
		holdMatrix[i] = make([]bool, col)
		scanMatrix[i] = make([]bool, col)
	}
	return &MX{
		HoldMatrix: holdMatrix,
		ScanMatrix: scanMatrix,
	}
}

func (this MX) Init(ctx context.Context, keyboard *Keyboard) {

	diodeDirection := keyboard.KeyScanExtraConfigs["diode_direction"].(string)
	switch diodeDirection {
	case "col2row":
		// rowは出力、colは入力プルアップ
		for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {
			keyboard.MatrixPins["col"].Gpios[col].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
		for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
			keyboard.MatrixPins["row"].Gpios[row].Configure(machine.PinConfig{Mode: machine.PinOutput})
		}
	case "row2col":
		// colは出力、rowは入力プルアップ
		for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
			keyboard.MatrixPins["row"].Gpios[row].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
		for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {
			keyboard.MatrixPins["col"].Gpios[col].Configure(machine.PinConfig{Mode: machine.PinOutput})
		}
	}
	println("init mx complete")

}

func (this MX) Scan(ctx context.Context, keyboard *Keyboard) bool {

	diodeDirection := keyboard.KeyScanExtraConfigs["diode_direction"].(string)
	var isMatrixUpdate bool
	switch diodeDirection {
	case "col2row":
		for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
			rowPin := keyboard.MatrixPins["row"].Gpios[row]
			rowPin.Low()
			for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {
				colPin := keyboard.MatrixPins["col"].Gpios[col]
				if !colPin.Get() {
					if !this.ScanMatrix[row][col] {
						// keydown
						isMatrixUpdate = true
					}
					this.ScanMatrix[row][col] = true
				} else {
					if this.ScanMatrix[row][col] {
						// keyup
						isMatrixUpdate = true
					}
					this.ScanMatrix[row][col] = false
				}
			}
			rowPin.High()
		}
	case "row2col":
		for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {
			colPin := keyboard.MatrixPins["col"].Gpios[col]
			colPin.Low()
			for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
				rowPin := keyboard.MatrixPins["row"].Gpios[row]
				if !rowPin.Get() {
					if !this.ScanMatrix[row][col] {
						// keydown
						isMatrixUpdate = true
					}
					this.ScanMatrix[row][col] = true
				} else {
					if this.ScanMatrix[row][col] {
						// keyup
						isMatrixUpdate = true
					}
					this.ScanMatrix[row][col] = false
				}
			}
			colPin.High()
		}
	}

	return isMatrixUpdate
}

func (this MX) Print(ctx context.Context, keyboard *Keyboard) string {

	var result []string
	for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
		for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {
			if this.ScanMatrix[row][col] {
				result = append(result, "1 ")
			} else {
				result = append(result, "0 ")
			}
		}
		result = append(result, "\n")
	}
	result = append(result, "\n\n")
	return strings.Join(result, "")
}

func (this MX) GetEventsAfterScan(keyboard *Keyboard) []KeyscanEvent {
	res := []KeyscanEvent{}
	for row := 0; row < keyboard.MatrixPins["row"].Len; row++ {
		for col := 0; col < keyboard.MatrixPins["col"].Len; col++ {

			if this.ScanMatrix[row][col] && !this.HoldMatrix[row][col] {
				// keydown
				res = append(res, KeyscanEvent{
					Row:      uint8(row),
					Col:      uint8(col),
					HoldFlag: false,
				})
				this.HoldMatrix[row][col] = true
			}

			if !this.ScanMatrix[row][col] && this.HoldMatrix[row][col] {
				// keyup
				res = append(res, KeyscanEvent{
					Row:      uint8(row),
					Col:      uint8(col),
					HoldFlag: true,
				})
				this.HoldMatrix[row][col] = false
			}
		}
	}
	return res
}

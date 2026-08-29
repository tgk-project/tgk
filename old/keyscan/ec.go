package keyscan

import (
	"machine"
)

// KeyScanECはKeyScanMatrixを実装します。
type KeyScanEC struct { // implements tgk.KeyScanMatrix
	rowPins     []machine.Pin
	colChannels []int
	muxSelPins  []machine.Pin
	adc         machine.ADC
	adcGain     int
	nowPushing  [][]bool
	nowRelease  [][]bool
	config      KeyboardConfig

	dischargePin machine.Pin
	powerPin     machine.Pin
	muxEnPin     machine.Pin
}

// NewECKeyScan はECキースキャンを作成します
func NewECKeyScan() *KeyScanEC {
	return &KeyScanEC{
		rowPins:      []machine.Pin{machine.D12, machine.D11, machine.D10, machine.D9, machine.D6},
		colChannels:  []int{4, 6, 7, 5, 1, 0},
		muxSelPins:   []machine.Pin{machine.D1, machine.D0, machine.D2},
		adc:          machine.ADC{Pin: machine.A0},
		adcGain:      12,
		dischargePin: machine.D22,
		powerPin:     machine.D25,
		muxEnPin:     machine.D24,
	}
}

func (s *KeyScanEC) Init(config KeyboardConfig) error {
	s.config = config

	// マトリックスの状態を初期化
	s.nowPushing = make([][]bool, config.Matrix.Rows)
	s.nowRelease = make([][]bool, config.Matrix.Rows)
	for i := range s.nowPushing {
		s.nowPushing[i] = make([]bool, config.Matrix.Cols)
		s.nowRelease[i] = make([]bool, config.Matrix.Cols)
	}

	// 追加設定を動的に読み込む
	if config.KeyScanExtraConfigs.ADCGain > 0 {
		s.adcGain = config.KeyScanExtraConfigs.ADCGain
	}

	// colChannelの設定を読み込む
	if len(config.KeyScanExtraConfigs.ColChannel) > 0 {
		s.colChannels = make([]int, len(config.KeyScanExtraConfigs.ColChannel))
		for i, ch := range config.KeyScanExtraConfigs.ColChannel {
			// 文字列を整数に変換
			chInt, err := stringToInt(ch)
			if err != nil {
				return err
			}
			s.colChannels[i] = chInt
		}
	}

	// ピン設定を動的に読み込む
	if len(config.MatrixPins.Rows) > 0 {
		// 設定ファイルからピン設定を読み込む
		s.rowPins = make([]machine.Pin, len(config.MatrixPins.Rows))
		for i, pinName := range config.MatrixPins.Rows {
			pin, err := stringToPinWithError(pinName)
			if err != nil {
				return err
			}
			s.rowPins[i] = pin
		}
	}

	// ピンの初期化
	s.dischargePin.Low()
	s.dischargePin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	s.powerPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s.powerPin.High()
	s.muxEnPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	s.muxEnPin.Low()

	for _, pin := range s.rowPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		pin.Low()
	}

	for _, pin := range s.muxSelPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		pin.Low()
	}

	machine.InitADC()
	adcCfg := machine.ADCConfig{
		Reference: 3600,
	}
	s.adc.Configure(adcCfg)

	return nil
}

func (s *KeyScanEC) Scan() bool {
	var isMatrixUpdate bool

	for col := range s.colChannels {
		for row := range s.rowPins {
			s.dischargePin.Configure(machine.PinConfig{Mode: machine.PinOutput})

			ch := s.colChannels[col]
			s.muxSelPins[0].Set(ch&1 == 1)
			s.muxSelPins[1].Set(ch&2 == 2)
			s.muxSelPins[2].Set(ch&4 == 4)

			for _, pin := range s.rowPins {
				pin.Low()
			}

			s.dischargePin.Configure(machine.PinConfig{Mode: machine.PinInput})
			s.rowPins[row].High()

			val := s.analogRead()

			// 閾値を超えた場合はキーが押されている
			if val > uint16(s.config.KeyScanExtraConfigs.PushThreshold) {
				if !s.nowPushing[row][col] {
					// keydown
					isMatrixUpdate = true
					s.nowRelease[row][col] = false
				}
				s.nowPushing[row][col] = true
			} else if val < uint16(s.config.KeyScanExtraConfigs.ReleaseThreshold) {
				if s.nowPushing[row][col] {
					// keyup
					isMatrixUpdate = true
					s.nowRelease[row][col] = true
				}
				s.nowPushing[row][col] = false
			}
		}
	}

	return isMatrixUpdate
}

func (s *KeyScanEC) GetNowPushing() [][]bool {
	return s.nowPushing
}

func (s *KeyScanEC) GetNowRelease() [][]bool {
	return s.nowRelease
}

func (s *KeyScanEC) Print() string {
	paper := ""
	for row := 0; row < len(s.rowPins); row++ {
		for col := 0; col < len(s.colChannels); col++ {
			if s.nowPushing[row][col] {
				paper += "1 "
			} else {
				paper += "0 "
			}
		}
		paper += "\n"
	}
	return paper
}

func (s *KeyScanEC) analogRead() uint16 {
	return s.adc.Get() / uint16(s.adcGain)
}

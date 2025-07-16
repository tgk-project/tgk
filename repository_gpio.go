package tgk

import (
	"machine"
	"strings"
)

// GPIORepository はGPIOリポジトリのインターフェースです
type GPIORepository interface {
	// Init はGPIOリポジトリを初期化します
	Init() error
	// ConfigurePin はピンを設定します
	ConfigurePin(pin machine.Pin, mode machine.PinMode) error
	// SetPin はピンの状態を設定します
	SetPin(pin machine.Pin, high bool) error
	// GetPin はピンの状態を取得します
	GetPin(pin machine.Pin) (bool, error)
	// ToPin は文字列形式のピン名をmachine.Pin型に変換します
	ToPin(pin string) machine.Pin
	// ToGpios はマトリックスピンの設定をGpioConfigに変換します
	ToGpios(matrixPins map[string][]string) map[string]GpioConfig
}

// GpioConfig はGPIOの設定を表します
type GpioConfig struct {
	Len   int           `json:"len"`
	Gpios []machine.Pin `json:"gpios"`
}

// NewGPIORepository はGPIOリポジトリを作成します
func NewGPIORepository() GPIORepository {
	return &gpioRepository{}
}

// gpioRepository はGPIOリポジトリの実装です
type gpioRepository struct {
	configuredPins map[machine.Pin]machine.PinMode
}

func (r *gpioRepository) Init() error {
	r.configuredPins = make(map[machine.Pin]machine.PinMode)
	return nil
}

func (r *gpioRepository) ConfigurePin(pin machine.Pin, mode machine.PinMode) error {
	pin.Configure(machine.PinConfig{Mode: mode})
	r.configuredPins[pin] = mode
	return nil
}

func (r *gpioRepository) SetPin(pin machine.Pin, high bool) error {
	if high {
		pin.High()
	} else {
		pin.Low()
	}
	return nil
}

func (r *gpioRepository) GetPin(pin machine.Pin) (bool, error) {
	return pin.Get(), nil
}

func (r *gpioRepository) ToPin(pin string) machine.Pin {
	// HACK: 本来ならrefrectを使って動的にコードを生成するべきだが、マイコン上で動作させるということと、速度観点からパターンマッチを使う
	switch pin {
	case "D0":
		return machine.D0
	case "D1":
		return machine.D1
	case "D2":
		return machine.D2
	case "D3":
		return machine.D3
	case "D4":
		return machine.D4
	case "D5":
		return machine.D5
	case "D6":
		return machine.D6
	case "D7":
		return machine.D7
	case "D8":
		return machine.D8
	case "D9":
		return machine.D9
	case "D10":
		return machine.D10
	case "D11":
		return machine.D11
	case "D12":
		return machine.D12
	case "D13":
		return machine.D13
	case "D14":
		return machine.D14
	case "D15":
		return machine.D15
	case "D16":
		return machine.D16
	case "D17":
		return machine.D17
	case "D18":
		return machine.D18
	case "D19":
		return machine.D19
	case "D20":
		return machine.D20
	case "D21":
		return machine.D21
	case "D22":
		return machine.D22
	case "D23":
		return machine.D23
	case "D24":
		return machine.D24
	case "D25":
		return machine.D25
	case "D26":
		return machine.D26
	case "D27":
		return machine.D27
	case "D28":
		return machine.D28
	case "D29":
		return machine.D29
		// machineパッケージには30以上も存在するが、キーボードとして使われるマイコンに30以上のピンが存在するものは少ないので、ここで止めてる
	}
	panic("invalid pin")
}

func (r *gpioRepository) ToGpios(matrixPins map[string][]string) map[string]GpioConfig {
	gpios := map[string]GpioConfig{}
	for key, pins := range matrixPins {
		// rowsとかcolsになってるので、最後のsを削除
		trueKey := strings.TrimSuffix(key, "s")
		gpioPins := []machine.Pin{}
		for _, pin := range pins {
			gpioPins = append(gpioPins, r.ToPin(pin))
		}
		gpios[trueKey] = GpioConfig{
			Len:   len(pins),
			Gpios: gpioPins,
		}
	}
	return gpios
}

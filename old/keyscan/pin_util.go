package keyscan

import "machine"

// ErrInvalidPin は無効なピン名が指定された場合のエラーです
type ErrInvalidPin struct {
	Pin string
}

func (e ErrInvalidPin) Error() string {
	return "invalid pin: " + e.Pin
}

// stringToPin は文字列形式のピン名をmachine.Pin型に変換します
func stringToPin(pin string) machine.Pin {
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
	// panicではなくエラーを返す
	panic("invalid pin: " + pin)
}

// stringToPinWithError は文字列形式のピン名をmachine.Pin型に変換します
// 無効なピン名が指定された場合はエラーを返します
func stringToPinWithError(pin string) (machine.Pin, error) {
	// HACK: 本来ならrefrectを使って動的にコードを生成するべきだが、マイコン上で動作させるということと、速度観点からパターンマッチを使う
	switch pin {
	case "D0":
		return machine.D0, nil
	case "D1":
		return machine.D1, nil
	case "D2":
		return machine.D2, nil
	case "D3":
		return machine.D3, nil
	case "D4":
		return machine.D4, nil
	case "D5":
		return machine.D5, nil
	case "D6":
		return machine.D6, nil
	case "D7":
		return machine.D7, nil
	case "D8":
		return machine.D8, nil
	case "D9":
		return machine.D9, nil
	case "D10":
		return machine.D10, nil
	case "D11":
		return machine.D11, nil
	case "D12":
		return machine.D12, nil
	case "D13":
		return machine.D13, nil
	case "D14":
		return machine.D14, nil
	case "D15":
		return machine.D15, nil
	case "D16":
		return machine.D16, nil
	case "D17":
		return machine.D17, nil
	case "D18":
		return machine.D18, nil
	case "D19":
		return machine.D19, nil
	case "D20":
		return machine.D20, nil
	case "D21":
		return machine.D21, nil
	case "D22":
		return machine.D22, nil
	case "D23":
		return machine.D23, nil
	case "D24":
		return machine.D24, nil
	case "D25":
		return machine.D25, nil
	case "D26":
		return machine.D26, nil
	case "D27":
		return machine.D27, nil
	case "D28":
		return machine.D28, nil
	case "D29":
		return machine.D29, nil
		// machineパッケージには30以上も存在するが、キーボードとして使われるマイコンに30以上のピンが存在するものは少ないので、ここで止めてる
	}
	return 0, ErrInvalidPin{Pin: pin}
}

package tests

import (
	"machine"
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewGPIORepository(t *testing.T) {
	tests := []struct {
		name string
		want tgk.GPIORepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewGPIORepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGPIORepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gpioRepository_Init(t *testing.T) {
	tests := []struct {
		name    string
		r       tgk.GPIORepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Init(); (err != nil) != tt.wantErr {
				t.Errorf("gpioRepository.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gpioRepository_ConfigurePin(t *testing.T) {
	type args struct {
		pin  machine.Pin
		mode machine.PinMode
	}
	tests := []struct {
		name    string
		r       tgk.GPIORepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.ConfigurePin(tt.args.pin, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("gpioRepository.ConfigurePin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gpioRepository_SetPin(t *testing.T) {
	type args struct {
		pin  machine.Pin
		high bool
	}
	tests := []struct {
		name    string
		r       tgk.GPIORepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.SetPin(tt.args.pin, tt.args.high); (err != nil) != tt.wantErr {
				t.Errorf("gpioRepository.SetPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gpioRepository_GetPin(t *testing.T) {
	type args struct {
		pin machine.Pin
	}
	tests := []struct {
		name    string
		r       tgk.GPIORepository
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetPin(tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("gpioRepository.GetPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gpioRepository.GetPin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gpioRepository_ToPin(t *testing.T) {
	type args struct {
		pin string
	}
	tests := []struct {
		name string
		r    tgk.GPIORepository
		args args
		want machine.Pin
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ToPin(tt.args.pin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gpioRepository.ToPin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gpioRepository_ToGpios(t *testing.T) {
	type args struct {
		matrixPins map[string][]string
	}
	tests := []struct {
		name string
		r    tgk.GPIORepository
		args args
		want map[string]tgk.GpioConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ToGpios(tt.args.matrixPins); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gpioRepository.ToGpios() = %v, want %v", got, tt.want)
			}
		})
	}
}

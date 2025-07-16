package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewHIDService(t *testing.T) {
	tests := []struct {
		name string
		want tgk.HIDService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewHIDService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHIDService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hidService_RegisterHID(t *testing.T) {
	type args struct {
		hid tgk.HIDInterface
	}
	tests := []struct {
		name string
		s    tgk.HIDService
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.RegisterHID(tt.args.hid)
		})
	}
}

func Test_hidService_Init(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.HIDService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("hidService.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hidService_SendKeyboardReport(t *testing.T) {
	type args struct {
		report []byte
	}
	tests := []struct {
		name    string
		s       tgk.HIDService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.SendKeyboardReport(tt.args.report); (err != nil) != tt.wantErr {
				t.Errorf("hidService.SendKeyboardReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

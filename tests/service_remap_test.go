package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewRemapService(t *testing.T) {
	tests := []struct {
		name string
		want tgk.RemapService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewRemapService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRemapService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_remapService_Init(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.RemapService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("remapService.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_remapService_Remap(t *testing.T) {
	type args struct {
		keycode int
	}
	tests := []struct {
		name    string
		s       tgk.RemapService
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Remap(tt.args.keycode)
			if (err != nil) != tt.wantErr {
				t.Errorf("remapService.Remap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("remapService.Remap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_remapService_SetRemap(t *testing.T) {
	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name    string
		s       tgk.RemapService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.SetRemap(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("remapService.SetRemap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_remapService_ClearRemap(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.RemapService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ClearRemap(); (err != nil) != tt.wantErr {
				t.Errorf("remapService.ClearRemap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

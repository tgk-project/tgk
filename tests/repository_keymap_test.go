package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewKeymapRepository(t *testing.T) {
	type args struct {
		configRepo tgk.ConfigRepository
	}
	tests := []struct {
		name string
		args args
		want tgk.KeymapRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewKeymapRepository(tt.args.configRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeymapRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keymapRepository_Init(t *testing.T) {
	tests := []struct {
		name    string
		r       tgk.KeymapRepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Init(); (err != nil) != tt.wantErr {
				t.Errorf("keymapRepository.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_keymapRepository_Load(t *testing.T) {
	tests := []struct {
		name    string
		r       tgk.KeymapRepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Load(); (err != nil) != tt.wantErr {
				t.Errorf("keymapRepository.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_keymapRepository_Save(t *testing.T) {
	tests := []struct {
		name    string
		r       tgk.KeymapRepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Save(); (err != nil) != tt.wantErr {
				t.Errorf("keymapRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_keymapRepository_GetKeymap(t *testing.T) {
	tests := []struct {
		name string
		r    tgk.KeymapRepository
		want [][]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetKeymap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keymapRepository.GetKeymap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keymapRepository_SetKeymap(t *testing.T) {
	type args struct {
		keymap [][]int
	}
	tests := []struct {
		name    string
		r       tgk.KeymapRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.SetKeymap(tt.args.keymap); (err != nil) != tt.wantErr {
				t.Errorf("keymapRepository.SetKeymap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

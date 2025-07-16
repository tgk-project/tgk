package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewConfigRepository(t *testing.T) {
	tests := []struct {
		name string
		want tgk.ConfigRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewConfigRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configRepository_Init(t *testing.T) {
	tests := []struct {
		name    string
		r       tgk.ConfigRepository
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Init(); (err != nil) != tt.wantErr {
				t.Errorf("configRepository.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configRepository_Load(t *testing.T) {
	type args struct {
		configFileName string
	}
	tests := []struct {
		name    string
		r       tgk.ConfigRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Load(tt.args.configFileName); (err != nil) != tt.wantErr {
				t.Errorf("configRepository.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configRepository_GetConfig(t *testing.T) {
	tests := []struct {
		name string
		r    tgk.ConfigRepository
		want *tgk.KeyboardConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configRepository.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

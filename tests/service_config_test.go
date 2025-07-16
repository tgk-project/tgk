package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewConfigService(t *testing.T) {
	type args struct {
		configRepository tgk.ConfigRepository
	}
	tests := []struct {
		name string
		args args
		want tgk.ConfigService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewConfigService(tt.args.configRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configService_GetConfig(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.ConfigService
		want    *tgk.KeyboardConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("configService.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configService.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configService_SetLoadConfigName(t *testing.T) {
	type args struct {
		configFileName string
	}
	tests := []struct {
		name string
		s    tgk.ConfigService
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetLoadConfigName(tt.args.configFileName)
		})
	}
}

func Test_configService_Load(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.ConfigService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Load(); (err != nil) != tt.wantErr {
				t.Errorf("configService.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

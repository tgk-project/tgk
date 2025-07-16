package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewTGKManager(t *testing.T) {
	type args struct {
		hidService     tgk.HIDService
		keyScanService tgk.KeyScanService
		layerService   tgk.LayerService
		remapService   tgk.RemapService
		splitService   tgk.SplitService
		configService  tgk.ConfigService
	}
	tests := []struct {
		name string
		args args
		want tgk.TGKManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewTGKManager(tt.args.hidService, tt.args.keyScanService, tt.args.layerService, tt.args.remapService, tt.args.splitService, tt.args.configService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTGKManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tgkManager_Init(t *testing.T) {
	tests := []struct {
		name    string
		m       tgk.TGKManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Init(); (err != nil) != tt.wantErr {
				t.Errorf("tgkManager.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tgkManager_SetLoadConfigName(t *testing.T) {
	type args struct {
		configFileName string
	}
	tests := []struct {
		name string
		m    tgk.TGKManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetLoadConfigName(tt.args.configFileName)
		})
	}
}

func Test_tgkManager_Task(t *testing.T) {
	tests := []struct {
		name    string
		m       tgk.TGKManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Task(); (err != nil) != tt.wantErr {
				t.Errorf("tgkManager.Task() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

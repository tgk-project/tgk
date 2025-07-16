package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewLayerService(t *testing.T) {
	type args struct {
		keymapRepo tgk.KeymapRepository
		configRepo tgk.ConfigRepository
	}
	tests := []struct {
		name string
		args args
		want tgk.LayerService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewLayerService(tt.args.keymapRepo, tt.args.configRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLayerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layerService_Init(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.LayerService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("layerService.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_layerService_GetCurrentLayer(t *testing.T) {
	tests := []struct {
		name string
		s    tgk.LayerService
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetCurrentLayer(); got != tt.want {
				t.Errorf("layerService.GetCurrentLayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layerService_SetLayer(t *testing.T) {
	type args struct {
		layer int
	}
	tests := []struct {
		name    string
		s       tgk.LayerService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.SetLayer(tt.args.layer); (err != nil) != tt.wantErr {
				t.Errorf("layerService.SetLayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_layerService_GetKeymap(t *testing.T) {
	type args struct {
		layer int
		key   int
	}
	tests := []struct {
		name    string
		s       tgk.LayerService
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetKeymap(tt.args.layer, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("layerService.GetKeymap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("layerService.GetKeymap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layerService_IsLayerKey(t *testing.T) {
	type args struct {
		key int
	}
	tests := []struct {
		name string
		s    tgk.LayerService
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsLayerKey(tt.args.key); got != tt.want {
				t.Errorf("layerService.IsLayerKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layerService_GetLayerName(t *testing.T) {
	tests := []struct {
		name string
		s    tgk.LayerService
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetLayerName(); got != tt.want {
				t.Errorf("layerService.GetLayerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layerService_LayerTask(t *testing.T) {
	type args struct {
		layer  int
		isHold bool
	}
	tests := []struct {
		name    string
		s       tgk.LayerService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.LayerTask(tt.args.layer, tt.args.isHold); (err != nil) != tt.wantErr {
				t.Errorf("layerService.LayerTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

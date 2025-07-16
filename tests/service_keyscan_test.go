package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewKeyScanService(t *testing.T) {
	type args struct {
		keyScan tgk.KeyScanMatrix
	}
	tests := []struct {
		name string
		args args
		want tgk.KeyScanService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewKeyScanService(tt.args.keyScan); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeyScanService() = %v, want %v", got, tt.want)
			}
		})
	}
}

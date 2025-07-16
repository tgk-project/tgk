package tests

import (
	"reflect"
	"testing"

	"github.com/Diwamoto/tgk"
)

func TestNewSplitService(t *testing.T) {
	tests := []struct {
		name string
		want tgk.SplitService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tgk.NewSplitService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSplitService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitService_Init(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.SplitService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("splitService.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_splitService_IsMaster(t *testing.T) {
	tests := []struct {
		name string
		s    tgk.SplitService
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsMaster(); got != tt.want {
				t.Errorf("splitService.IsMaster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitService_Send(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       tgk.SplitService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Send(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("splitService.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_splitService_Receive(t *testing.T) {
	tests := []struct {
		name    string
		s       tgk.SplitService
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Receive()
			if (err != nil) != tt.wantErr {
				t.Errorf("splitService.Receive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitService.Receive() = %v, want %v", got, tt.want)
			}
		})
	}
}

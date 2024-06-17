// Package metrics ...
package metrics

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
)

func TestGetClient(t *testing.T) {
	defer func() {
		clientFactory = make(map[string]Client)
	}()
	RegisterClient(EmptyMetricsClient)
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  Client
		want1 bool
	}{
		{
			name:  "normal",
			args:  args{name: "empty"},
			want:  EmptyMetricsClient,
			want1: true,
		},
		{
			name:  "normal",
			args:  args{name: "n"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetClient(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRegisterSendDataHook(t *testing.T) {
	defer func() {
		sendDataHook = nil
	}()
	type args struct {
		handler func(metadata *Metadata, data [][]string) error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{handler: nil},
		},
		{
			name: "normal",
			args: args{handler: func(metadata *Metadata, data [][]string) error {
				return nil
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterSendDataHook(tt.args.handler)
		})
	}
}

func TestSendData(t *testing.T) {
	type args struct {
		ctx      context.Context
		metadata *Metadata
		data     [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty data",
			args: args{
				ctx:      context.TODO(),
				metadata: nil,
				data:     nil,
			},
			wantErr: false,
		},
		{
			name: "panic, nil metadata",
			args: args{
				ctx:      context.TODO(),
				metadata: nil,
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sample",
			args: args{
				ctx: context.TODO(),
				metadata: &Metadata{
					SamplingInterval: math.MaxUint32,
				},
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sample",
			args: args{
				ctx: context.TODO(),
				metadata: &Metadata{
					SamplingInterval: 0,
				},
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not client",
			args: args{
				ctx: context.TODO(),
				metadata: &Metadata{
					SamplingInterval: 1,
				},
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendData(tt.args.ctx, tt.args.metadata, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SendData error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendData1(t *testing.T) {
	defer func() {
		sendDataHook = nil
		clientFactory = make(map[string]Client)
	}()
	RegisterClient(EmptyMetricsClient)
	sendDataHook = func(metadata *Metadata, data [][]string) error {
		return nil
	}
	type args struct {
		ctx      context.Context
		metadata *Metadata
		data     [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				metadata: &Metadata{
					SamplingInterval:  1,
					MetricsPluginName: "empty",
				},
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendData(tt.args.ctx, tt.args.metadata, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// hook 报错中断上报
func TestSendData2(t *testing.T) {
	defer func() {
		sendDataHook = nil
		clientFactory = make(map[string]Client)
	}()
	RegisterClient(EmptyMetricsClient)
	sendDataHook = func(metadata *Metadata, data [][]string) error {
		return errors.Errorf("mock metrics err")
	}
	type args struct {
		ctx      context.Context
		metadata *Metadata
		data     [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				metadata: &Metadata{
					SamplingInterval:  1,
					MetricsPluginName: "empty",
				},
				data: [][]string{
					{
						"123", "234",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendData(tt.args.ctx, tt.args.metadata, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SendData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type empty struct{}

func (a empty) LogExposure(ctx context.Context, metadata *Metadata, exposureGroup *protoc_event_server.ExposureGroup) error {
	return nil
}

func (a empty) LogEvent(ctx context.Context, metadata *Metadata, eventGroup *protoc_event_server.EventGroup) error {
	return nil
}

func (a empty) LogMonitorEvent(ctx context.Context, metadata *Metadata, eventGroup *protoc_event_server.MonitorEventGroup) error {
	return nil
}

// Init TODO
func (a empty) Init(ctx context.Context, config *protoc_cache_server.MetricsInitConfig) error {
	return nil
}

// SendData TODO
func (a empty) SendData(ctx context.Context, metadata *Metadata, data [][]string) error {
	return nil
}

// Name TODO
func (a empty) Name() string {
	return "empty"
}

var (
	// EmptyMetricsClient TODO
	// Deprecated: for test
	EmptyMetricsClient = &empty{}
)

// Package env ...
package env

import (
	"testing"

	"github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
)

func TestInvokePath(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{skip: 4},
			want: ":0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvokePath(tt.args.skip); got != tt.want {
				t.Errorf("InvokePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

var invalidAddr = []byte{
	0x7f,
}

func TestRegisterAddr(t *testing.T) {
	type args struct {
		envType Type
		addr    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				envType: TypePrd,
				addr:    DefaultAddrPrd,
			},
			wantErr: false,
		},
		{
			name: "normal",
			args: args{
				envType: TypePrd,
				addr:    string(invalidAddr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterAddr(tt.args.envType, tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("RegisterAddr() error = %v, wantErr %v, addr=%s", err, tt.wantErr, tt.args.addr)
			}
		})
	}
}

func TestGetAddr(t *testing.T) {
	type args struct {
		envType Type
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{envType: "xx"},
			want: DefaultAddrPrd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAddr(tt.args.envType); got != tt.want {
				t.Errorf("GetAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterDMPAddr(t *testing.T) {
	type args struct {
		envType Type
		addr    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				envType: TypePrd,
				addr:    string(invalidAddr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterDMPAddr(tt.args.envType, tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("RegisterDMPAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDMPAddr(t *testing.T) {
	type args struct {
		envType Type
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{envType: "xxx"},
			want: DefaultDMPAddrPrd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDMPAddr(tt.args.envType); got != tt.want {
				t.Errorf("GetDMPAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrMsg(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{err: errors.Errorf("mock err")},
			want: "mock err",
		},
		{
			name: "normal",
			args: args{err: nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrMsg(tt.args.err); got != tt.want {
				t.Errorf("ErrMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventStatus(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want protoc_event_server.MonitorEventStatus
	}{
		{
			name: "normal",
			args: args{err: errors.Errorf("mock err")},
			want: protoc_event_server.MonitorEvent_STATUS_UNEXPECTED,
		},
		{
			name: "normal",
			args: args{err: nil},
			want: protoc_event_server.MonitorEvent_STATUS_SUCCESS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EventStatus(tt.args.err); got != tt.want {
				t.Errorf("EventStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONString(t *testing.T) {
	type args struct {
		source interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{source: nil},
			want: "",
		},
		{
			name: "normal",
			args: args{source: "xxx"},
			want: `"xxx"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONString(tt.args.source); got != tt.want {
				t.Errorf("JSONString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamplingInterval(t *testing.T) {
	type args struct {
		config *protoc_cache_server.MetricsConfig
		err    error
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "normal",
			args: args{
				config: &protoc_cache_server.MetricsConfig{
					SamplingInterval:    100000,
					ErrSamplingInterval: 1,
				},
				err: nil,
			},
			want: 100000,
		},
		{
			name: "normal",
			args: args{
				config: &protoc_cache_server.MetricsConfig{
					SamplingInterval:    100000,
					ErrSamplingInterval: 1,
				},
				err: errors.Errorf("mock err"),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SamplingInterval(tt.args.config, tt.args.err); got != tt.want {
				t.Errorf("SamplingInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

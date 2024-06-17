// Package abc ...
package abc

import (
	"context"
	"reflect"
	"testing"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/testdata"
)

var (
	projectIDList = []string{"123"}
	projectID     = "123"
)

func TestGetGlobalConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *internal.GlobalConfig
		wantErr bool
	}{
		{
			name:    "normal",
			want:    internal.C,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGlobalConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGlobalConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitWithRPCMock(t *testing.T) {
	ctx := context.Background()
	mockClient := testdata.MockCacheClient(t)
	type args struct {
		ctx           context.Context
		projectIDList []string
		opts          []InitOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "projectID is required",
			args: args{
				ctx:           ctx,
				projectIDList: nil,
				opts:          nil,
			},
			wantErr: true,
		},
		{
			name: "pass",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts:          []InitOption{WithRegisterCacheClient(mockClient), WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil)},
			},
			wantErr: false,
		},
		{
			name: "client is required",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts:          []InitOption{WithRegisterCacheClient(nil)},
			},
			wantErr: true,
		},
		{
			name: "dmpClient is required",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts:          []InitOption{WithRegisterDMPClient(nil)},
			},
			wantErr: true,
		},
		{
			name: "pass",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts: []InitOption{
					WithRegisterDMPClient(testdata.MockEmptyDMPClient),
					WithRegisterCacheClient(mockClient), WithSecretKey("")},
			},
			wantErr: false,
		},
		{
			name: "client should not be nil",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts:          []InitOption{WithRegisterMetricsPlugin(nil, nil)},
			},
			wantErr: true,
		},
		{
			name: "pass",
			args: args{
				ctx:           ctx,
				projectIDList: projectIDList,
				opts: []InitOption{
					WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
					WithRegisterDMPClient(client.NewDMPClient()),
					WithRegisterCacheClient(mockClient)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Release()
			if err := Init(tt.args.ctx, tt.args.projectIDList, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		ctx           context.Context
		projectIDList []string
		opts          []InitOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "projectID is required",
			args: args{
				ctx:           nil,
				projectIDList: nil,
				opts:          nil,
			},
			wantErr: true,
		},
		{
			name: "401",
			args: args{
				ctx:           context.Background(),
				projectIDList: []string{"not exist"},
				opts:          nil,
			},
			wantErr: true,
		},
		{
			name: "401",
			args: args{
				ctx:           context.Background(),
				projectIDList: []string{"123"},
				opts:          nil,
			},
			wantErr: true,
		},
		{
			name: "401",
			args: args{
				ctx:           context.Background(),
				projectIDList: []string{"123"},
				opts:          []InitOption{WithEnvType(env.TypePrd)},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.args.ctx, tt.args.projectIDList, tt.args.opts...)
			t.Logf("%v", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

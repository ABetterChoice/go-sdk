// Package abc ...
package abc

import (
	"context"
	"encoding/base64"
	"fmt"
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

func buildTokenWithEnv(envValue string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"TAB_DATA_ENV":"%s"}`, envValue)))
	return fmt.Sprintf("%s.%s.signature", header, payload)
}

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

func TestWithSecretKeyTokenDataEnv(t *testing.T) {
	testToken := buildTokenWithEnv("test")
	grayToken := buildTokenWithEnv("gray")

	t.Run("parse TAB_DATA_ENV from token", func(t *testing.T) {
		config := &internal.GlobalConfig{EnvType: env.TypePrd}
		err := WithSecretKey(testToken)(config)
		if err != nil {
			t.Fatalf("WithSecretKey() error = %v", err)
		}
		if config.TokenDataEnv != "test" {
			t.Fatalf("TokenDataEnv = %s, want %s", config.TokenDataEnv, "test")
		}
		if config.EnvType != env.TypePrd {
			t.Fatalf("EnvType = %s, want %s", config.EnvType, env.TypePrd)
		}
	})

	t.Run("keep EnvType unchanged after parsing token env", func(t *testing.T) {
		config := &internal.GlobalConfig{EnvType: env.TypePrd}
		if err := WithEnvType(env.TypePrd)(config); err != nil {
			t.Fatalf("WithEnvType() error = %v", err)
		}
		if err := WithSecretKey(testToken)(config); err != nil {
			t.Fatalf("WithSecretKey() error = %v", err)
		}
		if config.EnvType != env.TypePrd {
			t.Fatalf("EnvType = %s, want %s", config.EnvType, env.TypePrd)
		}
	})

	t.Run("accept non standard token data env value", func(t *testing.T) {
		config := &internal.GlobalConfig{EnvType: env.TypePrd}
		if err := WithSecretKey(grayToken)(config); err != nil {
			t.Fatalf("WithSecretKey() error = %v", err)
		}
		if config.TokenDataEnv != "gray" {
			t.Fatalf("TokenDataEnv = %s, want %s", config.TokenDataEnv, "gray")
		}
	})

	t.Run("clear stale token data env when parsing fails", func(t *testing.T) {
		config := &internal.GlobalConfig{TokenDataEnv: "3"}
		if err := WithSecretKey("invalid_token")(config); err != nil {
			t.Fatalf("WithSecretKey() error = %v", err)
		}
		if config.TokenDataEnv != "" {
			t.Fatalf("TokenDataEnv = %s, want empty", config.TokenDataEnv)
		}
	})
}

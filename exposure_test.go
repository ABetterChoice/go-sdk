// Package abc ...
package abc

import (
	"context"
	"fmt"
	"testing"

	"github.com/abetterchoice/go-sdk/testdata"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/stretchr/testify/assert"
)

func TestExposureExperiment(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		result    *ExperimentResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal", // no data
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				result:    nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "normal", // no data
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				result: &ExperimentResult{
					userCtx: &userContext{
						err: nil,
					},
					Group: &Group{
						ID:            301001001,
						Key:           "301001001",
						ExperimentKey: "301001",
						LayerKey:      "doubleHashLayerTag",
						IsControl:     true,
						params: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
						sceneIDList: nil,
						UnitIDType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "normal", // no data
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				result: &ExperimentResult{
					userCtx: &userContext{
						err: nil,
					},
					Group: &Group{
						ID:            301001001,
						Key:           "301001001",
						ExperimentKey: "301001",
						LayerKey:      "doubleHashLayerTag",
						IsControl:     true,
						params: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
						sceneIDList: []int64{1, 2, 3},
						UnitIDType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, LogExperimentExposure(tt.args.ctx, tt.args.projectID, tt.args.result), fmt.Sprintf("LogExperimentExposure(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.result))
		})
	}
}

func TestExposureExperiments(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		list      *ExperimentList
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				list: &ExperimentList{
					userCtx: &userContext{
						err:          nil,
						expandedData: map[string]string{"keyA": "value1"},
					},
					Data: map[string]*Group{
						"doubleHashLayer1": &Group{
							ID:            301001001,
							Key:           "301001001",
							ExperimentKey: "301001",
							LayerKey:      "doubleHashLayer1",
							IsControl:     true,
							params: map[string]string{
								"key1": "value1",
							},
							sceneIDList: nil,
							UnitIDType:  0,
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, LogExperimentsExposure(tt.args.ctx, tt.args.projectID, tt.args.list), fmt.Sprintf("LogExperimentsExposure(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.list))
		})
	}
}

func TestExposureFeatureFlag(t *testing.T) {
	type args struct {
		ctx         context.Context
		projectID   string
		featureFlag *FeatureFlag
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				featureFlag: &FeatureFlag{
					ConfigResult: &ConfigResult{
						userCtx: &userContext{},
						Config: &Config{
							Key:          "remoteConfig1",
							Value:        &Value{data: []byte("hello feature flag")},
							remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
							unitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				featureFlag: &FeatureFlag{
					ConfigResult: &ConfigResult{
						userCtx: &userContext{},
						Config: &Config{
							Key:          "bitmapTest",
							Value:        &Value{data: []byte("hello feature flag")},
							remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["bitmapTest"],
							unitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, LogFeatureFlagExposure(tt.args.ctx, tt.args.projectID, tt.args.featureFlag), fmt.Sprintf("LogFeatureFlagExposure(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.featureFlag))
		})
	}
}

func TestExposureRemoteConfig(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		config    *ConfigResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				config: &ConfigResult{
					userCtx: &userContext{},
					Config: &Config{
						Key:          "remoteKey1",
						Value:        &Value{data: []byte("hello remoteKey value")},
						remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
						unitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "normal", // 带场景
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				config: &ConfigResult{
					userCtx: &userContext{},
					Config: &Config{
						Key:          "remoteKey1",
						Value:        &Value{data: []byte("hello remoteKey value")},
						remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
						unitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "normal", // 带场景，场景没有配置 metricsConfig
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				config: &ConfigResult{
					userCtx: &userContext{},
					Config: &Config{
						Key:          "bitmapTest",
						Value:        &Value{data: []byte("hello remoteKey value")},
						remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["bitmapTest"],
						unitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, LogRemoteConfigExposure(tt.args.ctx, tt.args.projectID, tt.args.config), fmt.Sprintf("LogRemoteConfigExposure(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.config))
		})
	}
}

func BenchmarkInt64Join(b *testing.B) {
	var source = []int64{1, 2, 3, 4, 4, 5, 6, 3, 3, 4, 5, 5, 523424, 23, 4}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = int64ListJoin(source, ";")
		}
	})
}

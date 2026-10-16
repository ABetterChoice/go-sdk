// Package config ...
package config

import (
	"context"
	"reflect"
	"testing"

	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/abetterchoice/go-sdk/testdata"
	"github.com/abetterchoice/protoc_cache_server"
)

var (
	projectID     = "123"
	projectIDList = []string{"123"}
)

func Test_executor_GetRemoteConfig(t *testing.T) {
	mockInitLocalCache(t)
	type args struct {
		ctx       context.Context
		projectID string
		key       string
		options   *experiment.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Value
		wantErr bool
	}{
		{
			name: "projectID not found",
			args: args{
				ctx:       context.TODO(),
				projectID: "emptyKey",
				key:       "emptyKey",
				options:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "remoteConfig not found",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "emptyKey",
				options:   &experiment.Options{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "remoteConfig1",
				options: &experiment.Options{
					Application: nil,
				},
			},
			want: &Value{
				Data:           []byte("remoteConfig1-condition1"),
				IsOverrideList: false,
				IsDefault:      false,
				RemoteConfig:   testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
				UnitIDType:     protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "override",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "remoteConfig1",
				options: &experiment.Options{
					Application: nil,
					UnitID:      "overrideUnitID",
				},
			},
			want: &Value{
				Data:           []byte("hitOverrideResult"),
				IsOverrideList: true,
				IsDefault:      false,
				RemoteConfig:   testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
				UnitIDType:     protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "bitmapTest",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "bitmapTest",
				options: &experiment.Options{
					Application: nil,
				},
			},
			want: &Value{
				Data:         []byte("bitmapTestDefaultValue"),
				IsDefault:    true,
				RemoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["bitmapTest"],
				UnitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "withTag",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "withTag",
				options: &experiment.Options{
					Application: nil,
				},
			},
			want: &Value{
				Data:         []byte("withTagDefaultValue"),
				IsDefault:    true,
				RemoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["withTag"],
				UnitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "withTag",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "withTag",
				options: &experiment.Options{
					Application: nil,
					AttributeTag: map[string][]string{
						"tagKey1": []string{"ios"},
					},
				},
			},
			want: &Value{
				Data:         []byte("withTag-condition1"),
				IsDefault:    false,
				RemoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["withTag"],
				UnitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "withTag",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "withTag",
				options: &experiment.Options{
					Application: nil,
					AttributeTag: map[string][]string{
						"tagKey1": []string{"ios", "iphone"},
					},
				},
			},
			want: &Value{
				Data:         []byte("withTagDefaultValue"),
				IsDefault:    true,
				RemoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["withTag"],
				UnitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
		{
			name: "withExperiment",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "withExperiment",
				options: &experiment.Options{
					Application: nil,
					AttributeTag: map[string][]string{
						"tagKey1": []string{"ios", "iphone"},
					},
				},
			},
			want: &Value{
				Data:         []byte("withExperiment-condition1"),
				IsDefault:    false,
				RemoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["withExperiment"],
				UnitIDType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.GetRemoteConfig(tt.args.ctx, tt.args.projectID, tt.args.key, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRemoteConfig() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func mockInitLocalCache(t *testing.T) {
	cacheClient := testdata.MockCacheClient(t)
	client.RegisterCacheClient(cacheClient)
	client.RegisterDMPClient(testdata.MockEmptyDMPClient)
	err := cache.InitLocalCache(context.TODO(), projectIDList)
	if err != nil {
		t.Fatalf("initLocalCache:%v", err)
	}
}

func Test_executor_getHashSource(t *testing.T) {
	type args struct {
		unitType protoc_cache_server.UnitIDType
		options  *experiment.Options
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			if got := e.getHashSource(tt.args.unitType, tt.args.options); got != tt.want {
				t.Errorf("getHashSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_getRemoteConfigValue(t *testing.T) {
	type args struct {
		ctx     context.Context
		config  *protoc_cache_server.RemoteConfig
		options *experiment.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.getRemoteConfigValue(tt.args.ctx, tt.args.config, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRemoteConfigValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRemoteConfigValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHitConditionBucketInfo(t *testing.T) {
	type args struct {
		bucketNum  int64
		bucketInfo *protoc_cache_server.BucketInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			if got := e.isHitConditionBucketInfo(tt.args.bucketNum, tt.args.bucketInfo); got != tt.want {
				t.Errorf("isHitConditionBucketInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_processCondition(t *testing.T) {
	type args struct {
		ctx       context.Context
		condition *protoc_cache_server.Condition
		options   *experiment.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Value
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, got1, err := e.processCondition(tt.args.ctx, tt.args.condition, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("processCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processCondition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("processCondition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_executor_processOverrideList(t *testing.T) {
	type args struct {
		config  *protoc_cache_server.RemoteConfig
		options *experiment.Options
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, _, got1 := e.processOverrideList(tt.args.config, tt.args.options)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processOverrideList() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("processOverrideList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_processConditionExperiment(t *testing.T) {
	type args struct {
		ctx       context.Context
		condition *protoc_cache_server.Condition
		options   *experiment.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processConditionExperiment(tt.args.ctx, tt.args.condition, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("processConditionExperiment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processConditionExperiment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

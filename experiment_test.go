// Package abc ...
package abc

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/testdata"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Test GetExperiment
func Test_userContext_GetExperiment(t *testing.T) {
	err := Init(context.Background(), projectIDList, WithRegisterCacheClient(testdata.MockCacheClient(t)), WithRegisterDMPClient(testdata.MockEmptyDMPClient))
	assert.Nil(t, err)
	type fields struct {
		err           error
		tags          map[string][]string
		unitID        string
		decisionID    string
		newUnitID     string
		newDecisionID string
		expandedData  map[string]string
	}
	type args struct {
		ctx       context.Context
		projectID string
		layerKey  string
		opts      []ExperimentOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ExperimentResult
		wantErr bool
	}{
		{
			name:    "business not found",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ctx err not nil",
			fields: fields{
				err: errors.Errorf("ctx err"),
			},
			args: args{
				projectID: projectID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "invalid layerKey",
			fields: fields{},
			args: args{
				projectID: projectID,
			},
			want:    nil,
			wantErr: true,
		},
		{ // The simplest possible test example
			name:   "normal",
			fields: fields{},
			args: args{
				ctx:       nil,
				projectID: projectID,
				layerKey:  "overrideLayer",
				opts:      nil,
			},
			want: &ExperimentResult{
				Group: &Group{
					ID:             100002001,
					Key:            "100002001",
					ExperimentKey:  "100002",
					LayerKey:       "overrideLayer",
					IsDefault:      false,
					IsControl:      true,
					IsOverrideList: false,
					params:         map[string]string{"key1": "100002001"},
					sceneIDList:    nil,
					UnitIDType:     protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				},
				userCtx: &userContext{
					err:           nil,
					tags:          nil,
					unitID:        "",
					decisionID:    "",
					newUnitID:     "",
					newDecisionID: "",
					expandedData:  nil,
				},
			},
			wantErr: false,
		},
		{ // Non-natural full traffic layer Missing situation
			name:   "invalid layer",
			fields: fields{},
			args: args{
				ctx:       nil,
				projectID: projectID,
				layerKey:  "subDomain-multiDomain1-multiLayer1",
				opts:      nil,
			},
			want:    nil,
			wantErr: false,
		},
		{ // Non-natural full traffic layer hit situation
			name:   "invalid layer",
			fields: fields{},
			args: args{
				ctx:       nil,
				projectID: projectID,
				layerKey:  "subDomain-holdoutDomain1-singleLayer",
				opts:      nil,
			},
			want: &ExperimentResult{
				userCtx: &userContext{},
				Group: &Group{
					ID:            200002001,
					Key:           "200002001",
					ExperimentKey: "200002",
					LayerKey:      "subDomain-holdoutDomain1-singleLayer",
					IsControl:     true,
					params: map[string]string{
						"key1": "200002001",
					},
					UnitIDType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				},
			},
			wantErr: false,
		},
		{ // If the experiment key is specified and it is not consistent with the expected hit, the result will return the hit experiment on the layer.
			name:   "normal",
			fields: fields{},
			args: args{
				ctx:       nil,
				projectID: projectID,
				layerKey:  "overrideLayer",
				opts:      []ExperimentOption{WithExperimentKey("100001")},
			},
			want: &ExperimentResult{
				userCtx: &userContext{
					err:           nil,
					tags:          nil,
					unitID:        "",
					decisionID:    "",
					newUnitID:     "",
					newDecisionID: "",
					expandedData:  nil,
				},
				Group: &Group{
					ID:             100002001,
					Key:            "100002001",
					ExperimentKey:  "100002",
					LayerKey:       "overrideLayer",
					IsDefault:      false,
					IsControl:      true,
					IsOverrideList: false,
					params:         map[string]string{"key1": "100002001"},
					sceneIDList:    nil,
					UnitIDType:     protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				},
			},
			wantErr: false,
		},
		{ // The opts scenario conflicts with layerKey. Under the logical AND condition, no related layers satisfy both the scenario and layerKey. Therefore, no experiment is hit and invalid layerKey is returned.
			name: "invalid layerKey",
			fields: fields{
				err:           nil,
				tags:          nil,
				unitID:        "",
				decisionID:    "",
				newUnitID:     "",
				newDecisionID: "",
				expandedData:  nil,
			},
			args: args{
				ctx:       nil,
				projectID: projectID,
				layerKey:  "overrideLayer",
				opts: []ExperimentOption{
					WithSceneIDList([]int64{1, 2, 3}),
					WithAutomatic(true),
					WithSceneID(4),
					WithLayerKey("overrideLayer"),
					WithLayerKeyList([]string{"overrideLayer"}),
					WithIsPreparedDMPTag(false),
					WithIsDisableDMP(false),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &userContext{
				err:           tt.fields.err,
				tags:          tt.fields.tags,
				unitID:        tt.fields.unitID,
				decisionID:    tt.fields.decisionID,
				newUnitID:     tt.fields.newUnitID,
				newDecisionID: tt.fields.newDecisionID,
				expandedData:  tt.fields.expandedData,
			}
			got, err := c.GetExperiment(tt.args.ctx, tt.args.projectID, tt.args.layerKey, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("got=%+v", got)
				t.Errorf("GetExperiment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExperiment() got = \n%+v, want \n%+v", env.JSONString(got), env.JSONString(tt.want))
			}
		})
	}
}

// Test whether the traffic ratio meets expectations
func Test_userContext_GetExperiment1(t *testing.T) {
	Release()
	// defer Release()
	err := Init(context.Background(), projectIDList, WithRegisterCacheClient(testdata.MockCacheClient(t)), WithRegisterDMPClient(testdata.MockEmptyDMPClient))
	assert.Nil(t, err)
	type fields struct {
		err           error
		tags          map[string][]string
		unitID        string
		decisionID    string
		newUnitID     string
		newDecisionID string
		expandedData  map[string]string
	}
	type args struct {
		ctx       context.Context
		projectID string
		layerKey  string
		opts      []ExperimentOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ExperimentResult
		wantErr assert.ErrorAssertionFunc
	}{
		{ // The simplest possible test example
			name:   "normal",
			fields: fields{},
			args: args{
				ctx:       context.Background(),
				projectID: projectID,
				opts:      nil,
			},
			want: nil,
		},
	}
	var result = make(map[int64]int64)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := 10000
			r := 0.03
			for i := 0; i < testTime; i++ {
				uuidStr := uuid.New().String()
				c := &userContext{
					err:           tt.fields.err,
					tags:          tt.fields.tags,
					unitID:        uuidStr,
					decisionID:    uuidStr,
					newUnitID:     uuidStr,
					newDecisionID: uuidStr,
					expandedData:  tt.fields.expandedData,
				}
				got, err := c.GetExperiments(tt.args.ctx, tt.args.projectID, tt.args.opts...)
				assert.Nil(t, err)
				assert.NotNil(t, got)
				for _, experiment := range got.Data {
					result[experiment.ID]++
				}
			}
			for groupID, hitNum := range result {
				switch groupID {
				case 100001001: // 层默认实验 预计有 60% 的流量
					assert.True(t, equalWithRange(hitNum, 6000, testTime, r))
				case 100002001: // 预计有 10% 的流量
					assert.True(t, equalWithRange(hitNum, 1000, testTime, r))
				case 100002002: // 预计有 10% 的流量
					assert.True(t, equalWithRange(hitNum, 1000, testTime, r))
				case 100003001: // 预计有 10% 的流量
					assert.True(t, equalWithRange(hitNum, 1000, testTime, r))
				case 100003002: // 预计有 10% 的流量
					assert.True(t, equalWithRange(hitNum, 1000, testTime, r))
				case 101001001:
					assert.True(t, equalWithRange(hitNum, 6000, testTime, r))
				case 200001001:
					assert.True(t, equalWithRange(hitNum, 300, testTime, r))
				case 200002001:
					assert.True(t, equalWithRange(hitNum, 50, testTime, r))
				case 200002002:
					assert.True(t, equalWithRange(hitNum, 50, testTime, r))
				case 200003001:
					assert.True(t, equalWithRange(hitNum, 50, testTime, r))
				case 200003002:
					assert.True(t, equalWithRange(hitNum, 50, testTime, r))
				}
			}
			t.Logf("%+v", result)
		})
	}
}

func equalWithRange(actual, want int64, total int, r float64) bool {
	if want == 0 {
		return actual == want
	}
	return math.Abs(float64(actual-want))/float64(total) < r || math.Abs(float64(actual-want))/float64(want) < 0.08
}

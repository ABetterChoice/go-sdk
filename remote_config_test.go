// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments, user configuration data retrieval,
// user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/abetterchoice/go-sdk/testdata"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_userContext_GetRemoteConfig(t *testing.T) {
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
		key       string
		opts      []ConfigOption
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *ConfigResult
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "ctx err", // userCtx err
			fields: fields{
				err: errors.Errorf("mock err"),
			},
			args: args{
				ctx:       context.TODO(),
				projectID: "123",
				key:       "remoteConfig1",
			},
			wantResult: nil,
			wantErr: assert.ErrorAssertionFunc(func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			}),
		},
		{
			name:   "projectID not found", // userCtx err
			fields: fields{},
			args: args{
				ctx:       context.TODO(),
				projectID: "emptyProjectID",
				key:       "remoteConfig1",
				opts:      []ConfigOption{forTestConfigOption()},
			},
			wantResult: nil,
			wantErr: assert.ErrorAssertionFunc(func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			}),
		},
		{
			name:   "normal", // userCtx err
			fields: fields{},
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				key:       "remoteConfig1",
				opts:      []ConfigOption{forTestConfigOption()},
			},
			wantResult: &ConfigResult{
				userCtx: &userContext{},
				Config: &Config{
					Key:          "remoteConfig1",
					Value:        &Value{data: []byte("remoteConfig1-condition1")},
					remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
					unitIDType:   protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				},
			},
			wantErr: assert.ErrorAssertionFunc(func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			}),
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
			gotResult, err := c.GetRemoteConfig(tt.args.ctx, tt.args.projectID, tt.args.key, tt.args.opts...)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRemoteConfig(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.key)) {
				t.Errorf("invalid err:%v", err)
				return
			}
			if gotResult == nil {
				return
			}
			gotResultBody, _ := json.Marshal(gotResult)
			wantResultBody, _ := json.Marshal(tt.wantResult)
			if !reflect.DeepEqual(gotResultBody, wantResultBody) {
				t.Errorf("GetRemoteConfig() got=\n%s, want=\n%s", gotResultBody, wantResultBody)
			}
			assert.NotNil(t, gotResult.Byte())
			assert.NotEmpty(t, gotResult.String())
			int64Value, err := gotResult.GetInt64()
			assert.Zero(t, int64Value)
			assert.NotNil(t, err)
			int64Value = gotResult.GetInt64WithDefault(1)
			assert.Equal(t, int64Value, int64(1))
			assert.Equal(t, gotResult.MustGetInt64(), int64(0))
			assert.NotNil(t, gotResult.MustGetBool(), false)
			boolValue, err := gotResult.GetBool()
			assert.False(t, boolValue)
			assert.NotNil(t, err)
			assert.Equal(t, gotResult.GetBoolWithDefault(true), true)
			float64Value, err := gotResult.GetFloat64()
			assert.Equal(t, float64Value, float64(0))
			assert.NotNil(t, err)
			assert.Equal(t, gotResult.GetFloat64WithDefault(10.9), 10.9)
			assert.Equal(t, gotResult.MustGetFloat64(), float64(0))
			jsonMapValue, err := gotResult.GetJSONMap()
			assert.Nil(t, jsonMapValue)
			assert.NotNil(t, err)
			assert.Nil(t, gotResult.MustGetJSONMap())
			assert.Nil(t, gotResult.GetJSONMapWithDefault(nil))
			// assert.Equalf(t, tt.wantResult, gotResult, "GetRemoteConfig(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.key)
		})
	}
}

func forTestConfigOption() ConfigOption {
	return func(options *experiment.Options) error {
		return nil
	}
}

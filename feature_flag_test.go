// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments, user configuration data retrieval,
// user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"fmt"
	"testing"

	"github.com/abetterchoice/go-sdk/testdata"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/stretchr/testify/assert"
)

func Test_userContext_GetFeatureFlag(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		want    *FeatureFlag
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "normal", // normal
			fields: fields{},
			args: args{
				ctx:       context.TODO(),
				projectID: "emptyProjectID",
				key:       "empty",
				opts:      []ConfigOption{forTestConfigOption()},
			},
			want: nil,
			wantErr: assert.ErrorAssertionFunc(func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			}),
		},
		{
			name:   "normal", // normal
			fields: fields{},
			args: args{
				ctx:       context.TODO(),
				projectID: "123",
				key:       "remoteConfig1",
				opts:      []ConfigOption{forTestConfigOption()},
			},
			want: &FeatureFlag{&ConfigResult{
				userCtx: &userContext{},
				Config: &Config{
					Key:          "remoteConfig1",
					Value:        &Value{data: []byte("remoteConfig1-condition1")},
					remoteConfig: testdata.NormalTabConfig.ConfigData.RemoteConfigIndex["remoteConfig1"],
					unitIDType:   protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				},
			}},
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
			got, err := c.GetFeatureFlag(tt.args.ctx, tt.args.projectID, tt.args.key, tt.args.opts...)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFeatureFlag(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.key)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFeatureFlag(%v, %v, %v)", tt.args.ctx, tt.args.projectID, tt.args.key)
			if got == nil {
				return
			}
			got.MustGetBool()
		})
	}
}

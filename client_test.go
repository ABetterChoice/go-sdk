// Package abc ...
package abc

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test the normal conditions of setting various IDs
func TestNewUserContext(t *testing.T) {
	type args struct {
		unitID string
		opts   []Attribution
	}
	tests := []struct {
		name string
		args args
		want Context
	}{
		{
			name: "unitID is required",
			args: args{
				unitID: "",
				opts:   nil,
			},
			want: &userContext{
				err:  fmt.Errorf("unitID is required"),
				tags: map[string][]string{},
			},
		},
		{
			name: "pass",
			args: args{
				unitID: "unitID",
				opts:   nil,
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "pass2",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithNewUnitID("newUnitID")},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "newUnitID",
				newDecisionID: "newUnitID",
			},
		},
		{
			name: "pass2",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithDecisionID("decisionID")},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "decisionID",
				newUnitID:     "unitID",
				newDecisionID: "decisionID",
			},
		},
		{
			name: "pass2",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithNewDecisionID("newDecisionID")},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "newDecisionID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserContext(tt.args.unitID, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test and set various ID error reports
func TestNewUserContext2(t *testing.T) {
	type args struct {
		unitID string
		opts   []Attribution
	}
	tests := []struct {
		name string
		args args
		want Context
	}{
		{
			name: "decisionID is required",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithDecisionID("")},
			},
			want: &userContext{
				err:           fmt.Errorf("decisionID is required"),
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "newUnitID is required",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithNewUnitID("")},
			},
			want: &userContext{
				err:           fmt.Errorf("newUnitID is required"),
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "newDecisionID is required",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithNewDecisionID("")},
			},
			want: &userContext{
				err:           fmt.Errorf("newDecisionID is required"),
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserContext(tt.args.unitID, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// tagInfo
func TestNewUserContext3(t *testing.T) {
	type args struct {
		unitID string
		opts   []Attribution
	}
	tests := []struct {
		name string
		args args
		want Context
	}{
		{
			name: "empty tags",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithTags(nil)},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "tags1",
			args: args{
				unitID: "unitID",
				opts: []Attribution{WithTags(map[string][]string{
					"sexy":    {"man"},
					"age":     {"27"},
					"country": {"cn", "us", "ru"},
				})},
			},
			want: &userContext{
				err: nil,
				tags: map[string][]string{
					"sexy":    {"man"},
					"age":     {"27"},
					"country": {"cn", "us", "ru"},
				},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "merge tag",
			args: args{
				unitID: "unitID",
				opts: []Attribution{WithTagKV("decisionID", "unitID"), WithTags(map[string][]string{
					"sexy":    {"man"},
					"age":     {"27"},
					"country": {"cn", "us", "ru"},
				}), WithTagKV("name", "unitID")},
			},
			want: &userContext{
				err: nil,
				tags: map[string][]string{
					"sexy":       {"man"},
					"age":        {"27"},
					"name":       {"unitID"},
					"country":    {"cn", "us", "ru"},
					"decisionID": {"unitID"},
				},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserContext(tt.args.unitID, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Set extended information
func TestNewUserContext4(t *testing.T) {
	type args struct {
		unitID string
		opts   []Attribution
	}
	tests := []struct {
		name string
		args args
		want Context
	}{
		{
			name: "empty expandedData",
			args: args{
				unitID: "unitID",
				opts:   []Attribution{WithExpandedData(nil)},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
			},
		},
		{
			name: "expandedData",
			args: args{
				unitID: "unitID",
				opts: []Attribution{WithExpandedData(map[string]string{
					"keyA": "valueA",
				})},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
				expandedData: map[string]string{
					"keyA": "valueA",
				},
			},
		},
		{
			name: "merge expandedData",
			args: args{
				unitID: "unitID",
				opts: []Attribution{WithExpandedData(map[string]string{
					"keyA": "valueA",
				}), WithExpandedData(map[string]string{
					"keyB": "valueB",
				})},
			},
			want: &userContext{
				err:           nil,
				tags:          map[string][]string{},
				unitID:        "unitID",
				decisionID:    "unitID",
				newUnitID:     "unitID",
				newDecisionID: "unitID",
				expandedData: map[string]string{
					"keyA": "valueA",
					"keyB": "valueB",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserContext(tt.args.unitID, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test client mock
func TestClientMock(t *testing.T) {
	mockCtx := NewMockContext(gomock.NewController(t))
	mockCtx.EXPECT().GetExperiments(gomock.Any(), gomock.Any()).Return(&ExperimentList{
		userCtx: nil,
		Data:    nil,
	}, nil)
	mockCtx.EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ExperimentResult{
		userCtx: nil,
		Group:   nil,
	}, nil)
	mockCtx.EXPECT().GetRemoteConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ConfigResult{
		userCtx: nil,
		Config:  nil,
	}, nil)
	mockCtx.EXPECT().GetFeatureFlag(gomock.Any(), gomock.Any(), gomock.Any()).Return(&FeatureFlag{}, nil)
	experimentList, err := mockCtx.GetExperiments(context.TODO(), "123")
	assert.Nil(t, err)
	assert.NotNil(t, experimentList)
	experimentResult, err := mockCtx.GetExperiment(context.TODO(), "123", "layerKey")
	assert.Nil(t, err)
	assert.NotNil(t, experimentResult)
	configResult, err := mockCtx.GetRemoteConfig(context.TODO(), "123", "configKey")
	assert.Nil(t, err)
	assert.NotNil(t, configResult)
	featureFlag, err := mockCtx.GetFeatureFlag(context.TODO(), "123", "featureFlagKey")
	assert.Nil(t, err)
	assert.NotNil(t, featureFlag)
}

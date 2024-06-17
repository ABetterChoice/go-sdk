// Package experiment ...
package experiment

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/testdata"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/stretchr/testify/assert"
)

var (
	projectID     = "123"
	projectIDList = []string{"123"}
)

func Test_executor_GetExperiments(t *testing.T) {
	mockInitLocalCache(t)
	type args struct {
		ctx       context.Context
		projectID string
		options   *Options
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Experiment
		wantErr bool
	}{
		{
			name: "abtest global domain",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					IsExposureLoggingAutomatic: true,
					OverrideList: map[string]int64{
						"overrideLayer":                      100001001,
						"subDomain-multiDomain1-multiLayer1": 201002002,
					},
					DMPTagResult:     map[string]bool{},
					IsPreparedDMPTag: true,
				},
			},
			want:    globalAbtestResult,
			wantErr: false,
		},
		{
			name: "abtest overrideLayer",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					IsExposureLoggingAutomatic: true,
					UnitID:                     "",
					DecisionID:                 "",
					LayerKeys: map[string]bool{
						"overrideLayer": true,
					},
					DMPTagResult:     map[string]bool{},
					IsPreparedDMPTag: true,
				},
			},
			want: map[string]*Experiment{
				"overrideLayer": &Experiment{
					Group: &protoccacheserver.Group{
						Id:            100002001,
						GroupKey:      "100002001",
						ExperimentId:  100002,
						ExperimentKey: "100002",
						Params: map[string]string{
							"key1": "100002001",
						},
						IsControl: true,
						LayerKey:  "overrideLayer",
						IssueInfo: &protoccacheserver.IssueInfo{
							IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
						},
						UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "abtest override",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					IsExposureLoggingAutomatic: true,
					UnitID:                     "overrideUnitID",
					NewUnitID:                  "newOverrideUnitID",
					OverrideList: map[string]int64{
						"overrideLayer": 100001001,
					},
					LayerKeys: map[string]bool{
						"overrideLayer": true,
					},
					DMPTagResult:     map[string]bool{},
					IsPreparedDMPTag: true,
				},
			},
			want: map[string]*Experiment{
				"overrideLayer": &Experiment{
					Group: &protoccacheserver.Group{
						Id:            100001001,
						GroupKey:      "100001001",
						ExperimentId:  100001,
						ExperimentKey: "100001",
						Params: map[string]string{
							"key1": "100001001",
						},
						IsDefault: true,
						LayerKey:  "overrideLayer",
						IssueInfo: &protoccacheserver.IssueInfo{
							IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
						},
						UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
					IsOverrideList: true,
				},
			},
			wantErr: false,
		},
		{
			name: "abtest global domain2",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					IsExposureLoggingAutomatic: true,
					OverrideList:               map[string]int64{},
					DecisionID:                 "123",
					DMPTagResult:               map[string]bool{},
					IsPreparedDMPTag:           false,
				},
			},
			want:    globalAbtestResult2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.GetExperiments(tt.args.ctx, tt.args.projectID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExperiments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// gotBody, err := json.Marshal(got)
			// assert.Nil(t, err)
			// wantBody, err := json.Marshal(tt.want)
			// assert.Nil(t, err)
			// assert.Equal(t, string(gotBody), string(wantBody))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExperiments() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

var (
	// decisionID is empty, hash hit subDomain-holdoutDomain
	globalAbtestResult = map[string]*Experiment{
		"doubleHashLayerCityTag": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            303001001,
				GroupKey:      "303001001",
				ExperimentId:  303001,
				ExperimentKey: "303001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerCityTag",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_CITY_TAG,
					TagListGroup: []*protoccacheserver.TagList{
						{
							TagList: []*protoccacheserver.Tag{
								{
									Key:         "dmpTagTest",
									TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
									Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
									Value:       "dmpCodeTest3",
									DmpPlatform: 3,
									UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
								},
							},
						},
					},
				},
				SceneIdList: nil,
				UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"doubleHashLayerPercentage": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            302001001,
				GroupKey:      "302001001",
				ExperimentId:  302001,
				ExperimentKey: "302001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerPercentage",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"doubleHashLayerTag": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            301001001,
				GroupKey:      "301001001",
				ExperimentId:  301001,
				ExperimentKey: "301001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerTag",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_TAG,
					TagListGroup: []*protoccacheserver.TagList{
						{
							TagList: []*protoccacheserver.Tag{
								{
									Key:         "dmpTagTest",
									TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
									Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
									Value:       "dmpCodeTest3",
									DmpPlatform: 3,
									UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
								},
							},
						},
					},
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"multiLayer2": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            101002001,
				GroupKey:      "101002001",
				ExperimentId:  101002,
				ExperimentKey: "101002",
				Params: map[string]string{
					"key1": "101002001",
				},
				IsControl: true,
				LayerKey:  "multiLayer2",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"overrideLayer": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            100001001,
				GroupKey:      "100001001",
				ExperimentId:  100001,
				ExperimentKey: "100001",
				Params: map[string]string{
					"key1": "100001001",
				},
				IsDefault: true,
				LayerKey:  "overrideLayer",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			IsOverrideList: true,
		},
		"subDomain-holdoutDomain1-singleLayer": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            200002001,
				GroupKey:      "200002001",
				ExperimentId:  200002,
				ExperimentKey: "200002",
				Params: map[string]string{
					"key1": "200002001",
				},
				IsControl: true,
				LayerKey:  "subDomain-holdoutDomain1-singleLayer",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"subDomain-multiDomain1-multiLayer1": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            201002002,
				GroupKey:      "201002002",
				ExperimentId:  201002,
				ExperimentKey: "201002",
				Params: map[string]string{
					"key1": "201002002",
				},
				LayerKey: "subDomain-multiDomain1-multiLayer1",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
			IsOverrideList: true,
		},
	}
	// decisionID is 123, hash hit subDomain-multiDomain
	globalAbtestResult2 = map[string]*Experiment{
		"doubleHashLayerCityTag": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            303001001,
				GroupKey:      "303001001",
				ExperimentId:  303001,
				ExperimentKey: "303001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerCityTag",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_CITY_TAG,
					TagListGroup: []*protoccacheserver.TagList{
						{
							TagList: []*protoccacheserver.Tag{
								{
									Key:         "dmpTagTest",
									TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
									Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
									Value:       "dmpCodeTest3",
									DmpPlatform: 3,
									UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
								},
							},
						},
					},
				},
				SceneIdList: nil,
				UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"doubleHashLayerPercentage": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            302001001,
				GroupKey:      "302001001",
				ExperimentId:  302001,
				ExperimentKey: "302001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerPercentage",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"doubleHashLayerTag": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            301001001,
				GroupKey:      "301001001",
				ExperimentId:  301001,
				ExperimentKey: "301001",
				IsControl:     true,
				LayerKey:      "doubleHashLayerTag",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_TAG,
					TagListGroup: []*protoccacheserver.TagList{
						{
							TagList: []*protoccacheserver.Tag{
								{
									Key:         "dmpTagTest",
									TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
									Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
									Value:       "dmpCodeTest3",
									DmpPlatform: 3,
									UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
								},
							},
						},
					},
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"multiLayer2": &Experiment{
			Group: &protoccacheserver.Group{
				Id:        env.DefaultGlobalGroupID,
				GroupKey:  env.DefaultGlobalGroupKey,
				IsDefault: true,
				LayerKey:  "multiLayer2",
			},
		},
		"overrideLayer": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            100001001,
				GroupKey:      "100001001",
				ExperimentId:  100001,
				ExperimentKey: "100001",
				Params: map[string]string{
					"key1": "100001001",
				},
				IsDefault: true,
				LayerKey:  "overrideLayer",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"subDomain-multiDomain1-multiLayer1": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            201001001,
				GroupKey:      "201001001",
				ExperimentId:  201001,
				ExperimentKey: "201001",
				Params: map[string]string{
					"key1": "201001001",
				},
				IsDefault: true,
				LayerKey:  "subDomain-multiDomain1-multiLayer1",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
		"subDomain-multiDomain1-multiLayer2": &Experiment{
			Group: &protoccacheserver.Group{
				Id:            202001001,
				GroupKey:      "202001001",
				ExperimentId:  202001,
				ExperimentKey: "202001",
				Params: map[string]string{
					"key1": "202001001",
				},
				IsDefault: true,
				LayerKey:  "subDomain-multiDomain1-multiLayer2",
				IssueInfo: &protoccacheserver.IssueInfo{
					IssueType: protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE,
				},
				UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
			},
		},
	}
)

func mockInitLocalCache(t *testing.T) {
	cacheClient := testdata.MockCacheClient(t)
	client.RegisterCacheClient(cacheClient)
	client.RegisterDMPClient(testdata.MockEmptyDMPClient)
	err := cache.InitLocalCache(context.TODO(), projectIDList)
	if err != nil {
		t.Fatalf("initLocalCache:%v", err)
	}
}

func Test_executor_GetExperiments3(t *testing.T) {
	mockInitLocalCache(t)
	type args struct {
		ctx       context.Context
		projectID string
		options   *Options
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Experiment
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					LayerKeys: map[string]bool{
						"doubleHashLayerTag": true,
					},
					AttributeTag: nil,
				},
			},
			want: map[string]*Experiment{
				"doubleHashLayerTag": &Experiment{
					Group: &protoccacheserver.Group{
						Id:            301001001,
						GroupKey:      "301001001",
						IsControl:     true,
						ExperimentKey: "301001",
						ExperimentId:  301001,
						LayerKey:      "doubleHashLayerTag",
						IssueInfo: &protoccacheserver.IssueInfo{
							IssueType: protoccacheserver.IssueType_ISSUE_TYPE_TAG,
							TagListGroup: []*protoccacheserver.TagList{
								{
									TagList: []*protoccacheserver.Tag{
										{
											Key:         "dmpTagTest",
											TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
											Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
											Value:       "dmpCodeTest3",
											DmpPlatform: 3,
											UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
										},
									},
								},
							},
						},
						UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.GetExperiments(tt.args.ctx, tt.args.projectID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExperiments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody, err := json.Marshal(got)
			if err != nil {
				t.Errorf("json marshal fail")
				return
			}
			wantBody, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("json marshal fail")
				return
			}
			if !reflect.DeepEqual(gotBody, wantBody) {
				t.Errorf("GetExperiments() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_executor_GetExperiments2(t *testing.T) {
	mockInitLocalCache(t)
	type args struct {
		ctx       context.Context
		projectID string
		options   *Options
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Experiment
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					LayerKeys: map[string]bool{
						"doubleHashLayerTag": true,
					},
					AttributeTag: nil,
				},
			},
			want: map[string]*Experiment{
				"doubleHashLayerTag": &Experiment{
					Group: &protoccacheserver.Group{
						Id:            301001001,
						GroupKey:      "301001001",
						ExperimentId:  301001,
						ExperimentKey: "301001",
						IsControl:     true,
						LayerKey:      "doubleHashLayerTag",
						IssueInfo: &protoccacheserver.IssueInfo{
							IssueType: protoccacheserver.IssueType_ISSUE_TYPE_TAG,
							TagListGroup: []*protoccacheserver.TagList{
								{
									TagList: []*protoccacheserver.Tag{
										{
											Key:         "dmpTagTest",
											TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
											Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
											Value:       "dmpCodeTest3",
											DmpPlatform: 3,
											UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
										},
									},
								},
							},
						},
						UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.GetExperiments(tt.args.ctx, tt.args.projectID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExperiments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody, _ := json.Marshal(got)
			wantBody, _ := json.Marshal(tt.want)
			if !reflect.DeepEqual(gotBody, wantBody) {
				t.Errorf("GetExperiments() got = \n%v, want \n%v", got, tt.want)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetExperiments() got = \n%v, want \n%v", got, tt.want)
			// }
		})
	}
}

func Test_executor_GetExperiments1(t *testing.T) {
	mockInitLocalCache(t)
	type args struct {
		ctx       context.Context
		projectID string
		options   *Options
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Experiment
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:       context.TODO(),
				projectID: projectID,
				options: &Options{
					LayerKeys: map[string]bool{
						"doubleHashLayerTag": true,
					},
					AttributeTag: nil,
				},
			},
			want: map[string]*Experiment{
				"doubleHashLayerTag": &Experiment{
					Group: &protoccacheserver.Group{
						Id:            301001001,
						GroupKey:      "301001001",
						ExperimentId:  301001,
						ExperimentKey: "301001",
						IsControl:     true,
						LayerKey:      "doubleHashLayerTag",
						IssueInfo: &protoccacheserver.IssueInfo{
							IssueType: protoccacheserver.IssueType_ISSUE_TYPE_TAG,
							TagListGroup: []*protoccacheserver.TagList{
								{
									TagList: []*protoccacheserver.Tag{
										{
											Key:         "dmpTagTest",
											TagType:     protoccacheserver.TagType_TAG_TYPE_DMP,
											Operator:    protoccacheserver.Operator_OPERATOR_FALSE,
											Value:       "dmpCodeTest3",
											DmpPlatform: 3,
											UnitIdType:  protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
										},
									},
								},
							},
						},
						UnitIdType: protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{}
			got, err := e.GetExperiments(tt.args.ctx, tt.args.projectID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExperiments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBody, _ := json.Marshal(got)
			wantBody, _ := json.Marshal(tt.want)
			assert.Equal(t, gotBody, wantBody)
			// if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetExperiments() got = \n%v, want \n%v", got, tt.want)
			// }
		})
	}
}

func BenchmarkDMPResultKeyFormat(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dmpTagResultKeyFormat(protoccacheserver.UnitIDType_UNIT_ID_TYPE_DEFAULT,
				3, "uuid.New().String()[:20]", &Options{
					UnitID:    "unitID",
					NewUnitID: "newUnitID",
				})
		}
	})
}

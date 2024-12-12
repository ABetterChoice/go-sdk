// Package testdata Test data Do not use
package testdata

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

/*
Unit test data rules
- layerKey uniformly starts with test and ends with the actual feature function being tested
For example, test-multiLayerDomain-overrideLayer identifies the overrideLayer layer of the test data under the multi-layer domain
- There are multiple domains under a domain, identified by suffixes such as 1, 2, and 3,
- Experimental version ID rules {{A}}{{B}}{{C}} ABC are 3 digits respectively 100_001_001 identifies the experimental version with ID 100001001 under the experiment with ID 100001 under layer ID 100
*/
var (
	projectID                          = "123"
	mockSuccessMsg                     = "mock success"
	subDomainMultiDomain1MultiLayer2   = "subDomain-multiDomain1-multiLayer2"
	subDomainMultiDomain1MultiLayer1   = "subDomain-multiDomain1-multiLayer1"
	subDomainHoldoutDomain1SingleLayer = "subDomain-holdoutDomain1-singleLayer"
	// EmptyTabConfig TODO
	// Deprecated: Test data Do not use
	EmptyTabConfig = &protoc_cache_server.TabConfig{
		ExperimentData: nil,
		ConfigData:     nil,
		ControlData:    nil,
	}
	// NormalTabConfig TODO
	// Deprecated: Test data Do not use
	NormalTabConfig = &protoc_cache_server.TabConfig{
		ExperimentData: &protoc_cache_server.ExperimentData{
			DefaultGroupId: -1,
			OverrideList: map[string]*protoc_cache_server.LayerToGroupID{
				"overrideID": {LayerToGroupId: map[string]int64{
					"overrideLayer": 100001001,
				}},
				"newOverrideUnitID": {LayerToGroupId: map[string]int64{
					"overrideLayer": 100001001,
				}},
			},
			GlobalDomain: &protoc_cache_server.Domain{
				Metadata: &protoc_cache_server.DomainMetadata{
					Key:        "globalDomain",
					DomainType: protoc_cache_server.DomainType_DOMAIN_TYPE_DOMAIN,
					HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
					HashSeed:   5080801,
					UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
					BucketSize: 100,
					TrafficRangeList: []*protoc_cache_server.TrafficRange{
						{
							Left:  1,
							Right: 100,
						},
					},
				},
				MultiLayerDomainList: []*protoc_cache_server.MultiLayerDomain{
					{
						Metadata: &protoc_cache_server.DomainMetadata{
							Key:        "multiDomain1",
							DomainType: protoc_cache_server.DomainType_DOMAIN_TYPE_MULTILAYER,
							HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:   2357911,
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
							BucketSize: 100,
							TrafficRangeList: []*protoc_cache_server.TrafficRange{
								{
									Left:  1,
									Right: 100,
								},
							}},
						LayerList: []*protoc_cache_server.Layer{
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key: "overrideLayer",
									DefaultGroup: &protoc_cache_server.Group{
										Id:            100001001,
										GroupKey:      "100001001",
										ExperimentId:  100001,
										ExperimentKey: "100001",
										Params: map[string]string{
											"key1": "100001001",
											// "variant_key_test": "exp-value",
										},
										IsDefault: true,
										IsControl: false,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									HashType:    protoc_cache_server.HashType_HASH_TYPE_SINGLE,
									HashMethod:  protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:    23579,
									SceneIdList: nil,
									UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize:  10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									100001001: { // Layer Default Experiment
										Id:            100001001,
										GroupKey:      "100001001",
										ExperimentId:  100001,
										ExperimentKey: "100001",
										Params: map[string]string{
											"key1": "100001001",
											// "variant_key_test": "exp-value",
										},
										IsDefault: true,
										IsControl: false,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									100002001: { // The control group of the experiment 100002
										Id:            100002001,
										GroupKey:      "100002001",
										ExperimentId:  100002,
										ExperimentKey: "100002",
										Params: map[string]string{
											"key1": "100002001",
										},
										IsDefault: false,
										IsControl: true,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									100002002: { // Experimental Group of the experiment 100002
										Id:            100002002,
										GroupKey:      "100002002",
										ExperimentId:  100002,
										ExperimentKey: "100002",
										Params: map[string]string{
											"key1": "100002002",
										},
										IsDefault: false,
										IsControl: false,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									100003001: { // The control group of the experiment 100003
										Id:            100003001,
										GroupKey:      "100003001",
										ExperimentId:  100003,
										ExperimentKey: "100003",
										Params: map[string]string{
											"key1": "100003001",
										},
										IsDefault: false,
										IsControl: true,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									100003002: { // The control group of the experiment 100003
										Id:            100003002,
										GroupKey:      "100003002",
										ExperimentId:  100003,
										ExperimentKey: "100003",
										Params: map[string]string{
											"key1": "100003002",
										},
										IsDefault: false,
										IsControl: false,
										LayerKey:  "overrideLayer",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								// Whether the index of the experiment TODO is verified when double hashing is used. Double hashing is required, and single hashing can be empty.
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									100001: &protoc_cache_server.Experiment{
										Id:  100001,
										Key: "100001",
										GroupIdIndex: map[int64]bool{
											100001001: true,
										},
									},
									100002: &protoc_cache_server.Experiment{
										Id:  100002,
										Key: "100002",
										GroupIdIndex: map[int64]bool{
											100002001: true,
											100002002: true,
										},
									},
									100003: &protoc_cache_server.Experiment{
										Id:  100003,
										Key: "100003",
										GroupIdIndex: map[int64]bool{
											100003001: true,
											100003002: true,
										},
									},
								},
							},
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key:          "multiLayer2",
									DefaultGroup: nil,
									HashType:     protoc_cache_server.HashType_HASH_TYPE_SINGLE,
									HashMethod:   protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:     23579,
									SceneIdList:  nil,
									UnitIdType:   protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize:   10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									101002001: {
										Id:            101002001,
										GroupKey:      "101002001",
										ExperimentId:  101002,
										ExperimentKey: "101002",
										Params: map[string]string{
											"key1": "101002001",
										},
										IsDefault: false,
										IsControl: true,
										LayerKey:  "multiLayer2",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									101002002: {
										Id:            101002002,
										GroupKey:      "101002002",
										ExperimentId:  101002,
										ExperimentKey: "101002",
										Params: map[string]string{
											"key1": "101002002",
										},
										IsDefault: false,
										IsControl: false,
										LayerKey:  "multiLayer2",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									101003001: {
										Id:            101003001,
										GroupKey:      "101003001",
										ExperimentId:  101003,
										ExperimentKey: "101003",
										Params: map[string]string{
											"key1": "101003001",
										},
										IsDefault: false,
										IsControl: true,
										LayerKey:  "multiLayer2",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									101003002: {
										Id:            101003002,
										GroupKey:      "101003002",
										ExperimentId:  101003,
										ExperimentKey: "101003",
										Params: map[string]string{
											"key1": "101003002",
										},
										IsDefault: false,
										IsControl: false,
										LayerKey:  "multiLayer2",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										SceneIdList: nil,
										UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								// Whether the index of the experiment TODO is verified when double hashing is used. Double hashing is required, and single hashing can be empty.
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									101002: {
										Id:  101002,
										Key: "101002",
										GroupIdIndex: map[int64]bool{
											101002001: true,
											101002002: true,
										},
									},
									101003: {
										Id:  101003,
										Key: "101003",
										GroupIdIndex: map[int64]bool{
											101003001: true,
											101003002: true,
										},
									},
								},
							},
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key:        "doubleHashLayerTag",
									HashType:   protoc_cache_server.HashType_HASH_TYPE_DOUBLE,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   5713,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									301001001: &protoc_cache_server.Group{
										Id:            301001001,
										GroupKey:      "301001001",
										ExperimentId:  301001,
										ExperimentKey: "301001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerTag",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "dmpTagTest",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
															Operator:    protoc_cache_server.Operator_OPERATOR_FALSE,
															Value:       "dmpCodeTest3",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									301001002: &protoc_cache_server.Group{
										Id:            301001002,
										GroupKey:      "301001002",
										ExperimentId:  301001,
										ExperimentKey: "301001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerTag",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "dmpTagTest",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
															Operator:    protoc_cache_server.Operator_OPERATOR_FALSE,
															Value:       "dmpCodeTest3",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									301001: &protoc_cache_server.Experiment{
										HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
										HashSeed:   5713,
										Id:         301001,
										Key:        "301001",
										BucketSize: 10000,
										IssueType:  protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
										GroupIdIndex: map[int64]bool{
											301001001: true,
											301001002: true,
										},
									},
								},
							},
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key:        "doubleHashLayerPercentage",
									HashType:   protoc_cache_server.HashType_HASH_TYPE_DOUBLE,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   5713,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									302001001: &protoc_cache_server.Group{
										Id:            302001001,
										GroupKey:      "302001001",
										ExperimentId:  302001,
										ExperimentKey: "302001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerPercentage",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									302001002: &protoc_cache_server.Group{
										Id:            302001002,
										GroupKey:      "302001002",
										ExperimentId:  302001,
										ExperimentKey: "302001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerPercentage",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									302001: &protoc_cache_server.Experiment{
										HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
										HashSeed:   5713,
										Id:         302001,
										Key:        "302001",
										BucketSize: 10000,
										IssueType:  protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
										GroupIdIndex: map[int64]bool{
											302001001: true,
											302001002: true,
										},
									},
								},
							},
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key:        "doubleHashLayerCityTag",
									HashType:   protoc_cache_server.HashType_HASH_TYPE_DOUBLE,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   5713,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									303001001: &protoc_cache_server.Group{
										Id:            303001001,
										GroupKey:      "303001001",
										ExperimentId:  303001,
										ExperimentKey: "303001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerCityTag",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_CITY_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "dmpTagTest",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
															Operator:    protoc_cache_server.Operator_OPERATOR_FALSE,
															Value:       "dmpCodeTest3",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									303001002: &protoc_cache_server.Group{
										Id:            303001002,
										GroupKey:      "303001002",
										ExperimentId:  303001,
										ExperimentKey: "303001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerCityTag",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_CITY_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "dmpTagTest",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
															Operator:    protoc_cache_server.Operator_OPERATOR_FALSE,
															Value:       "dmpCodeTestHit",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									303001: &protoc_cache_server.Experiment{
										HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
										HashSeed:   5713,
										Id:         303001,
										Key:        "303001",
										BucketSize: 10000,
										IssueType:  protoc_cache_server.IssueType_ISSUE_TYPE_CITY_TAG,
										GroupIdIndex: map[int64]bool{
											303001001: true,
											303001002: true,
										},
									},
								},
							},
							{
								Metadata: &protoc_cache_server.LayerMetadata{
									Key:        "doubleHashLayerTDMPagValue",
									HashType:   protoc_cache_server.HashType_HASH_TYPE_DOUBLE,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   5713,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 10000,
								},
								GroupIndex: map[int64]*protoc_cache_server.Group{
									304001001: &protoc_cache_server.Group{
										Id:            304001001,
										GroupKey:      "304001001",
										ExperimentId:  304001,
										ExperimentKey: "304001",
										IsControl:     true,
										LayerKey:      "doubleHashLayerTDMPagValue",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "tagTest123",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_STRING,
															Operator:    protoc_cache_server.Operator_OPERATOR_EQ,
															Value:       "123",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
															TagOrigin:   protoc_cache_server.TagOrigin_TAG_ORIGIN_DMP,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
									304001002: &protoc_cache_server.Group{
										Id:            304001002,
										GroupKey:      "304001002",
										ExperimentId:  304001,
										ExperimentKey: "304001",
										IsControl:     false,
										LayerKey:      "doubleHashLayerTDMPagValue",
										IssueInfo: &protoc_cache_server.IssueInfo{
											IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_CITY_TAG,
											TagListGroup: []*protoc_cache_server.TagList{
												{
													TagList: []*protoc_cache_server.Tag{
														{
															Key:         "tagTest123",
															TagType:     protoc_cache_server.TagType_TAG_TYPE_STRING,
															Operator:    protoc_cache_server.Operator_OPERATOR_EQ,
															Value:       "123",
															DmpPlatform: 3,
															UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
															TagOrigin:   protoc_cache_server.TagOrigin_TAG_ORIGIN_DMP,
														},
													},
												},
											},
										},
										UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									},
								},
								ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
									304001: &protoc_cache_server.Experiment{
										HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
										HashSeed:   5713,
										Id:         304001,
										Key:        "304001",
										BucketSize: 10000,
										IssueType:  protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
										GroupIdIndex: map[int64]bool{
											304001001: true,
											304001002: true,
										},
									},
								},
							},
						},
					},
				},
				DomainList: []*protoc_cache_server.Domain{
					&protoc_cache_server.Domain{
						Metadata: &protoc_cache_server.DomainMetadata{
							Key:        "subDomain1",
							DomainType: protoc_cache_server.DomainType_DOMAIN_TYPE_DOMAIN,
							HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:   5080803,
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
							BucketSize: 100,
							TrafficRangeList: []*protoc_cache_server.TrafficRange{
								{
									Left:  1,
									Right: 100,
								},
							},
						},
						HoldoutDomainList: []*protoc_cache_server.HoldoutDomain{
							{
								Metadata: &protoc_cache_server.DomainMetadata{
									Key:        "subDomain-holdoutDomain1",
									DomainType: protoc_cache_server.DomainType_DOMAIN_TYPE_HOLDOUT,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   5080804,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 100,
									TrafficRangeList: []*protoc_cache_server.TrafficRange{
										{
											Left:  1,
											Right: 5,
										},
									}},
								LayerList: []*protoc_cache_server.Layer{
									{
										Metadata: &protoc_cache_server.LayerMetadata{
											Key: subDomainHoldoutDomain1SingleLayer,
											DefaultGroup: &protoc_cache_server.Group{
												Id:            200001001,
												GroupKey:      "200001001",
												ExperimentId:  200001,
												ExperimentKey: "200001",
												Params: map[string]string{
													"key1": "200001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											HashType:    protoc_cache_server.HashType_HASH_TYPE_SINGLE,
											HashMethod:  protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
											HashSeed:    23579,
											SceneIdList: nil,
											UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											BucketSize:  10000,
										},
										GroupIndex: map[int64]*protoc_cache_server.Group{
											200001001: {
												Id:            200001001,
												GroupKey:      "200001001",
												ExperimentId:  200001,
												ExperimentKey: "200001",
												Params: map[string]string{
													"key1": "200001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											200002001: {
												Id:            200002001,
												GroupKey:      "200002001",
												ExperimentId:  200002,
												ExperimentKey: "200002",
												Params: map[string]string{
													"key1": "200002001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											200002002: {
												Id:            200002002,
												GroupKey:      "200002002",
												ExperimentId:  200002,
												ExperimentKey: "200002",
												Params: map[string]string{
													"key1": "200002002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											200003001: {
												Id:            200003001,
												GroupKey:      "200003001",
												ExperimentId:  200003,
												ExperimentKey: "200003",
												Params: map[string]string{
													"key1": "200003001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											200003002: {
												Id:            200003002,
												GroupKey:      "200003002",
												ExperimentId:  200003,
												ExperimentKey: "200003",
												Params: map[string]string{
													"key1": "200003002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainHoldoutDomain1SingleLayer,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
										},
										// Whether the index of the experiment TODO is verified when double hashing is used. Double hashing is required, and single hashing can be empty.
										ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
											200001: {
												Id:  200001,
												Key: "200001",
												GroupIdIndex: map[int64]bool{
													200001001: true,
												},
											},
											200002: {
												Id:  200002,
												Key: "200002",
												GroupIdIndex: map[int64]bool{
													200002001: true,
													200002002: true,
												},
											},
											200003: {
												Id:  200003,
												Key: "200003",
												GroupIdIndex: map[int64]bool{
													200003001: true,
													200003002: true,
												},
											},
										},
									},
								},
							},
						},
						MultiLayerDomainList: []*protoc_cache_server.MultiLayerDomain{
							{
								Metadata: &protoc_cache_server.DomainMetadata{
									Key:        "subDomain-multiDomain1",
									DomainType: protoc_cache_server.DomainType_DOMAIN_TYPE_MULTILAYER,
									HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
									HashSeed:   2357919,
									UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
									BucketSize: 100,
									TrafficRangeList: []*protoc_cache_server.TrafficRange{
										{
											Left:  6,
											Right: 100,
										},
									}},
								LayerList: []*protoc_cache_server.Layer{
									{
										Metadata: &protoc_cache_server.LayerMetadata{
											Key: subDomainMultiDomain1MultiLayer1,
											DefaultGroup: &protoc_cache_server.Group{
												Id:            201001001,
												GroupKey:      "201001001",
												ExperimentId:  201001,
												ExperimentKey: "201001",
												Params: map[string]string{
													"key1": "201001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											HashType:    protoc_cache_server.HashType_HASH_TYPE_SINGLE,
											HashMethod:  protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
											HashSeed:    23579,
											SceneIdList: nil,
											UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											BucketSize:  10000,
										},
										GroupIndex: map[int64]*protoc_cache_server.Group{
											201001001: {
												Id:            201001001,
												GroupKey:      "201001001",
												ExperimentId:  201001,
												ExperimentKey: "201001",
												Params: map[string]string{
													"key1": "201001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											201002001: {
												Id:            201002001,
												GroupKey:      "201002001",
												ExperimentId:  201002,
												ExperimentKey: "201002",
												Params: map[string]string{
													"key1": "201002001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											201002002: {
												Id:            201002002,
												GroupKey:      "201002002",
												ExperimentId:  201002,
												ExperimentKey: "201002",
												Params: map[string]string{
													"key1": "201002002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											201003001: {
												Id:            201003001,
												GroupKey:      "201003001",
												ExperimentId:  201003,
												ExperimentKey: "201003",
												Params: map[string]string{
													"key1": "201003001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest1",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											201003002: {
												Id:            201003002,
												GroupKey:      "201003002",
												ExperimentId:  201003,
												ExperimentKey: "201003",
												Params: map[string]string{
													"key1": "201003002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer1,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest1",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
																&protoc_cache_server.Tag{
																	Key:        "numberTagTest",
																	TagType:    protoc_cache_server.TagType_TAG_TYPE_NUMBER,
																	Operator:   protoc_cache_server.Operator_OPERATOR_GTE,
																	Value:      "95",
																	UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
										},
										// Whether the index of the experiment TODO is verified when double hashing is used. Double hashing is required, and single hashing can be empty.
										ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
											201001: {
												Id:  201001,
												Key: "201001",
												GroupIdIndex: map[int64]bool{
													201001001: true,
												},
											},
											201002: {
												Id:  201002,
												Key: "201002",
												GroupIdIndex: map[int64]bool{
													201002001: true,
													201002002: true,
												},
											},
											201003: {
												Id:  201003,
												Key: "201003",
												GroupIdIndex: map[int64]bool{
													201003001: true,
													201003002: true,
												},
											},
										},
									},
									{
										Metadata: &protoc_cache_server.LayerMetadata{
											Key: subDomainMultiDomain1MultiLayer2,
											DefaultGroup: &protoc_cache_server.Group{
												Id:            202001001,
												GroupKey:      "202001001",
												ExperimentId:  202001,
												ExperimentKey: "202001",
												Params: map[string]string{
													"key1": "202001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											HashType:    protoc_cache_server.HashType_HASH_TYPE_DOUBLE,
											HashMethod:  protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
											HashSeed:    23579,
											SceneIdList: nil,
											UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											BucketSize:  10000,
										},
										GroupIndex: map[int64]*protoc_cache_server.Group{
											202001001: {
												Id:            202001001,
												GroupKey:      "202001001",
												ExperimentId:  202001,
												ExperimentKey: "202001",
												Params: map[string]string{
													"key1": "202001001",
												},
												IsDefault: true,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											202002001: {
												Id:            202002001,
												GroupKey:      "201002001",
												ExperimentId:  202002,
												ExperimentKey: "202002",
												Params: map[string]string{
													"key1": "202002001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest2",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
																&protoc_cache_server.Tag{
																	Key:        "numberTagTest",
																	TagType:    protoc_cache_server.TagType_TAG_TYPE_NUMBER,
																	Operator:   protoc_cache_server.Operator_OPERATOR_GTE,
																	Value:      "95",
																	UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											202002002: {
												Id:            202002002,
												GroupKey:      "202002002",
												ExperimentId:  202002,
												ExperimentKey: "202002",
												Params: map[string]string{
													"key1": "202002002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest2",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
																&protoc_cache_server.Tag{
																	Key:        "numberTagTest",
																	TagType:    protoc_cache_server.TagType_TAG_TYPE_NUMBER,
																	Operator:   protoc_cache_server.Operator_OPERATOR_GTE,
																	Value:      "95",
																	UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											202003001: {
												Id:            202003001,
												GroupKey:      "202003001",
												ExperimentId:  202003,
												ExperimentKey: "202003",
												Params: map[string]string{
													"key1": "202003001",
												},
												IsDefault: false,
												IsControl: true,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest1",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
											202003002: {
												Id:            202003002,
												GroupKey:      "202003002",
												ExperimentId:  202003,
												ExperimentKey: "202003",
												Params: map[string]string{
													"key1": "202003002",
												},
												IsDefault: false,
												IsControl: false,
												LayerKey:  subDomainMultiDomain1MultiLayer2,
												IssueInfo: &protoc_cache_server.IssueInfo{
													IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
													TagListGroup: []*protoc_cache_server.TagList{
														{
															TagList: []*protoc_cache_server.Tag{
																&protoc_cache_server.Tag{
																	Key:         "dmpTest1",
																	TagType:     protoc_cache_server.TagType_TAG_TYPE_DMP,
																	Operator:    protoc_cache_server.Operator_OPERATOR_TRUE,
																	Value:       "dmpCodeTest1",
																	DmpPlatform: 3,
																	UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
																&protoc_cache_server.Tag{
																	Key:        "numberTagTest",
																	TagType:    protoc_cache_server.TagType_TAG_TYPE_NUMBER,
																	Operator:   protoc_cache_server.Operator_OPERATOR_GTE,
																	Value:      "95",
																	UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
																},
															},
														},
													},
												},
												SceneIdList: nil,
												UnitIdType:  protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
										},
										// Whether the index of the experiment TODO is verified when double hashing is used. Double hashing is required, and single hashing can be empty.
										ExperimentIndex: map[int64]*protoc_cache_server.Experiment{
											202001: {
												Id:         202001,
												Key:        "202001",
												HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
												HashSeed:   257913,
												BucketSize: 10000,
												GroupIdIndex: map[int64]bool{
													202001001: true,
												},
											},
											202002: {
												Id:         202002,
												Key:        "202002",
												HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
												HashSeed:   257917,
												BucketSize: 10000,
												GroupIdIndex: map[int64]bool{
													202002001: true,
													202002002: true,
												},
											},
											202003: {
												Id:         202003,
												Key:        "202003",
												HashMethod: protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
												HashSeed:   257919,
												BucketSize: 10000,
												GroupIdIndex: map[int64]bool{
													202003001: true,
													202003002: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ConfigData: &protoc_cache_server.RemoteConfigData{
			RemoteConfigIndex: map[string]*protoc_cache_server.RemoteConfig{
				"remoteConfig1": &protoc_cache_server.RemoteConfig{
					Key:          "remoteConfig1",
					DefaultValue: []byte("remoteConfig1-defaultValue"),
					Version:      "v0.1.0",
					SceneIdList:  []int64{1, 2, 3},
					Type:         protoc_cache_server.RemoteConfigValueType_REMOTE_CONFIG_VALUE_TYPE_BYTES,
					OverrideList: map[string][]byte{
						"overrideUnitID": []byte("hitOverrideResult"),
					},
					ConditionList: []*protoc_cache_server.Condition{
						&protoc_cache_server.Condition{
							Id:            1,
							Key:           "condition1",
							Value:         []byte("remoteConfig1-condition1"),
							HashMethod:    protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:      23579,
							ExperimentKey: "",
							BucketSize:    10000,
							BucketInfo: &protoc_cache_server.BucketInfo{
								BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
								TrafficRange: &protoc_cache_server.TrafficRange{
									Left:  1,
									Right: 10000,
								},
								Version:    "",
								ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
							},
							IssueInfo: &protoc_cache_server.IssueInfo{
								IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
							},
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
				"bitmapTest": &protoc_cache_server.RemoteConfig{
					Key:          "bitmapTest",
					DefaultValue: []byte("bitmapTestDefaultValue"),
					Version:      "v0.2.1",
					SceneIdList:  []int64{5, 6, 7},
					Type:         protoc_cache_server.RemoteConfigValueType_REMOTE_CONFIG_VALUE_TYPE_BYTES,
					OverrideList: map[string][]byte{
						"overrideUnitID": []byte("hitOverrideResult"),
					},
					ConditionList: []*protoc_cache_server.Condition{
						&protoc_cache_server.Condition{
							Id:            1,
							Key:           "condition1",
							Value:         []byte("bitmapTest-condition1"),
							HashMethod:    protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:      23579,
							ExperimentKey: "",
							BucketSize:    10000,
							BucketInfo: &protoc_cache_server.BucketInfo{
								BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
								Version:    "",
								Bitmap:     MockGenBitmap(1, 10000),
								ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
							},
							IssueInfo: &protoc_cache_server.IssueInfo{
								IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
							},
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
				"withTag": &protoc_cache_server.RemoteConfig{
					Key:          "withTag",
					DefaultValue: []byte("withTagDefaultValue"),
					Version:      "v0.2.1",
					SceneIdList:  nil,
					Type:         protoc_cache_server.RemoteConfigValueType_REMOTE_CONFIG_VALUE_TYPE_BYTES,
					OverrideList: map[string][]byte{
						"overrideUnitID": []byte("hitOverrideResult"),
					},
					ConditionList: []*protoc_cache_server.Condition{
						&protoc_cache_server.Condition{
							Id:            1,
							Key:           "condition1",
							Value:         []byte("withTag-condition1"),
							HashMethod:    protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:      23579,
							ExperimentKey: "",
							BucketSize:    10000,
							BucketInfo: &protoc_cache_server.BucketInfo{
								BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
								TrafficRange: &protoc_cache_server.TrafficRange{
									Left:  1,
									Right: 10000,
								},
								Version:    "",
								ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
							},
							IssueInfo: &protoc_cache_server.IssueInfo{
								IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_TAG,
								TagListGroup: []*protoc_cache_server.TagList{
									{
										TagList: []*protoc_cache_server.Tag{
											&protoc_cache_server.Tag{
												Key:        "tagKey1",
												TagType:    protoc_cache_server.TagType_TAG_TYPE_STRING,
												Operator:   protoc_cache_server.Operator_OPERATOR_EQ,
												Value:      "ios",
												UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
											},
										},
									},
								},
							},
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
				"withExperiment": &protoc_cache_server.RemoteConfig{
					Key:          "withExperiment",
					DefaultValue: []byte("withExperimentDefaultValue"),
					Version:      "v0.3.1",
					SceneIdList:  nil,
					Type:         protoc_cache_server.RemoteConfigValueType_REMOTE_CONFIG_VALUE_TYPE_BYTES,
					OverrideList: map[string][]byte{
						"overrideUnitID": []byte("hitOverrideResult"),
					},
					ConditionList: []*protoc_cache_server.Condition{
						&protoc_cache_server.Condition{
							Id:            1,
							Key:           "condition1",
							Value:         []byte("withExperiment-condition1"),
							HashMethod:    protoc_cache_server.HashMethod_HASH_METHOD_BKDR,
							HashSeed:      23579,
							ExperimentKey: "303001",
							BucketSize:    10000,
							BucketInfo: &protoc_cache_server.BucketInfo{
								BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
								TrafficRange: &protoc_cache_server.TrafficRange{
									Left:  1,
									Right: 10000,
								},
								Version:    "",
								ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
							},
							IssueInfo: &protoc_cache_server.IssueInfo{
								IssueType: protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE,
							},
							UnitIdType: protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT,
						},
					},
				},
			}},
		ControlData: &protoc_cache_server.ControlData{
			RefreshInterval:     3, // 3s
			IgnoreReportGroupId: map[int64]bool{100000: true},
			ExperimentMetricsConfig: map[int64]*protoc_cache_server.MetricsConfig{
				1: &protoc_cache_server.MetricsConfig{
					IsAutomatic:         true,
					IsEnable:            true,
					PluginName:          "empty",
					SamplingInterval:    100,
					ErrSamplingInterval: 1,
					Metadata: &protoc_cache_server.MetricsMetadata{
						ExpandedData: nil,
						Name:         "empty",
						Id:           "empty",
					},
				},
			},
			DefaultExperimentMetricsConfig: &protoc_cache_server.MetricsConfig{
				IsAutomatic:         true,
				IsEnable:            true,
				PluginName:          "pubsub",
				SamplingInterval:    1,
				ErrSamplingInterval: 1,
				Metadata: &protoc_cache_server.MetricsMetadata{
					ExpandedData: nil,
					Name:         "abc_exp_expose_test",
					Id:           "empty",
				},
			},
			RemoteConfigMetricsConfig: map[int64]*protoc_cache_server.MetricsConfig{
				1: &protoc_cache_server.MetricsConfig{
					IsAutomatic:         true,
					IsEnable:            true,
					PluginName:          "empty",
					SamplingInterval:    100,
					ErrSamplingInterval: 1,
					Metadata: &protoc_cache_server.MetricsMetadata{
						ExpandedData: nil,
						Name:         "empty",
						Id:           "empty",
					},
				},
			},
			DefaultRemoteConfigMetricsConfig: &protoc_cache_server.MetricsConfig{
				IsAutomatic:         true,
				IsEnable:            true,
				PluginName:          "empty",
				SamplingInterval:    100,
				ErrSamplingInterval: 1,
				Metadata: &protoc_cache_server.MetricsMetadata{
					ExpandedData: nil,
					Name:         "empty",
					Id:           "empty",
				},
			},
			FeatureFlagMetricsConfig: map[int64]*protoc_cache_server.MetricsConfig{
				1: &protoc_cache_server.MetricsConfig{
					IsAutomatic:         true,
					IsEnable:            true,
					PluginName:          "empty",
					SamplingInterval:    100,
					ErrSamplingInterval: 1,
					Metadata: &protoc_cache_server.MetricsMetadata{
						ExpandedData: nil,
						Name:         "empty",
						Id:           "empty",
					},
				},
			},
			DefaultFeatureFlagMetricsConfig: &protoc_cache_server.MetricsConfig{
				IsAutomatic:         true,
				IsEnable:            true,
				PluginName:          "empty",
				SamplingInterval:    100,
				ErrSamplingInterval: 1,
				Metadata: &protoc_cache_server.MetricsMetadata{
					ExpandedData: nil,
					Name:         "empty",
					Id:           "empty",
				},
			},
			EventMetricsConfig: &protoc_cache_server.MetricsConfig{
				IsAutomatic:         true,
				IsEnable:            true,
				PluginName:          "empty",
				SamplingInterval:    1,
				ErrSamplingInterval: 1,
				Metadata: &protoc_cache_server.MetricsMetadata{
					ExpandedData: nil,
					Name:         "empty",
					Id:           "empty",
				},
			},
			// MetricsInitConfigIndex: map[string]*protoc_cache_server.MetricsInitConfig{
			//	"pubsub": &protoc_cache_server.MetricsInitConfig{
			//		Kv: map[string]string{
			//			metrics.CredentialsJSON: "",
			//			metrics.TopicNameKey:    "abc_exp_expose_test",
			//			metrics.ProjectIDKey:    "abetterchoice",
			//		},
			//	},
			// },
		},
	}
)

var (
	// NormalExperimentBucketInfo TODO
	// Deprecated: for test,do not use
	NormalExperimentBucketInfo = map[int64]*protoc_cache_server.BucketInfo{
		202001: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1,
				Right: 1000,
			},
			Bitmap:     nil,
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202002: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1001,
				Right: 5000,
			},
			Bitmap:     nil,
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202003: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			Bitmap:     MockGenBitmap(6001, 10000),
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202004: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  6001,
				Right: 10000,
			},
			Bitmap:     nil,
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_DELETE,
		},
		202005: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  6001,
				Right: 10000,
			},
			Bitmap:     nil,
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		301001: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		302001: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		303001: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		304001: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  1,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
	}
	// NormalGroupBucketInfo TODO
	// Deprecated: for test,do not use
	NormalGroupBucketInfo = map[int64]*protoc_cache_server.BucketInfo{
		100001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 0, Right: 0},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		100002001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1, Right: 1000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		100002002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1001, Right: 2000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		100003001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 2001, Right: 3000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		100003002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 3001, Right: 4000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		101001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 0, Right: 0},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		101002001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1, Right: 1000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		101002002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1001, Right: 2000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		101003001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 2001, Right: 3000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		101003002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 3001, Right: 4000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		200001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 0, Right: 0},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		200002001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1, Right: 1000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		200002002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1001, Right: 2000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		200003001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 2001, Right: 3000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		200003002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 3001, Right: 4000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 0, Right: 0},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201002001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1, Right: 1000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201002002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1001, Right: 2000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201003001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(2001, 3000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201003002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(3001, 10000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		201004001: &protoc_cache_server.BucketInfo{ // 已下线实验组
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(3001, 4000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_DELETE,
		},
		202001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 0, Right: 0},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202002001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 1, Right: 5000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202002002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{Left: 5001, Right: 10000},
			Bitmap:       nil,
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202003001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(1001, 7000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		202003002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(7001, 9000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		301001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(0, 5000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		301001002: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(5001, 10000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		302001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(0, 5000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		302001002: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  5001,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		303001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(0, 5000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		303001002: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  5001,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		304001001: &protoc_cache_server.BucketInfo{
			BucketType:   protoc_cache_server.BucketType_BUCKET_TYPE_BITMAP,
			TrafficRange: nil,
			Bitmap:       MockGenBitmap(0, 5000),
			Version:      "",
			ModifyType:   protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
		304001002: &protoc_cache_server.BucketInfo{
			BucketType: protoc_cache_server.BucketType_BUCKET_TYPE_RANGE,
			TrafficRange: &protoc_cache_server.TrafficRange{
				Left:  5001,
				Right: 10000,
			},
			Version:    "",
			ModifyType: protoc_cache_server.ModifyType_MODIFY_UPDATE,
		},
	}
)

// MockGenBitmap for test
// Deprecated: for test
func MockGenBitmap(left, right uint64) []byte {
	bitmap := roaring.NewBitmap()
	bitmap.AddRange(left, right)
	bitmapByte, _ := bitmap.MarshalBinary()
	return bitmapByte
}

// MockCacheClient TODO
// Deprecated: MockCacheClient
func MockCacheClient(t gomock.TestReporter) *client.MockClient {
	mockClient := client.NewMockClient(gomock.NewController(t))
	// mock GetTabConfigData
	mockClient.EXPECT().GetTabConfigData(gomock.Any(), &protoc_cache_server.GetTabConfigReq{
		ProjectId:  projectID,
		UpdateType: protoc_cache_server.UpdateType_UPDATE_TYPE_COMPLETE,
		SdkVersion: env.SDKVersion,
	}).Return(&protoc_cache_server.GetTabConfigResp{
		Code:    protoc_cache_server.Code_CODE_SUCCESS,
		Message: mockSuccessMsg,
		TabConfigManager: &protoc_cache_server.TabConfigManager{
			ProjectId:  projectID,
			Version:    "",
			UpdateType: protoc_cache_server.UpdateType_UPDATE_TYPE_COMPLETE,
			TabConfig:  NormalTabConfig,
		},
	}, nil).AnyTimes()
	// mock BatchGetExperimentBucketInfo
	mockClient.EXPECT().BatchGetExperimentBucketInfo(gomock.Any(), gomock.Any()).
		Return(&protoc_cache_server.BatchGetExperimentBucketResp{
			Code:        protoc_cache_server.Code_CODE_SUCCESS,
			Message:     mockSuccessMsg,
			BucketIndex: NormalExperimentBucketInfo,
		}, nil).AnyTimes()
	mockClient.EXPECT().BatchGetGroupBucketInfo(gomock.Any(), gomock.Any()).Return(&protoc_cache_server.BatchGetGroupBucketResp{
		Code:        protoc_cache_server.Code_CODE_SUCCESS,
		Message:     mockSuccessMsg,
		BucketIndex: NormalGroupBucketInfo}, nil).AnyTimes()
	return mockClient
}

// MockCacheClientWithData TODO
// Deprecated: MockCacheClientWithData
func MockCacheClientWithData(t gomock.TestReporter, tabConfig *protoc_cache_server.TabConfig,
	experimentBucketInfoIndex, groupBucketInfoIndex map[int64]*protoc_cache_server.BucketInfo) *client.MockClient {
	mockClient := client.NewMockClient(gomock.NewController(t))
	// mock GetTabConfigData
	mockClient.EXPECT().GetTabConfigData(gomock.Any(), &protoc_cache_server.GetTabConfigReq{
		ProjectId:  projectID,
		UpdateType: protoc_cache_server.UpdateType_UPDATE_TYPE_COMPLETE,
		SdkVersion: env.SDKVersion,
	}).Return(&protoc_cache_server.GetTabConfigResp{
		Code:    protoc_cache_server.Code_CODE_SUCCESS,
		Message: mockSuccessMsg,
		TabConfigManager: &protoc_cache_server.TabConfigManager{
			ProjectId:  projectID,
			Version:    "",
			UpdateType: protoc_cache_server.UpdateType_UPDATE_TYPE_COMPLETE,
			TabConfig:  tabConfig,
		},
	}, nil).AnyTimes()
	// mock BatchGetExperimentBucketInfo
	mockClient.EXPECT().BatchGetExperimentBucketInfo(gomock.Any(), gomock.Any()).
		Return(&protoc_cache_server.BatchGetExperimentBucketResp{
			Code:        protoc_cache_server.Code_CODE_SUCCESS,
			Message:     mockSuccessMsg,
			BucketIndex: experimentBucketInfoIndex,
		}, nil).AnyTimes()
	mockClient.EXPECT().BatchGetGroupBucketInfo(gomock.Any(), gomock.Any()).Return(&protoc_cache_server.BatchGetGroupBucketResp{
		Code:        protoc_cache_server.Code_CODE_SUCCESS,
		Message:     mockSuccessMsg,
		BucketIndex: groupBucketInfoIndex}, nil).AnyTimes()
	return mockClient
}

// MockFakeCacheClient mock fake cacheClient
// Deprecated: for test
func MockFakeCacheClient(t *testing.T) *client.MockClient {
	mockClient := client.NewMockClient(gomock.NewController(t))
	mockClient.EXPECT().GetTabConfigData(gomock.Any(), gomock.Any()).Return(nil, errors.Errorf("mock err")).AnyTimes()
	return mockClient
}

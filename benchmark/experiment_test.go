// Package benchmark ...
package benchmark

//
// import (
//	"context"
//	"testing"
//
//	tab "github.com/abetterchoice/go-sdk"
//	"github.com/abetterchoice/go-sdk/testdata"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
// )
//
// func Init(b *testing.B) {
//	// log.SetLoggerLevel(log.ErrorLevel)
//	// cacheOldData := &protoc_sdk_api_server.GetTabConfigResp{}
//	// body, err := os.ReadFile("../testdata/5608_test.data")
//	// assert.Nil(b, err)
//	// err = json.Unmarshal(body, cacheOldData)
//	// assert.Nil(b, err)
//	//
//	// result := transform(cacheOldData)
//	// //t.Logf("%+v", result)
//	//
//	// testdata.NormalTabConfig = result.TabConfigManager.TabConfig
// }
//
// // BenchmarkUUID
// func BenchmarkUUID(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			_ = uuid.NewString()
//		}
//	})
// }
//
//
// func BenchmarkGetExperiments(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID)
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// for layerKey, experiment := range experiments.Data {
//			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
//			// }
//		}
//	})
// }
//
//
// func BenchmarkGetExpsFullFlowLayerKey(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID, tab.WithLayerKey("overrideLayer"))
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// var result = make(map[int64]int64)
//			// for _, experiment := range experiments.Data {
//			//	result[experiment.ID]++
//			// }
//			// b.Logf("%+v", result)
//		}
//	})
// }
//
//
// func BenchmarkGetExpsNotFFLayerKey(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID, tab.WithLayerKeyList([]string{"subDomain-holdoutDomain1-singleLayer"}))
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// var result = make(map[int64]int64)
//			// for _, experiment := range experiments.Data {
//			//	result[experiment.ID]++
//			// }
//			// b.Logf("%+v", result)
//		}
//	})
// }
//
//
// func BenchmarkGetExpsFFLayerKeyAndID(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext("overrideID")
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID, tab.WithLayerKey("overrideLayer"))
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// for layerKey, experiment := range experiments.Data {
//			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
//			// }
//		}
//	})
// }
//
//
// func BenchmarkGetExperimentWithoutReport(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
//		tab.WithDisableReport(true))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID)
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// for layerKey, experiment := range experiments.Data {
//			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
//			// }
//		}
//	})
// }
//
//
// func BenchmarkGetExperimentWithoutDMP(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID, tab.WithIsDisableDMP(true))
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// for layerKey, experiment := range experiments.Data {
//			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
//			// }
//		}
//	})
// }
//
//
// func BenchmarkGetExpReportAndWithoutDMP(b *testing.B) {
//	defer tab.Release()
//	Init(b)
//	mockClient := testdata.MockCacheClient(b)
//	err := tab.Init(context.Background(), projectIDList,
//		tab.WithRegisterCacheClient(mockClient),
//		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
//		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
//		tab.WithDisableReport(true))
//	assert.Nil(b, err)
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			userCtx := tab.NewUserContext(uuid.NewString())
//			experiments, err := userCtx.GetExperiments(context.TODO(), projectID, tab.WithIsDisableDMP(true))
//			if err != nil {
//				b.Fatalf("GetExperiments fail:%v", err)
//				return
//			}
//			assert.NotNil(b, experiments)
//			// for layerKey, experiment := range experiments.Data {
//			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
//			// }
//		}
//	})
// }
//
// var (
//	projectIDList = []string{"123"}
//	projectID     = "123"
// )

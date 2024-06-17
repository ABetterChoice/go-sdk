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
//	// // 转换成新协议数据结构
//	// result := transform(cacheOldData)
//	// //t.Logf("%+v", result)
//	// // 对同一份数据不同数据格式进行新老 SDK 版本性能测试
//	// testdata.NormalTabConfig = result.TabConfigManager.TabConfig
// }
//
// // BenchmarkUUID 基线
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
// // 进行实验分流基准测试
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
// // 指定满流量层做实验分流
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
// // 指定非满流量层
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
// // 白名单，指定满流量层做实验分流
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
// // 屏蔽上报
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
// // 屏蔽 dmp
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
// // 屏蔽上报
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

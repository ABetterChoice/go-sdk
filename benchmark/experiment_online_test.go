// Package benchmark ...
package benchmark

import (
	"context"
	"testing"

	tab "github.com/abetterchoice/go-sdk"
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/testdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var onlineProjectIDList = []string{"101"}
var onlineProjectID = "101"

// BenchmarkOnlineUUID 基线
func BenchmarkOnlineUUID(b *testing.B) {
	env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	defer tab.Release()
	err := tab.Init(context.Background(), onlineProjectIDList,
		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
	assert.Nil(b, err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.NewString()
		}
	})
}

// 进行实验分流基准测试
func BenchmarkOnlineGetExperiments(b *testing.B) {
	env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	defer tab.Release()
	err := tab.Init(context.Background(), onlineProjectIDList,
		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
	assert.Nil(b, err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userCtx := tab.NewUserContext(uuid.NewString())
			experiments, err := userCtx.GetExperiments(context.TODO(), onlineProjectID)
			if err != nil {
				b.Fatalf("GetExperiments fail:%v", err)
				return
			}
			assert.NotNil(b, experiments)
			// for layerKey, experiment := range experiments.Data {
			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
			// }
		}
	})
}

// 屏蔽上报
func BenchmarkOnlineGetExperimentWithoutReport(b *testing.B) {
	env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	defer tab.Release()
	err := tab.Init(context.Background(), onlineProjectIDList,
		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
		tab.WithDisableReport(true))
	assert.Nil(b, err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userCtx := tab.NewUserContext(uuid.NewString())
			experiments, err := userCtx.GetExperiments(context.TODO(), onlineProjectID)
			if err != nil {
				b.Fatalf("GetExperiments fail:%v", err)
				return
			}
			assert.NotNil(b, experiments)
			// for layerKey, experiment := range experiments.Data {
			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
			// }
		}
	})
}

// 屏蔽 dmp
func BenchmarkOnlineGetExperimentWithoutDMP(b *testing.B) {
	env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	defer tab.Release()
	err := tab.Init(context.Background(), onlineProjectIDList,
		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil))
	assert.Nil(b, err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userCtx := tab.NewUserContext(uuid.NewString())
			experiments, err := userCtx.GetExperiments(context.TODO(), onlineProjectID, tab.WithIsDisableDMP(true))
			if err != nil {
				b.Fatalf("GetExperiments fail:%v", err)
				return
			}
			assert.NotNil(b, experiments)
			// for layerKey, experiment := range experiments.Data {
			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
			// }
		}
	})
}

// 屏蔽上报
func BenchmarkOnlineGetExpReportAndWithoutDMP(b *testing.B) {
	env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	defer tab.Release()
	err := tab.Init(context.Background(), onlineProjectIDList,
		tab.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		tab.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
		tab.WithDisableReport(true))
	assert.Nil(b, err)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userCtx := tab.NewUserContext(uuid.NewString())
			experiments, err := userCtx.GetExperiments(context.TODO(), onlineProjectID, tab.WithIsDisableDMP(true))
			if err != nil {
				b.Fatalf("GetExperiments fail:%v", err)
				return
			}
			assert.NotNil(b, experiments)
			// for layerKey, experiment := range experiments.Data {
			//	b.Logf("[layerKey=%v]%+v", layerKey, experiment)
			// }
		}
	})
}

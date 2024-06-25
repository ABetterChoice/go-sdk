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

// BenchmarkOnlineUUID Baseline
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

// Conducting experimental offload benchmarks
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

// Without report
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

// without dmp
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

// Without DMP
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

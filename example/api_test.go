// Package example ...
package example

import (
	"context"
	"encoding/base64"
	"github.com/abetterchoice/go-sdk/plugin/log"
	"net/http"
	"net/http/pprof"
	"sync"
	"testing"
	"time"

	sdk "github.com/abetterchoice/go-sdk"
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/testdata"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTGLog(t *testing.T) {
	t.Skip()
	pID := "36"
	// env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
	// Initialize the sdk. The global initialization only needs to be done once.
	log.SetLoggerLevel(log.ErrorLevel)
	err := sdk.Init(context.TODO(), []string{pID},
		sdk.WithEnvType(env.TypePrd),
		sdk.WithSecretKey(""),
		// sdk.WithRegisterMetricsPlugin(tglog.Client, &protoc_cache_server.MetricsInitConfig{
		//	Region: "",
		//	Addr:   "127.0.0.1:8888",
		//	Kv:     nil,
		// }),
	)
	if err != nil {
		t.Fatalf("init fail:%v", err)
	}
	// sdk.Release()
	time.Sleep(20 * time.Second)
}

func TestInit(t *testing.T) {
	t.Skip()
	for {
		pID := "1111"
		// Initialize the sdk. The global initialization only needs to be done once.
		err := sdk.Init(context.TODO(), []string{pID},
			sdk.WithEnvType(env.TypePrd),
			sdk.WithSecretKey(""))
		// sdk.WithRegisterCacheClient(testdata.MockCacheClient(t)),
		if err != nil {
			t.Fatalf("init fail:%v", err)
		}
		sdk.Release()
		time.Sleep(2 * time.Second)
	}
}

func TestMonitorEvent(t *testing.T) {
	g := &protoc_event_server.MonitorEventGroup{}
	for i := 0; i < 5; i++ {
		g.Events = append(g.Events, &protoc_event_server.MonitorEvent{
			Time:      time.Now().Unix(),
			Ip:        "127.0.0.1",
			ProjectId: "6666",
			EventName: "" +
				"init",
			Latency:    12312 + float32(i),
			StatusCode: 0,
			Message:    "success",
			SdkType:    env.SDKType,
			SdkVersion: env.Version,
			InvokePath: "/root",
			InputData:  "",
			OutputData: "",
			ExtInfo:    nil,
		})
	}
	data, err := proto.Marshal(g)
	if err != nil {
		t.Fatalf("proto marshal fail:%v", err)
	}
	t.Logf("%s", base64.StdEncoding.EncodeToString(data))
}

var (
	projectID = "20001"
)

func TestSDKV2(t *testing.T) {
	t.Skip()
	pID := "6666"
	defer sdk.Release()
	// Initialize the sdk. The global initialization only needs to be done once.
	err := sdk.Init(context.TODO(), []string{pID},
		sdk.WithEnvType(env.TypePrd),
		sdk.WithSecretKey(""))
	// sdk.WithRegisterCacheClient(testdata.MockCacheClient(t)),
	assert.Nil(t, err)
	// Build an API client, userCtx is equivalent to the abstraction of a user
	userCtx := sdk.NewUserContext("21314212312eresadf",
		sdk.WithTags(map[string][]string{
			"country":  []string{"cn"},
			"strtest":  []string{"rthdfgh2"},
			"booltest": []string{"false"},
		}), sdk.WithTagKV("age", "18"))
	t.Run("abc", func(t *testing.T) {
		exp, err := userCtx.GetFeatureFlag(context.TODO(), pID, "yancisong_test_ff_1")
		if err != nil {
			t.Fatalf("getExperiment fail:%v", err)
		}
		t.Logf("%+v", exp.IsOverrideList)
		select {}
	})
}

func TestSDK(t *testing.T) {
	t.Skip()
	defer sdk.Release()
	// Initialize the sdk. The global initialization only needs to be done once.
	err := sdk.Init(context.TODO(), []string{projectID},
		sdk.WithEnvType(env.TypePrd),
		sdk.WithSecretKey(""))
	// sdk.WithRegisterCacheClient(testdata.MockCacheClient(t)),
	assert.Nil(t, err)
	userCtx := sdk.NewUserContext("sfjklajf45",
		sdk.WithTags(map[string][]string{
			"country":  []string{"cn"},
			"strtest":  []string{"rthdfgh2"},
			"booltest": []string{"false"},
		}), sdk.WithTagKV("age", "18"))
	t.Run("abc", func(t *testing.T) {
		exp, err := userCtx.GetExperiments(context.TODO(), projectID, sdk.WithAutomatic(false))
		if err != nil {
			t.Fatalf("getExperiment fail:%v", err)
		}
		err = sdk.LogExperimentsExposure(context.TODO(), projectID, exp)
		if err != nil {
			t.Fatalf("LogExperimentsExposure fail:%v", err)
		}
		t.Logf("%+v", exp)
		select {}
	})
	t.Run("experiment", func(t *testing.T) {
		group, err := userCtx.GetExperiment(context.TODO(), projectID, "overrideLayer")
		assert.Nil(t, err)
		assert.Equal(t, int64(0), group.MustGetInt64("color"))
		assert.Equal(t, true, group.GetBoolWithDefault("switch", true))
		// exposure
		sdk.LogExperimentExposure(context.TODO(), projectID, group)
		assert.Nil(t, err)
	})
	t.Run("experiments", func(t *testing.T) {
		experimentList, err := userCtx.GetExperiments(context.TODO(), "6666", sdk.WithSceneID(6))
		assert.Nil(t, err)
		for _, experiment := range experimentList.Data {
			assert.Equal(t, experiment.MustGetBool("switch"), false)
		}
		// exposure
		assert.Nil(t, err)
	})
	t.Run("remoteConfig", func(t *testing.T) {
		config, err := userCtx.GetRemoteConfig(context.TODO(), "6666", "remoteConfigKey")
		assert.Nil(t, err)
		assert.Equal(t, nil, config.Byte())
		assert.Equal(t, 100, config.GetInt64WithDefault(100))
		// exposure
		assert.Nil(t, err)
	})
	t.Run("featureFlag", func(t *testing.T) {
		config, err := userCtx.GetFeatureFlag(context.TODO(), projectID, "remoteConfigKey")
		assert.Nil(t, err)
		assert.Equal(t, false, config.MustGetBool())
		assert.Nil(t, err)
	})
}

func TestExperimentPprof(t *testing.T) {
	t.Skip()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})

		http.ListenAndServe(":6066", mux)
		// panic("fail")
	}()
	times := 150
	defer sdk.Release()
	// testdata.NormalTabConfig.ControlData.EventMetricsConfig.IsEnable = false
	err := sdk.Init(context.TODO(), []string{"6666"},
		sdk.WithRegisterMetricsPlugin(testdata.EmptyMetricsClient, nil),
		// sdk.WithRegisterCacheClient(testdata.MockCacheClient(t)),
		// sdk.WithRegisterDMPClient(testdata.MockEmptyDMPClient),
		sdk.WithEnvType(env.TypePrd),
		// sdk.WithDisableReport(true),
		sdk.WithSecretKey("secretKey"))
	assert.Nil(t, err)
	sg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		sg.Add(1)
		go func() {
			defer sg.Done()
			start := time.Now()
			var hitResult = make(map[int64]int64)
			for j := 0; j < times; j++ {
				result, err := sdk.NewUserContext(uuid.New().String(),
					sdk.WithTagKV("genderx", "女2@"),
					sdk.WithTagKV("appversionx", "3.0.21"),
				).GetExperiments(context.TODO(),
					projectID, sdk.WithIsDisableDMP(true))
				if err != nil {
					t.Errorf("GetExperiment fail:%v", err)
					return
				}
				sdk.LogExperimentsExposure(context.TODO(), projectID, result)
				assert.NotNil(t, result)
				for _, group := range result.Data {
					hitResult[group.ID]++
				}
			}
			t.Logf("%+v", hitResult)
			t.Logf("%v op/s", (time.Since(start) / time.Duration(times)).String())
		}()
	}
	sg.Wait()
}

func TestGetValueByVariantKey(t *testing.T) {
	t.Skip()
	defer sdk.Release()
	err := sdk.Init(context.TODO(), []string{"123"},
		sdk.WithRegisterCacheClient(testdata.MockCacheClient(t)),
		sdk.WithSecretKey("secretKey"))
	assert.Nil(t, err)

	userCtx := sdk.NewUserContext("123213AABB",
		sdk.WithTags(map[string][]string{
			"country": []string{"cn"},
		}), sdk.WithTagKV("age", "18"))
	t.Run("getValueByVariantKey", func(t *testing.T) {
		vr, err := userCtx.GetValueByVariantKey(context.TODO(), "123", "variant_key_test")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		flag := vr.GetBoolWithDefault(false)
		t.Logf("%v", flag)
	})
}

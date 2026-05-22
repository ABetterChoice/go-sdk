// Package example 演示如何在正式环境下一次性获取项目内所有远程配置。
//
// 默认通过 t.Skip 跳过，需要时设置以下环境变量后手动运行：
//
//	ABC_PROJECT_ID  目标项目 ID（必填）
//	ABC_SECRET_KEY  对应站点签发的 SecretKey（必填）
//	ABC_PRD_ADDR    可选，正式环境接入点。
//	                默认走 SDK 内置的 https://openapi.abetterchoice.ai；
//	                若项目部署在 SG 站点，请设为 https://openapi.sg.abetterchoice.ai
//
// 运行示例：
//
//	ABC_PROJECT_ID=60077 \
//	ABC_SECRET_KEY=eyJ... \
//	ABC_PRD_ADDR=https://openapi.sg.abetterchoice.ai \
//	go test -run TestListAllRemoteConfigs -v ./example
package example

import (
	"context"
	"os"
	"sort"
	"testing"

	sdk "github.com/abetterchoice/go-sdk"
	"github.com/abetterchoice/go-sdk/env"
)

func TestListAllRemoteConfigs(t *testing.T) {
	projectID := os.Getenv("ABC_PROJECT_ID")
	secretKey := os.Getenv("ABC_SECRET_KEY")
	if projectID == "" || secretKey == "" {
		t.Skip("set ABC_PROJECT_ID and ABC_SECRET_KEY to run this example")
	}

	if addr := os.Getenv("ABC_PRD_ADDR"); addr != "" {
		if err := env.RegisterAddr(env.TypePrd, addr); err != nil {
			t.Fatalf("RegisterAddr(%s) fail: %v", addr, err)
		}
	}

	defer sdk.Release()
	if err := sdk.Init(context.TODO(), []string{projectID},
		sdk.WithEnvType(env.TypePrd),
		sdk.WithSecretKey(secretKey),
		sdk.WithDisableReport(true),
	); err != nil {
		t.Fatalf("Init fail: %v", err)
	}

	configs, err := sdk.GetAllRemoteConfigs(projectID)
	if err != nil {
		t.Fatalf("GetAllRemoteConfigs fail: %v", err)
	}

	t.Logf("[projectID=%s] total remote configs: %d", projectID, len(configs))
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		t.Logf("  %s = %q", k, configs[k].String())
	}
}

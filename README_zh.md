# Go SDK

ABetterChoice Go SDK 用于服务端实验分流、Feature Flag 评估、远程配置读取和曝光上报。

## 支持平台

- Linux
- macOS
- Windows

## 前置条件

- Go 1.17+
- 已在 ABetterChoice 控制台创建项目
- `ProjectID` 和 `SecretKey`

参考文档：
- [创建项目](https://docs.abetterchoice.ai/guide/getting-started/create-project)
- [创建实验](https://docs.abetterchoice.ai/guide/features/create-experiment)
- [Feature Flags](https://docs.abetterchoice.ai/guide/features/feature-flags)

## 安装

```bash
go get github.com/abetterchoice/go-sdk
```

## 快速开始（推荐）

```go
package main

import (
    "context"
    "log"

    abc "github.com/abetterchoice/go-sdk"
)

func main() {
    ctx := context.Background()
    projectID := "YOUR_PROJECT_ID"
    secretKey := "YOUR_SECRET_KEY"
    unitID := "YOUR_UNIT_ID"

    // 每个进程初始化一次
    err := abc.Init(ctx, []string{projectID}, abc.WithSecretKey(secretKey))
    if err != nil {
        log.Fatalf("abc.Init failed: %v", err)
    }
    defer abc.Release()

    // 构建用户上下文，后续评估复用
    userCtx := abc.NewUserContext(unitID, abc.WithTagKV("region", "sea"))

    // 1) 推荐：通过全局唯一参数 key 读取
    value, err := userCtx.GetValueByVariantKey(ctx, projectID, "drop_rate")
    if err != nil {
        log.Printf("GetValueByVariantKey failed: %v", err)
    } else {
        dropRate := value.GetFloat64WithDefault(0.05)
        log.Printf("drop_rate=%v", dropRate)
    }

    // 2) 按层读取实验
    exp, err := userCtx.GetExperiment(ctx, projectID, "matchmaking_layer")
    if err != nil {
        log.Printf("GetExperiment failed: %v", err)
    } else if exp != nil {
        mmrWindow := exp.GetInt64WithDefault("mmr_window", 100)
        log.Printf("mmr_window=%d", mmrWindow)
    }

    // 3) Feature Flag
    ff, err := userCtx.GetFeatureFlag(ctx, projectID, "new_shop")
    if err != nil {
        log.Printf("GetFeatureFlag failed: %v", err)
    } else if ff != nil && ff.MustGetBool() {
        log.Printf("new_shop enabled")
    }
}
```

## API 选型建议

| 场景 | API | 说明 |
| --- | --- | --- |
| 通过参数 key 读取远程参数 | `GetValueByVariantKey` | 业务代码不依赖层名，支持实验优先 + 配置回退。 |
| 需要拿某个已知层的完整分流结果 | `GetExperiment` | 返回组信息及层内参数。 |
| 需要项目下全部命中层 | `GetExperiments` | 适用于批量透传或诊断。 |
| 需要项目下所有远程配置快照 | `GetAllRemoteConfigs` | 一次读取本地缓存中的全部远程配置键值。 |
| 简单开关控制 | `GetFeatureFlag` | 布尔开关语义清晰，支持类型化读取。 |

## 初始化与生命周期

### `Init(ctx, projectIDList, opts...)`

`Init` 会执行远程拉取与本地缓存初始化，建议在进程启动时调用一次。

常用选项：

- `WithSecretKey(secretKey)`（必填）
- `WithEnvType(envType)`（可选，默认正式环境）
- `WithDisableReport(true|false)`（可选，全局曝光开关）
- `WithRegionCode(regionCode)`（可选，按区域拉取配置）
- `WithRegisterCacheClient(...)` / `WithRegisterDMPClient(...)` / `WithRegisterMetricsPlugin(...)`（高级扩展）

示例：

```go
err := abc.Init(
    context.Background(),
    []string{"YOUR_PROJECT_ID"},
    abc.WithSecretKey("YOUR_SECRET_KEY"),
    abc.WithDisableReport(false),
)
```

### `Release()`

进程退出前调用 `Release()`，清理内存状态。

## 用户上下文与属性

通过 `NewUserContext(unitID, opts...)` 构建用户上下文。

常见属性选项：

- `WithTags(map[string][]string)`
- `WithTagKV(key, value)`
- `WithDecisionID(decisionID)`
- `WithNewUnitID(newUnitID)` / `WithNewDecisionID(newDecisionID)`（迁移场景）
- `WithExpandedData(map[string]string)`（补充曝光数据）

示例：

```go
userCtx := abc.NewUserContext(
    "player_1001",
    abc.WithTagKV("channel", "appstore"),
    abc.WithTagKV("server_id", "s1"),
)
```

## 评估 API

### 获取 Feature Flag

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
featureFlag, err := userCtx.GetFeatureFlag(context.TODO(), "project_id", "new_feature_flag")
if err == nil && featureFlag != nil {
    flagValue := featureFlag.MustGetBool()
    _ = flagValue
}
```

### 按层获取实验

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
experiment, err := userCtx.GetExperiment(context.TODO(), "project_id", "abc_layer_name")
if err == nil && experiment != nil {
    shouldShowBanner := experiment.GetBoolWithDefault("should_show_banner", false)
    _ = shouldShowBanner
}
```

### 按参数 key 获取值（`GetValueByVariantKey`）

`GetValueByVariantKey` 的解析顺序：

1. 若参数 key 归属于一个或多个实验层，按“每层最早实验 ID”升序遍历；
2. 返回第一个命中非默认组的层；
3. 若没有命中非默认组，回退到配置读取路径（`GetRemoteConfig` 逻辑）。

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
result, err := userCtx.GetValueByVariantKey(context.TODO(), "project_id", "should_show_banner")
if err == nil && result != nil {
    shouldShowBanner := result.GetBoolWithDefault(false)
    _ = shouldShowBanner

    // 便于上报与排障：
    // LayerKey / ExperimentKey / ExperimentID / GroupKey / VariantID / ConfigKey
    _ = result.Detail
}
```

### 获取项目下全部实验命中

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
experiments, err := userCtx.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
if err == nil && experiments != nil {
    _ = experiments
}
```

### 获取项目下全部远程配置

```go
configs, err := abc.GetAllRemoteConfigs("project_id")
if err == nil {
    for key, value := range configs {
        _ = key
        _ = value.String()
    }
}
```

## 曝光策略

默认情况下，实验/配置/Feature Flag 查询会自动上报曝光。

如果你是提前拉取，只有真正展示时才希望上报，建议改为手动曝光。

### API 级关闭自动曝光

```go
experiments, err := userCtx.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
if err == nil && experiments != nil {
    _ = abc.LogExperimentsExposure(context.TODO(), "project_id", experiments)
}
```

### 手动曝光 API

- `LogExperimentExposure(ctx, projectID, experimentResult)`
- `LogExperimentsExposure(ctx, projectID, experimentList)`
- `LogFeatureFlagExposure(ctx, projectID, featureFlag)`
- `LogRemoteConfigExposure(ctx, projectID, configResult)`

### 全局关闭曝光

```go
err := abc.Init(
    context.Background(),
    []string{"YOUR_PROJECT_ID"},
    abc.WithSecretKey("YOUR_SECRET_KEY"),
    abc.WithDisableReport(true),
)
```

## 高级选项

### 实验选项

可用于 `GetExperiment(s)` 和 `GetValueByVariantKey`：

- `WithLayerKey(layerKey)` / `WithLayerKeyList(layerKeys)`
- `WithSceneID(sceneID)` / `WithSceneIDList(sceneIDs)`
- `WithAutomatic(isAutomatic)`
- `WithIsPreparedDMPTag(isPreparedDMPTag)`
- `WithIsDisableDMP(isDisableDMP)`

### 配置选项

配置 API 与实验选项共享同一类型，同时提供：

- `WithIsPreparedDMPTagConfigOpt(...)`
- `WithIsDisableDMPConfigOpt(...)`

### 多项目注册

初始化后可继续注册其他项目：

```go
err := abc.RegisterProjectIDs(context.Background(), []string{"PROJECT_B", "PROJECT_C"})
```

## 排查指南

### 1) `Init` 失败

- 检查 `SecretKey` 与 `ProjectID` 是否匹配。
- 检查网络连通性。
- 如果使用自定义环境地址，确认 `env.RegisterAddr(...)` 配置正确。

### 2) 总是拿到默认值

- 检查 `unitID` 是否为空或不稳定。
- 检查 key 是否在控制台存在（`layerKey` / `variantKey` / `featureKey`）。
- 检查用户 tags 是否满足受众条件。

### 3) 没有曝光数据

- 确认没有误开 `WithDisableReport(true)`。
- 如果使用手动模式（`WithAutomatic(false)`），确认调用了手动曝光 API。
- 确认后端指标插件与采样配置可用。

### 4) 类型转换失败

- 优先使用带默认值的 getter（如 `GetBoolWithDefault`、`GetInt64WithDefault`）。
- 检查控制台参数类型是否与代码预期一致。

## FAQ

### `GetExperiment` 和 `GetValueByVariantKey` 怎么选？

优先拿实验分流语义时，用 `GetExperiment`；优先按参数 key 读取值时，用 `GetValueByVariantKey`。

### `GetRemoteConfig` 还支持吗？

支持。它主要用于兼容历史逻辑；新接入建议按场景优先使用 `GetFeatureFlag` 或 `GetValueByVariantKey`。

### 测试环境能否关闭曝光？

可以，初始化时设置 `WithDisableReport(true)`。

## API 参考（核心导出）

- 初始化：`Init`, `Release`, `RegisterProjectIDs`, `GetGlobalConfig`
- 用户上下文：`NewUserContext`, `WithTags`, `WithTagKV`, `WithDecisionID`, `WithNewUnitID`, `WithNewDecisionID`, `WithExpandedData`
- 评估：`GetExperiment`, `GetExperiments`, `GetFeatureFlag`, `GetValueByVariantKey`, `GetAllRemoteConfigs`, `GetRemoteConfig`
- 手动曝光：`LogExperimentExposure`, `LogExperimentsExposure`, `LogFeatureFlagExposure`, `LogRemoteConfigExposure`

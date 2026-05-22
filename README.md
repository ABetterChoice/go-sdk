# Go SDK

ABetterChoice Go SDK for server-side experiment assignment, feature flag evaluation, remote config lookup,
and exposure reporting.

## Supported platforms

- Linux
- macOS
- Windows

## Prerequisites

- Go 1.17+
- A project created in ABetterChoice console
- `ProjectID` and `SecretKey`

Create project and experiment references:
- [Create Project](https://docs.abetterchoice.ai/guide/getting-started/create-project)
- [Create Experiment](https://docs.abetterchoice.ai/guide/features/create-experiment)
- [Feature Flags](https://docs.abetterchoice.ai/guide/features/feature-flags)

## Installation

```bash
go get github.com/abetterchoice/go-sdk
```

## Quick Start (recommended)

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

    // Initialize once per process.
    err := abc.Init(ctx, []string{projectID}, abc.WithSecretKey(secretKey))
    if err != nil {
        log.Fatalf("abc.Init failed: %v", err)
    }
    defer abc.Release()

    // Build user context for all evaluations.
    userCtx := abc.NewUserContext(unitID, abc.WithTagKV("region", "sea"))

    // 1) Recommended API: fetch by globally unique variant key.
    value, err := userCtx.GetValueByVariantKey(ctx, projectID, "drop_rate")
    if err != nil {
        log.Printf("GetValueByVariantKey failed: %v", err)
    } else {
        dropRate := value.GetFloat64WithDefault(0.05)
        log.Printf("drop_rate=%v", dropRate)
    }

    // 2) Layer-based experiment lookup.
    exp, err := userCtx.GetExperiment(ctx, projectID, "matchmaking_layer")
    if err != nil {
        log.Printf("GetExperiment failed: %v", err)
    } else if exp != nil {
        mmrWindow := exp.GetInt64WithDefault("mmr_window", 100)
        log.Printf("mmr_window=%d", mmrWindow)
    }

    // 3) Feature flag lookup.
    ff, err := userCtx.GetFeatureFlag(ctx, projectID, "new_shop")
    if err != nil {
        log.Printf("GetFeatureFlag failed: %v", err)
    } else if ff != nil && ff.MustGetBool() {
        log.Printf("new_shop enabled")
    }
}
```

## API selection guide

| Scenario | API | Why |
| --- | --- | --- |
| Remote parameter lookup by parameter key | `GetValueByVariantKey` | Decouples business code from layer names; supports experiment-first with fallback behavior. |
| Need full assignment for a known layer | `GetExperiment` | Returns group details and layer-level parameter access. |
| Need all assigned layers in a project | `GetExperiments` | Batch retrieval for downstream forwarding or diagnostic scenarios. |
| Need all remote config snapshots (non-experiment path) | `GetAllRemoteConfigs` | Reads all remote-config keys from local cache in one call. |
| Simple on/off feature control | `GetFeatureFlag` | Boolean-first feature gating with typed getters. |

## Initialization and lifecycle

### `Init(ctx, projectIDList, opts...)`

`Init` performs network fetches and local cache initialization. Call once at process startup.

Key options:

- `WithSecretKey(secretKey)` (required)
- `WithEnvType(envType)` (optional, default production)
- `WithDisableReport(true|false)` (optional, global exposure reporting switch)
- `WithRegionCode(regionCode)` (optional, region-aware config delivery)
- `WithRegisterCacheClient(...)` / `WithRegisterDMPClient(...)` / `WithRegisterMetricsPlugin(...)` for advanced integration

Example:

```go
err := abc.Init(
    context.Background(),
    []string{"YOUR_PROJECT_ID"},
    abc.WithSecretKey("YOUR_SECRET_KEY"),
    abc.WithDisableReport(false),
)
```

### `Release()`

Call `Release()` on graceful shutdown to reset in-memory state.

## User context and attribution

Build user context with `NewUserContext(unitID, opts...)`.

Common attribution options:

- `WithTags(map[string][]string)`
- `WithTagKV(key, value)`
- `WithDecisionID(decisionID)`
- `WithNewUnitID(newUnitID)` / `WithNewDecisionID(newDecisionID)` for migration scenarios
- `WithExpandedData(map[string]string)` to enrich exposure logs

Example:

```go
userCtx := abc.NewUserContext(
    "player_1001",
    abc.WithTagKV("channel", "appstore"),
    abc.WithTagKV("server_id", "s1"),
)
```

## Evaluation APIs

### Get feature flag

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
featureFlag, err := userCtx.GetFeatureFlag(context.TODO(), "project_id", "new_feature_flag")
if err == nil && featureFlag != nil {
    flagValue := featureFlag.MustGetBool()
    _ = flagValue
}
```

### Get experiment by layer key

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
experiment, err := userCtx.GetExperiment(context.TODO(), "project_id", "abc_layer_name")
if err == nil && experiment != nil {
    shouldShowBanner := experiment.GetBoolWithDefault("should_show_banner", false)
    _ = shouldShowBanner
}
```

### Get value by variant key

`GetValueByVariantKey` resolves a globally unique parameter key with this order:

1. If the key belongs to one or more experiment layers, traverse layers by each layer's earliest experiment ID (ascending).
2. Return the first audience-matched non-default group.
3. If no layer matches a non-default group, fallback to config lookup (`GetRemoteConfig` path).

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
result, err := userCtx.GetValueByVariantKey(context.TODO(), "project_id", "should_show_banner")
if err == nil && result != nil {
    shouldShowBanner := result.GetBoolWithDefault(false)
    _ = shouldShowBanner

    // Detail fields help reporting and debugging:
    // LayerKey / ExperimentKey / ExperimentID / GroupKey / VariantID / ConfigKey
    _ = result.Detail
}
```

### Get all experiments in a project

```go
userCtx := abc.NewUserContext("{{.UnitID}}")
experiments, err := userCtx.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
if err == nil && experiments != nil {
    _ = experiments
}
```

### Get all remote configs in a project

```go
configs, err := abc.GetAllRemoteConfigs("project_id")
if err == nil {
    for key, value := range configs {
        _ = key
        _ = value.String()
    }
}
```

## Exposure strategy

By default, experiment/config/flag retrieval logs exposure automatically.

Use manual exposure if you prefetch values but only want to report when users actually see the feature.

### Disable auto exposure at API level

```go
experiments, err := userCtx.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
if err == nil && experiments != nil {
    _ = abc.LogExperimentsExposure(context.TODO(), "project_id", experiments)
}
```

### Manual exposure APIs

- `LogExperimentExposure(ctx, projectID, experimentResult)`
- `LogExperimentsExposure(ctx, projectID, experimentList)`
- `LogFeatureFlagExposure(ctx, projectID, featureFlag)`
- `LogRemoteConfigExposure(ctx, projectID, configResult)`

### Disable exposure globally

```go
err := abc.Init(
    context.Background(),
    []string{"YOUR_PROJECT_ID"},
    abc.WithSecretKey("YOUR_SECRET_KEY"),
    abc.WithDisableReport(true),
)
```

## Advanced options

### Experiment options

Use on `GetExperiment(s)` and `GetValueByVariantKey`:

- `WithLayerKey(layerKey)` / `WithLayerKeyList(layerKeys)`
- `WithSceneID(sceneID)` / `WithSceneIDList(sceneIDs)`
- `WithAutomatic(isAutomatic)`
- `WithIsPreparedDMPTag(isPreparedDMPTag)`
- `WithIsDisableDMP(isDisableDMP)`

### Config options

Config APIs share the same option type:

- `WithIsPreparedDMPTagConfigOpt(...)`
- `WithIsDisableDMPConfigOpt(...)`

### Multi-project registration

Register additional projects after init:

```go
err := abc.RegisterProjectIDs(context.Background(), []string{"PROJECT_B", "PROJECT_C"})
```

## Troubleshooting

### 1) `Init` failed

- Verify `SecretKey` and `ProjectID`.
- Verify network access to backend services.
- If using custom environment, confirm `env.RegisterAddr(...)` is configured correctly.

### 2) Always getting default values

- Check whether `unitID` is empty or unstable.
- Verify the key exists in console (`layerKey`, `variantKey`, or `featureKey`).
- Verify targeting conditions match user tags.

### 3) Exposure data missing

- Confirm `WithDisableReport(true)` is not enabled unintentionally.
- If using manual mode (`WithAutomatic(false)`), ensure manual log API is called.
- Confirm metrics plugin and sampling configuration are enabled in backend settings.

### 4) Value type conversion failure

- Use typed getters with defaults (`GetBoolWithDefault`, `GetInt64WithDefault`, etc.).
- Verify value type in console matches code expectation.

## FAQ

### Should I use `GetExperiment` or `GetValueByVariantKey`?

Use `GetExperiment` when your primary goal is to retrieve experiment assignment (recommended in code comments). Use `GetValueByVariantKey` when you want to read a parameter directly by globally unique key.

### Is `GetRemoteConfig` still supported?

`GetRemoteConfig` is kept for compatibility and remains supported. For new access paths, prefer `GetFeatureFlag` or `GetValueByVariantKey` based on your scenario.

### Can I turn off exposure globally in test environments?

Yes, use `WithDisableReport(true)` at `Init`.

## API reference (exported core APIs)

- Initialization: `Init`, `Release`, `RegisterProjectIDs`, `GetGlobalConfig`
- User context: `NewUserContext`, `WithTags`, `WithTagKV`, `WithDecisionID`, `WithNewUnitID`, `WithNewDecisionID`, `WithExpandedData`
- Evaluation: `GetExperiment`, `GetExperiments`, `GetFeatureFlag`, `GetValueByVariantKey`, `GetAllRemoteConfigs`, `GetRemoteConfig`
- Manual exposure: `LogExperimentExposure`, `LogExperimentsExposure`, `LogFeatureFlagExposure`, `LogRemoteConfigExposure`


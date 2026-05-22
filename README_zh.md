# Go SDK

## 支持的平台

Golang SDK 支持以下平台：

- Linux
- MacOS
- Windows

## 快速开始

如果你还没有项目，请先按照[文档](https://docs.abetterchoice.ai/guide/getting-started/create-project)创建一个 ABetterChoice 项目和若干实验。这一步可以暂时跳过，但后续章节会用到项目 ID 和实验层名。

## 安装

SDK 通过 `go get` 安装，需要 Go 1.17 或以上版本。

```
go get git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk
```

## SDK 初始化

安装完成后需要先初始化。`Init` 接受三个参数：

1. 第一个参数按 Go 惯例传 `context`。
2. 第二个参数传一个或多个项目 ID。
3. 第三个参数是 `InitOption`，用于灵活配置初始化项。目前必传的是 SecretKey，可以在控制台 setting 页的 API key 区域取到。

```go
// Package main ...
package main

import (
 "context"
 abc "git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk"
 "log"
)

abc.Init(context.TODO(), []string{"{{.ProjectID}}"}, abc.WithSecretKey("{{.SecretKey}}"))
```

*进阶：首次接入可跳过，按需再回看*
`InitOption` 还提供其他初始化项，详细列表见 `InitOption` 实现。除上面的 SecretKey 之外，比较常用的是关闭曝光上报。默认情况下，调用 `GetExperiment` 时 SDK 会自动上报一次曝光（即记录某个实验单元命中实验）；某些场景下这会带来过多预期外的曝光，稀释实验结果。
你可以在调用实验接口时不上报曝光，等到合适的时机再补上报。设置 `WithDisableReport(true)` 可以全局关闭自动曝光，也可以在单个接口调用上局部关闭，后文会展开。

```go
abc.Init(context.TODO(), []string{"{{.ProjectID}}"}, abc.WithSecretKey("{{.SecretKey}}"), abc.WithDisableReport(true))
```

*进阶：多区域接入*
ABC 部署在多个区域，SDK 默认指向全局（NA）端点。如果你的项目部署在新加坡（SG）区域，请在 `Init` **之前** 注册 SG 端点，并使用该区域签发的 SecretKey。

```go
import "git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk/env"

_ = env.RegisterAddr(env.TypePrd, "https://openapi.sg.abetterchoice.ai")
abc.Init(context.TODO(), []string{"{{.ProjectID}}"},
    abc.WithEnvType(env.TypePrd),
    abc.WithSecretKey("{{.SecretKey}}"))
```

## 获取 Feature Flag

本节介绍如何获取 feature flag 的取值。如果尚未创建 feature flag，请先按[文档](https://docs.abetterchoice.ai/guide/features/feature-flags)创建一个。
假设我们已在 `project_id` 项目下创建了名为 `new_feature_flag`、类型为 Boolean 的 feature flag，可以这样取值：

```go
abcUserContext, err := abc.NewUserContext("{{.UnitID}}")
featureFlag, err = abcUserContext.GetFeatureFlag(context.TODO(), "project_id", "new_feature_flag")
flagValue = featureFlag.MustGetBool()
```

## 获取所有远程配置

如果只想拿到某个项目下所有远程配置的快照——可以调用 `GetAllRemoteConfigs`。返回本地缓存中所有**非实验**远程配置 key 及其默认取值。

```go
configs, err := abc.GetAllRemoteConfigs("project_id")
if err == nil {
    for key, value := range configs {
        log.Infof("%s = %s", key, value.String())
    }
}
```

## 获取实验和上报曝光

本节介绍如何创建实验，并通过 SDK 提供的接口获取实验分流结果。


| 术语             | 含义                                                               |
| -------------- | ---------------------------------------------------------------- |
| UnitID         | 实验单元 ID，可以是用户 ID、Session ID 或机器 ID。SDK 会保证同一个 UnitID 始终落到同一个实验组。 |
| Layer（层）       | 一个层通常代表你产品中可用于实验的 100% 流量，可被进一步切分给一个或多个实验，同层内实验流量互斥。             |
| Experiment（实验） | 单个实验最多占用所在层 100% 的流量，会进一步包含两个或多个实验组。                             |
| Group（组）       | 同一个实验中，单元会被分到两个或多个组（如对照组、实验组）做对比。                                |
| Parameter（参数）  | 绑定在某一层上的取值，同层内不同实验可以共享同一组参数。                                     |


### 1. 创建实验

如果还没创建实验，请按照[创建实验文档](https://abetterchoice-test.woa.com/docs/create_an_experiment)在控制台创建一个。

### 2. 获取实验分流

假设第 1 步在项目 `project_id` 下创建了一个层名为 `abc_layer_name` 的实验，层内有一个名为 `should_show_banner` 的 Boolean 类型参数。注意基础版里层名默认就是实验名；如果在同一层下又创建了新实验，需要从实验页面读取层名。
当前提供三种方式拿到指定 unit id 的分流结果，所有接口都需要先构造一个携带 unit id 的 `UserContext`。

```go
abcUserContext, err := abc.NewUserContext("{{.UnitID}}")
```

#### a. 通过层 Key 获取实验

第一种方式是按层 Key 获取实验分流，并用强类型 API 取层内参数。当层/实验内有多个参数时这种方式比较方便。同一层下有多个实验时，单元只会落到其中一个，因此该接口能自动返回正确结果。这种方式还有一个好处：废弃旧实验、在同一层下新建实验时业务代码不用改。

```go
experiment, err := abcUserContext.GetExperiment(context.TODO(), "project_id", "abc_layer_name")
if err == nil {
    shouldShowBanner := experiment.GetBoolWithDefault("should_show_banner", false)
}
```

#### b. 通过参数 Key 获取实验

这种方式和上面类似，区别在于不按层名查，而是按参数名查（参数名在同一项目下是唯一的）。后续如果支持把参数从一层迁到另一层（比如发车层），这个接口会很有用。

当同一个参数 Key 被多个层暴露时，SDK 会按"每个层下最早创建的实验 ID（即该层最小的实验 ID）"给这些层排序，再依次判断：先返回第一个**人群命中且落在非默认组**的层，否则回落到第一个命中的默认组；多层命中不再报错。

```go
result, err := abcUserContext.GetValueByVariantKey(context.TODO(), "project_id", "should_show_banner")
if err == nil {
    shouldShowBanner := result.GetBoolWithDefault(false)
    // result.Detail 同时带有 LayerKey / ExperimentKey / ExperimentID / GroupKey / VariantID，可用于上报。
}
```

#### c. 获取项目下全部实验分流

某些场景下需要一次性拿到一个项目下所有的实验分流结果，SDK 提供了批量接口，返回值是 “层名 -> 该层内 unit 命中的实验分流” 的 map。前面提到这种方式容易导致曝光稀释，使用时请谨慎。可以按需通过 opts 关掉自动曝光（如下例），后续再用曝光接口手动补报。

```go
experiments, err := abcUserContext.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
```

### 3. 上报曝光

按上面方式获取分流时可以选择关闭自动曝光，避免稀释；之后在合适时机手动上报：

```go
abc.LogExperimentExposure(context.TODO(), "{{.ProjectID}}", experiment)
```

### 4. 完整代码示例

完整代码示例

```go
// Package main TODO
package main

import (
  "context"

  abc "git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk"
  "git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk/env"
  "git.woa.com/tencent_abtest/open-source/abetterchoice-go-sdk/plugin/log"
)

func main() {
  defer abc.Release()
  projectID := "6666"
  // 进程内只需初始化一次
  err := abc.Init(context.TODO(), []string{projectID},
  	abc.WithEnvType(env.TypePrd),
 	abc.WithSecretKey("{{.SecretKey}}"))
  if err != nil {
    log.Errorf("Init fail:%v", err)
    return
  }

  // 拉取项目下所有远程配置（不评估人群、不上报曝光）
  configs, _ := abc.GetAllRemoteConfigs(projectID)
  for k, v := range configs {
    log.Infof("remote_config %s = %s", k, v.String())
  }

  abcUserContext := abc.NewUserContext("unitID_demo")

  // a. 通过层 Key 拿实验，关掉自动曝光，后面手动上报
  experiment, err := abcUserContext.GetExperiment(context.TODO(), projectID, "abc_layer_name", abc.WithAutomatic(false))
  if err != nil {
    log.Errorf("GetExperiment fail:%v", err)
    return
  }
  log.Infof("should_show_banner=%v", experiment.GetBoolWithDefault("should_show_banner", false))
  _ = abc.LogExperimentExposure(context.TODO(), projectID, experiment)

  // b. 通过参数 Key 拿实验
  if result, err := abcUserContext.GetValueByVariantKey(context.TODO(), projectID, "should_show_banner"); err == nil {
    log.Infof("variant=%v layer=%s expID=%d",
      result.GetBoolWithDefault(false), result.Detail.LayerKey, result.Detail.ExperimentID)
  }

  // c. 批量拿项目下所有实验
  if experiments, err := abcUserContext.GetExperiments(context.TODO(), projectID, abc.WithAutomatic(false)); err == nil {
    for layerKey := range experiments.Data {
      log.Infof("hit layer=%s", layerKey)
    }
  }

  // Feature flag
  if featureFlag, err := abcUserContext.GetFeatureFlag(context.TODO(), projectID, "new_feature_flag"); err == nil {
    log.Infof("the feature flag value is %v", featureFlag.MustGetBool())
  }
}
```


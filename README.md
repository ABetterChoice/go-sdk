# Go SDK

## Supported platforms

The Golang SDK is compatible with the following platforms:
- Linux
- MacOS
- Windows

## Getting Started

If you haven't already, please follow the [doc](https://abetterchoice.ai/docs?wj-docs=%2Fprod%2Fdocs%2Fcreate_a_project) to create a ABetterChoice project and some experiments. While you can choose to bypass this step for now, please note that you may require the project ID and the names of the experiment layers in subsequent sections.

## Installation
In order to use Golang SDK, you need to install the SDK by `go get` CLI, which requires Go 1.17 or higher.

```
go get github.com/abetterchoice/go-sdk
```

## SDK Initialization

After installing the SDK, it's necessary to initialize it. The Init function requires three parameters:
1. The first parameter, as per Golang convention, is context.
2. The second parameter allows for the input of one or more project IDs.
3. The third parameter is of the type InitOption, offering flexible initialization settings. Currently, only the secret key is required. You can locate this under the API key section within the settings tab.
```go
// Package main ...
package main

import (
 "context"
 abc "github.com/abetterchoice/go-sdk"
 "log"
)

abc.Init(context.TODO(), []string{"{{.ProjectID}}"}, abc.WithSecretKey("{{.SecretKey}}")) 
```
*\[Advanced: This section can be skipped initially and revisited later when needed\]*
InitOption provides additional initialization configurations. For more details, please refer to the InitOption implementation. Apart from the secret key used in the preceding code, another commonly used option is to disable exposure logging. By default, the ABC SDK logs exposures (essentially, records of a specific experiment unit being exposed to your experiment) automatically when the `getExperiment` API is called. However, in certain scenarios, this might result in excessive unexpected exposure logs, potentially diluting your results.
To prevent this, you can fetch the experiment result without logging exposures when calling the `getExperiment` API, and log the exposure later at a more appropriate point. By setting `WithDisableReport()`, exposure logging will be globally disabled. We also provide a method to disable exposure logging at the individual API level, which will be explained later in this document.
```go
abc.Init(context.TODO(), []string{"{{.ProjectID}}"}, abc.WithSecretKey("{{.SecretKey}}").WithDisableReport()) 
```

## Checking feature flags

In this section, we will guide you through the process of retrieving the value of a feature flag. If you haven't already done so, please first follow our [documentation](https://abetterchoice.ai/docs?wj-docs=%2Fprod%2Fdocs%2Ffeature_flags) to create one.
Assuming we have already created a new feature flag named `new_feature_flag` under the project `project_id`, and its value type is Boolean, we can fetch the value in the following manner:
```go
abcUserContext, err := abc.NewUserContext("{{.UnitID}}")
featureFlag, err = abcUserContext.GetFeatureFlag(context.TODO(), "project_id", "new_feature_flag")
flagValue = featureFlag.MustGetBool()
```


## Getting experiments and logging Exposures

In this section, we will guide you to create a new experiment and fetch the experiment assignments via our experiment fetching APIs provided by our SDK.

| <div style="width:160px">Terms</div>        | Meaning           | 
| ------------- |:-------------:| 
| UnitID      | A unit ID can be a user ID, a session ID or a machine ID. The ABC SDK assigns a unit ID to the same group consistantly. | 
| Layer      | A layer usually represent 100% of the units that you can run experiment in your product, it can be further splited into one or more experiments, the traffic of the experiments under same layer are mutually exclusive.  |  
| Experiment | An experiment can take up to 100% of the layer traffic, it will further contain two or more experiment groups.     | 
| Group | Units within one experiment will be split into two or more groups, e.g. control and treatment, they will be compared against each other.      | 
| Parameter | These are the values binded to a particular layer. Experiments under the same layer can share the same parameters with each other.  | 


### 1. Creating a New Experiment

If you haven't done so already, please navigate to the console and create a new experiment by following our documentation on [how to create an experiment](https://abetterchoice-test.woa.com/docs/create_an_experiment).
### 2. Retrieving Experiment Assignments

Assume that in step 1, under the project `project_id`, we have already created an experiment with the layer name `abc_layer_name`. Within this layer, there is a parameter named `should_show_banner` of Boolean type. Please note that in our basic version, the layer name defaults to the experiment name. However, if you create new experiments under the same layer, you will need to find the layer name on the experiment page.
Currently, we offer three methods to fetch the traffic assignment for a specific unit id. All these APIs require you to establish a user context that encapsulates the unit ID.
```go
abcUserContext, err := abc.NewUserContext("{{.UnitID}}")
```
#### a. Retrieve Experiment by Layer Key

The first method involves obtaining the experiment assignment by the layer key and fetching the parameters within the layer using the strong type APIs. This is useful when you have multiple parameters within the layer/experiment. This API will automatically provide the correct result when there are multiple experiments under the same layer, as the unit will and can only fall into one of them. An additional advantage of this method is that it allows you to iterate through the experiments quickly by deallocating the old experiment and creating new ones under the same layer without modifying the code.
```go
experiment, err := abcUserContext.GetExperiment(context.TODO(), "project_id", "abc_layer_name")
if err == nil {
	shouldShowBanner := experiment.GetBoolWithDefault("should_show_banner", false)
}
```
#### b. Retrieve Experiment by Parameter Key

This method is similar to the previous one, with the only difference being that instead of fetching by layer name, we retrieve the result by parameter name, which is unique within the same project. We are planning to introduce features that may allow users to move the parameter from one layer to another (e.g., launcher layer), and this API will be useful in handling such cases.
```go
configParameter, err := abcUserContext.GetValueByVariantKey(context.TODO(), "project_id", "should_show_banner")
if err == nil {
	shouldShowBanner := configParameter.GetBoolWithDefault("should_show_banner", false)
}
```
#### c. Retrieve All Experiment Assignments Under the Project

In some cases, users may need to fetch all experiment assignment results under the same project. To accommodate this, we also provide a batch API, which will return a map linking the experiment layer name to the experiment assignment of the unit specific to that layer. As mentioned earlier, this could lead to dilution problems, so caution is advised when using this API. Depending on your use case, you might consider disabling exposure logging via the opts parameter as shown below. You can use the exposure logging API to manually log the exposure later.
```go
experiments, err := abcUserContext.GetExperiments(context.TODO(), "project_id", abc.WithAutomatic(false))
```
### 3. Logging Exposure Data

When retrieving the experiment assignments as described above, you have the option to disable exposure logging to avoid exposure dilution. You can manually log the exposure later at an appropriate time:
```go
abc.LogExperimentExposure(context.TODO(), "{{.ProjectID}}", experiment) 
```
### 4. Complete code example
Complete code example
```go
// Package main TODO
package main

import (
  "context"

  abc "github.com/abetterchoice/go-sdk"
  "github.com/abetterchoice/go-sdk/env"
  "github.com/abetterchoice/go-sdk/plugin/log"
)

func main() {
  defer abc.Release()
  projectID := "6666"
  // Initialize the SDK, and you only need do this once
  err := abc.Init(context.TODO(), []string{projectID},
  	abc.WithEnvType(env.TypePrd),
 	abc.WithSecretKey("{{.SecretKey}}"))
  if err != nil {
    log.Errorf("Init fail:%v", err)
    return
  }
  abcUserContext := abc.NewUserContext("unitID_demo")
  experiment, err := abcUserContext.GetExperiment(context.TODO(), projectID, "abc_layer_name")
  if err != nil {
    log.Errorf("GetExperiment fail:%v", err)
    return
  }
  log.Infof("%v", experiment.GetBoolWithDefault("should_show_banner", false))
  featureFlag, err = abcUserContext.GetFeatureFlag(context.TODO(), projectID, "new_feature_flag")
  flagValue = featureFlag.MustGetBool()
  log.Infof("the feature flag value is %v", flagValue)
}

```
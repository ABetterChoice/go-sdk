// Package abc TODO
package abc

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/abetterchoice/go-sdk/internal"
	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/config"
	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/pkg/errors"
)

// GetValueByVariantKey retrieves the parameter value using the globally unique parameter key.
//
// Experiments that carry this key, including config experiments, are resolved first.
// Whitelist wins over audience. If neither hits, GetRemoteConfig resolves config whitelist,
// then config conditions and the default, without entering an experiment.
func (c *userContext) GetValueByVariantKey(ctx context.Context, projectID string, key string,
	opts ...ExperimentOption) (*ValueResult, error) {
	options := defaultExperimentOptions // 拷贝，defaultExperimentOptions 作为模板保持不变
	if c.err != nil {
		return nil, c.err
	}
	c.fillOption(&options)
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, errors.Wrap(err, "opt")
		}
	}
	layerKeys, err := experiment.Executor.VariantKey2LayerKey(projectID, key)
	if err != nil {
		return nil, errors.Wrapf(err, "VariantKey2LayerKey")
	}
	vr := &ValueResult{
		Value: &Value{},
		Detail: &valueDetail{
			LayerKeys: layerKeys,
		},
	}
	exp, err := experiment.Executor.SelectVariantExperiment(ctx, projectID, key, &options)
	if err != nil {
		return nil, errors.Wrap(err, "selectVariantExperiment")
	}
	if exp != nil {
		return applyExperimentValue(vr, key, exp), nil
	}
	resolvedKey := resolveRemoteConfigKey(key, internal.C.TokenDataEnv)
	if remoteConfigExists(projectID, resolvedKey) {
		configValue, err := config.Executor.GetRemoteConfig(ctx, projectID, resolvedKey, &options)
		if err != nil {
			return nil, errors.Wrap(err, "GetRemoteConfig")
		}
		return applyConfigValue(vr, key, configValue), nil
	}
	configResult, err := c.GetRemoteConfig(ctx, projectID, key, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "GetRemoteConfig")
	}
	vr.Value = configResult.Value
	vr.Detail.ConfigKey = key
	return vr, nil
}

func remoteConfigExists(projectID, key string) bool {
	application := cache.GetApplication(projectID)
	if application == nil || application.TabConfig == nil || application.TabConfig.ConfigData == nil {
		return false
	}
	remoteConfig, ok := application.TabConfig.ConfigData.RemoteConfigIndex[key]
	return ok && remoteConfig != nil
}

func applyConfigValue(vr *ValueResult, key string, configValue *config.Value) *ValueResult {
	vr.data = configValue.Data
	if configValue.Experiment == nil {
		vr.Detail.ConfigKey = key
		return vr
	}
	return applyExperimentValue(vr, key, configValue.Experiment)
}

func applyExperimentValue(vr *ValueResult, key string, exp *experiment.Experiment) *ValueResult {
	if param, ok := exp.Params[key]; ok {
		vr.data = []byte(param)
	} else if len(vr.data) == 0 {
		vr.data = []byte(exp.Params[key])
	}
	vr.Detail.GroupKey = exp.GroupKey
	vr.Detail.VariantID = exp.Id
	vr.Detail.ExperimentID = exp.ExperimentId
	vr.Detail.ExperimentKey = exp.ExperimentKey
	vr.Detail.LayerKey = exp.LayerKey
	return vr
}

// ValueResult TODO
type ValueResult struct {
	*Value
	Detail *valueDetail
}

type valueDetail struct {
	ExperimentID  int64    // 非零说明走了实验获取参数值，命中的具体实验 ID
	ExperimentKey string   // 非空说明走了实验获取参数值，命中的具体实验
	VariantID     int64    // 非零说明走了实验获取参数值，命中的具体 variant（实验组）ID
	GroupKey      string   // 非空说明走了实验获取参数值 命中的具体实验版本（= variant key）
	LayerKeys     []string // 非空说明走了实验获取参数值，参数可能在多个互斥层上
	LayerKey      string   // 非空说明最终命中的实验层
	ConfigKey     string   // 非空说明走了配置获取参数值
}

// Value Parameter Value
type Value struct {
	data []byte
}

// Bytes Get specific configuration data. The original data is a snapshot of the local cache. To avoid tampering, a new copy of the data is copied here each time.
func (v *Value) Bytes() []byte {
	var result = make([]byte, len(v.data))
	copy(result, v.data)
	return result
}

// String Get specific configuration data, character value
func (v *Value) String() string {
	return string(v.data)
}

// GetInt64 Get int64 type value
func (v *Value) GetInt64() (int64, error) {
	return strconv.ParseInt(v.String(), 10, 64)
}

// GetInt64WithDefault Gets an int value, returning a default value if it fails.
func (v *Value) GetInt64WithDefault(defaultValue int64) int64 {
	result, err := v.GetInt64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetInt64 Get specific configuration data and convert it to int64 type. If the conversion fails, ignore it.
func (v *Value) MustGetInt64() int64 {
	result, _ := strconv.ParseInt(v.String(), 10, 64)
	return result
}

// MustGetBool Get specific configuration data, force conversion layer boolean type, ignore if forced conversion fails
func (v *Value) MustGetBool() bool {
	flag, _ := strconv.ParseBool(v.String())
	return flag
}

// GetBool Get Boolean type
func (v *Value) GetBool() (bool, error) {
	return strconv.ParseBool(v.String())
}

// GetBoolWithDefault Gets a boolean value, returning a default value if it fails
func (v *Value) GetBoolWithDefault(defaultValue bool) bool {
	result, err := v.GetBool()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64 Get float64 type data, if failed, return defaultValue
func (v *Value) GetFloat64() (float64, error) {
	return strconv.ParseFloat(v.String(), 64)
}

// GetFloat64WithDefault Get float64 type data, if failed, return defaultValue
func (v *Value) GetFloat64WithDefault(defaultValue float64) float64 {
	result, err := v.GetFloat64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetFloat64 Get float64 type data, return zero value if error occurs
func (v *Value) MustGetFloat64() float64 {
	result, _ := v.GetFloat64()
	return result
}

// GetJSONMap Get json map type data
func (v *Value) GetJSONMap() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	err := json.Unmarshal(v.data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJSONMapWithDefault Get json map type data, if failed, return defaultValue
func (v *Value) GetJSONMapWithDefault(defaultValue map[string]interface{}) map[string]interface{} {
	result, err := v.GetJSONMap()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetJSONMap Get json map type data, return zero value if error occurs
func (v *Value) MustGetJSONMap() map[string]interface{} {
	result, _ := v.GetJSONMap()
	return result
}

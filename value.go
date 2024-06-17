// Package abc TODO
package abc

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/pkg/errors"
)

// GetValueByVariantKey retrieves the parameter value using the globally unique parameter key.
//
//	// It first attempts to find the value among all parameters under all experiments within the given project ID.
//	// If the parameter value is not found, it continues to search for this parameter among all feature flags
//	// under the project ID until the value is found.
//	//
//	// Parameters:
//	// ctx: The context for the operation.
//	// projectID: The ID of the project where the search is performed.
//	// key: The globally unique key that identifies the parameter.
//	// opts: Additional experiment options.
//	//
//	// Returns:
//	// A pointer to the ValueResult object containing the found value and an error, if any occurred during the search.
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
	if len(layerKeys) != 0 { // 参数正在实验层上，获取实验命中的参数值
		opts = append(opts, WithLayerKeyList(layerKeys))
		experimentResult, err := c.GetExperiments(ctx, projectID, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "GetExperiment")
		}
		if len(experimentResult.Data) > 1 {
			return nil, errors.Errorf("invalid variantKey:%s, used by multiple layers:%+v", key, layerKeys)
		}
		for layerKey, group := range experimentResult.Data {
			data, ok := group.GetBytes(key)
			if !ok { // 获取层默认值
				data, err = experiment.Executor.GetVariantValue(projectID, layerKey, key)
				if err != nil {
					return nil, errors.Wrap(err, "getVariantValue")
				}
			}
			vr.data = data
			vr.Detail.GroupKey = group.Key
			vr.Detail.ExperimentKey = group.ExperimentKey
			vr.Detail.LayerKey = layerKey
			return vr, nil
		}
		return vr, nil
	}
	configResult, err := c.GetRemoteConfig(ctx, projectID, key, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "GetRemoteConfig")
	}
	vr.Value = configResult.Value
	vr.Detail.ConfigKey = key
	return vr, nil
}

// ValueResult TODO
type ValueResult struct {
	*Value
	Detail *valueDetail
}

type valueDetail struct {
	ExperimentKey string   // 非空说明走了实验获取参数值，命中的具体实验
	GroupKey      string   // 非空说明走了实验获取参数值 命中的具体实验版本
	LayerKeys     []string // 非空说明走了实验获取参数值，参数可能在多个互斥层上
	LayerKey      string   // 非空说明最终命中的实验层
	ConfigKey     string   // 非空说明走了配置获取参数值
}

// Value 参数值
type Value struct {
	data []byte
}

// Bytes 获取具体的配置数据，原始数据是本地缓存的一份快照，为避免被篡改，这里每次都会 copy 一份新的数据
func (v *Value) Bytes() []byte {
	var result = make([]byte, len(v.data))
	copy(result, v.data)
	return result
}

// String 获取具体的配置数据，字符值
func (v *Value) String() string {
	return string(v.data)
}

// GetInt64 获取 int64 类型值
func (v *Value) GetInt64() (int64, error) {
	return strconv.ParseInt(v.String(), 10, 64)
}

// GetInt64WithDefault 获取 int 类型值，如果失败，则返回默认值
func (v *Value) GetInt64WithDefault(defaultValue int64) int64 {
	result, err := v.GetInt64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetInt64 获取具体的配置数据，强转成 int64 类型，强转失败忽略
func (v *Value) MustGetInt64() int64 {
	result, _ := strconv.ParseInt(v.String(), 10, 64)
	return result
}

// MustGetBool 获取具体的配置数据，强转层 boolean 类型，强转失败忽略
func (v *Value) MustGetBool() bool {
	flag, _ := strconv.ParseBool(v.String())
	return flag
}

// GetBool 获取布尔类型
func (v *Value) GetBool() (bool, error) {
	return strconv.ParseBool(v.String())
}

// GetBoolWithDefault 获取布尔值，如果失败则返回默认值
func (v *Value) GetBoolWithDefault(defaultValue bool) bool {
	result, err := v.GetBool()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64 获取 float64 类型数据，如果失败则返回 defaultValue
func (v *Value) GetFloat64() (float64, error) {
	return strconv.ParseFloat(v.String(), 64)
}

// GetFloat64WithDefault 获取 float64 类型数据，如果失败则返回 defaultValue
func (v *Value) GetFloat64WithDefault(defaultValue float64) float64 {
	result, err := v.GetFloat64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetFloat64 获取 float64 类型数据，报错返回零值
func (v *Value) MustGetFloat64() float64 {
	result, _ := v.GetFloat64()
	return result
}

// GetJSONMap 获取 json map 类型数据
func (v *Value) GetJSONMap() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	err := json.Unmarshal(v.data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJSONMapWithDefault 获取 json map 类型数据，如果失败则返回 defaultValue
func (v *Value) GetJSONMapWithDefault(defaultValue map[string]interface{}) map[string]interface{} {
	result, err := v.GetJSONMap()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetJSONMap 获取 json map 类型数据，报错返回零值
func (v *Value) MustGetJSONMap() map[string]interface{} {
	result, _ := v.GetJSONMap()
	return result
}

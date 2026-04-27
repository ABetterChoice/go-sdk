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
// Resolution order:
//  1. If key is bound to one or more experiment layers, walk those layers in
//     ascending order of each layer's earliest experiment ID. Use the first
//     layer where the user is hit by a real experiment group (IsDefault=false)
//     AND that group has configured this key.
//  2. Otherwise (key is not on any layer, or no layer satisfies the condition
//     above), fall back to the remote config resolution flow
//     (override list / holdout / condition list / static default value).
//
// Parameters:
//
//	ctx: The context for the operation.
//	projectID: The ID of the project where the search is performed.
//	key: The globally unique key that identifies the parameter.
//	opts: Additional experiment options.
//
// Returns:
//
//	A pointer to the ValueResult object containing the resolved value, and an
//	error if any occurred during the lookup.
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
<<<<<<< HEAD
	if len(layerKeys) != 0 {
		// layerKeys 已按层上"最早实验 ID"升序排序（近似为最早实验的创建顺序）。
		// 优先使用第一个"命中真实验且实验组配置了 key"的层；
		// 任意层都不满足时 fallback 到参数本身的远程配置流程
		// （白名单 / holdout / condition / 参数静态默认值）。
		experimentOpts := make([]ExperimentOption, 0, len(opts)+1)
		experimentOpts = append(experimentOpts, opts...)
		experimentOpts = append(experimentOpts, WithLayerKeyList(layerKeys))
		experimentResult, err := c.GetExperiments(ctx, projectID, experimentOpts...)
		if err != nil {
			return nil, errors.Wrap(err, "GetExperiment")
		}
		for _, layerKey := range layerKeys {
			group, ok := experimentResult.Data[layerKey]
			if !ok || group == nil || group.IsDefault {
				continue
			}
			data, ok := group.GetBytes(key)
			if !ok {
				// 实验组没配该 key，视为该层在该 key 上未命中，继续下一层
				continue
			}
			vr.data = data
			vr.Detail.GroupKey = group.Key
			vr.Detail.VariantID = group.ID
			vr.Detail.ExperimentID = group.ExperimentID
			vr.Detail.ExperimentKey = group.ExperimentKey
			vr.Detail.LayerKey = layerKey
			return vr, nil
		}
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

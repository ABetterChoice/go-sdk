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
	if len(layerKeys) != 0 { // The parameter is on the experimental layer, get the parameter value hit by the experiment
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
	ExperimentKey string   // Non-empty means the experiment was completed to obtain parameter values, and the specific experiment was hit
	GroupKey      string   // Non-empty means the experiment was completed to obtain parameter values. The specific experimental version hit
	LayerKeys     []string // Non-empty means the experiment was completed to obtain parameter values. The parameters may be on multiple mutually exclusive layers
	LayerKey      string   // Non-empty means the experimental layer that was finally hit
	ConfigKey     string   // Non-empty means the configuration was completed to obtain parameter values
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

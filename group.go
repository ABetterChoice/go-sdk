// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval, user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"encoding/json"
	"strconv"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/protoc_cache_server"
)

// ExperimentList Experimental group list,
// which stores the return data of GetExperiments and identifies the experimental groups hit by unitID in each layer.
type ExperimentList struct {
	userCtx *userContext // Store user information
	// The experimental group hit by unitID in each layer, the key is layerKey,
	// and the value is the experimental group hit under the layer.
	Data map[string]*Group
}

// ExperimentResult Experimental offloading results,
// subsequent exposure reporting can be directly passed into the ExperimentResult instance
type ExperimentResult struct {
	userCtx *userContext // Store user information
	*Group               // Experimental group results for specific hits
}

// Group Experimental group information
type Group struct {
	// Experimental group ID Note: The database ID is not disclosed in principle, but for historical reasons
	ID int64 `json:"id"`

	// group Key
	Key string `json:"key"`

	// Experiment Key
	ExperimentKey string `json:"experimentKey"`

	// layer key
	LayerKey string `json:"layerKey"`

	// Whether it is a default experiment, layer default experiment or global default experiment is true
	IsDefault bool `json:"isDefault"`

	// Whether it is the control group, if not, it is the experimental group
	IsControl bool `json:"isControl"`

	// Whether it is an experimental group hit by the whitelist. If the experimental parameters conflict,
	// the whitelist takes precedence.
	IsOverrideList bool `json:"isOverrideList"`

	// Experimental parameters are not provided directly here to prevent concurrent reading and writing of the map,
	// while ensuring performance and avoiding unnecessary memory copies.
	// Provided externally through strongly typed API, such as [T]GetNumber(key string) T
	params map[string]string

	// Scene ID list, scene ID description
	sceneIDList []int64

	// Account system
	UnitIDType  protoc_cache_server.UnitIDType `json:"unitIdType"`
	holdoutData map[string]*Group
}

// SceneIDList Get scene ID list, deep copy
func (g *Group) SceneIDList() []int64 {
	if len(g.sceneIDList) == 0 {
		return nil
	}
	var result = make([]int64, 0, len(g.sceneIDList))
	for _, sceneID := range g.sceneIDList {
		result = append(result, sceneID)
	}
	return result
}

// Params deep copy params to prevent concurrent reading and writing of map
// At the same time, the sdk provides a variety of strongly typed APIs for easy use params
func (g *Group) Params() map[string]string {
	var result = make(map[string]string, len(g.params))
	for k, v := range g.params {
		result[k] = v
	}
	return result
}

// GetBool gets bool type data
func (g *Group) GetBool(key string) (bool, error) {
	source, ok := g.params[key]
	if !ok {
		return false, env.ErrParamKeyNotFound
	}
	return strconv.ParseBool(source)
}

// GetInt64 gets Int64 type data
func (g *Group) GetInt64(key string) (int64, error) {
	source, ok := g.params[key]
	if !ok {
		return 0, env.ErrParamKeyNotFound
	}
	return strconv.ParseInt(source, 10, 64)
}

// GetFloat64 gets float64 type data
func (g *Group) GetFloat64(key string) (float64, error) {
	source, ok := g.params[key]
	if !ok {
		return 0, env.ErrParamKeyNotFound
	}
	return strconv.ParseFloat(source, 64)
}

// GetJSONMap gets json map type data
func (g *Group) GetJSONMap(key string) (map[string]interface{}, error) {
	source, ok := g.params[key]
	if !ok {
		return nil, env.ErrParamKeyNotFound
	}
	var result = make(map[string]interface{})
	err := json.Unmarshal([]byte(source), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetString gets string type data
func (g *Group) GetString(key string) (string, error) {
	source, ok := g.params[key]
	if !ok {
		return "", env.ErrParamKeyNotFound
	}
	return source, nil
}

// MustGetBytes gets bytes, returns an empty array if it does not exist
func (g *Group) MustGetBytes(key string) []byte {
	source, _ := g.params[key]
	return []byte(source)
}

// GetBytes gets bytes
func (g *Group) GetBytes(key string) ([]byte, bool) {
	source, ok := g.params[key]
	if ok {
		return []byte(source), true
	}
	return nil, false
}

// GetBoolWithDefault gets bool type data, if it fails, it returns defaultValue
func (g *Group) GetBoolWithDefault(key string, defaultValue bool) bool {
	result, err := g.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetInt64WithDefault gets Int64 type data, and returns defaultValue if it fails.
func (g *Group) GetInt64WithDefault(key string, defaultValue int64) int64 {
	result, err := g.GetInt64(key)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64WithDefault gets float64 type data, and returns defaultValue if it fails.
func (g *Group) GetFloat64WithDefault(key string, defaultValue float64) float64 {
	result, err := g.GetFloat64(key)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetJSONMapWithDefault gets json map type data, and returns defaultValue if it fails.
func (g *Group) GetJSONMapWithDefault(key string,
	defaultValue map[string]interface{}) map[string]interface{} {
	result, err := g.GetJSONMap(key)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetStringWithDefault gets string type data, and returns defaultValue if it fails.
func (g *Group) GetStringWithDefault(key string, defaultValue string) string {
	result, err := g.GetString(key)
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetBool gets bool type data, returns zero value when error occurs
func (g *Group) MustGetBool(key string) bool {
	result, _ := g.GetBool(key)
	return result
}

// MustGetInt64 obtains Int64 type data and returns zero value if an error occurs.
func (g *Group) MustGetInt64(key string) int64 {
	result, _ := g.GetInt64(key)
	return result
}

// MustGetFloat64 obtains float64 type data and returns zero value when reporting an error.
func (g *Group) MustGetFloat64(key string) float64 {
	result, _ := g.GetFloat64(key)
	return result
}

// MustGetJSONMap obtains json map type data and returns zero value when reporting an error.
func (g *Group) MustGetJSONMap(key string) map[string]interface{} {
	result, _ := g.GetJSONMap(key)
	return result
}

// MustGetString obtains string type data and returns zero value if an error occurs.
func (g *Group) MustGetString(key string) string {
	result, _ := g.GetString(key)
	return result
}

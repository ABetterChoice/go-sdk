// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval, user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	"github.com/abetterchoice/go-sdk/internal/config"
	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/abetterchoice/go-sdk/plugin/log"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
)

// GetRemoteConfig gets the remote configuration logic development
// GetConfig gets the configuration hit by the user, specifies the projectID and configuration key
func (c *userContext) GetRemoteConfig(ctx context.Context, projectID string, key string,
	opts ...ConfigOption) (result *ConfigResult, err error) {
	options := defaultExperimentOptions // Copy, defaultExperimentOptions remains unchanged as template
	defer func(startTime time.Time) {
		latency := time.Since(startTime)
		if options.IsExposureLoggingAutomatic && !internal.C.IsDisableReport {
			exposureErr := asyncExposureRemoteConfig(projectID, result, protoc_event_server.ExposureType_EXPOSURE_TYPE_AUTOMATIC)
			if exposureErr != nil {
				log.Errorf("[projectID=%v]asyncExposureRemoteConfig fail:%v", projectID, exposureErr)
			}
		}
		exposureErr := asyncExposureRemoteConfigEvent(projectID, result, latency, env.JSONString(&options), err)
		if exposureErr != nil {
			log.Errorf("[projectID=%v]exposureRemoteConfigEvent fail:%v", projectID, exposureErr)
		}
	}(time.Now())
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
	configValue, err := config.Executor.GetRemoteConfig(ctx, projectID, key, &options)
	if err != nil {
		return nil, err
	}
	return &ConfigResult{
		userCtx: c,
		Config: &Config{
			Key:            key,
			Value:          &Value{data: configValue.Data},
			IsOverrideList: configValue.IsOverrideList,
			IsDefault:      configValue.IsDefault,
			Experiment:     convertGroup2Experiment(configValue.Experiment),
			remoteConfig:   configValue.RemoteConfig,
			unitIDType:     configValue.UnitIDType,
		},
	}, nil
}

// ConfigOption Gets the relevant Option of the hit configuration, including the specified scene ID
// WithSceneIDConfigOption Or if the configuration is not hit, return zero value or specify a default value
// WithDefaultValueConfigOption
type ConfigOption = ExperimentOption

// WithIsPreparedDMPTagConfigOpt sets whether to preprocess DMP tags
// If there is no dmp tag configured, or there is only one, preprocessing will not be enabled.
// Otherwise, it is enabled by default, and developers can use this option to manage it manually.
func WithIsPreparedDMPTagConfigOpt(isPreparedDMPTag bool) ConfigOption {
	return func(options *experiment.Options) error {
		options.IsPreparedDMPTag = isPreparedDMPTag
		return nil
	}
}

// WithIsDisableDMPConfigOpt sets whether to close DMP to prevent rpc operations.
// If the DMP label is clearly not needed, it can be closed, such as local testing, etc.
func WithIsDisableDMPConfigOpt(isDisableDMP bool) ConfigOption {
	return func(options *experiment.Options) error {
		options.IsDisableDMP = isDisableDMP
		return nil
	}
}

// ConfigResult TODO
type ConfigResult struct {
	userCtx *userContext `json:"-"`
	*Config
}

// Config configuration information
type Config struct {
	*Value

	// Configure key
	Key string `json:"key"`

	// Whether it is a whitelist hit
	IsOverrideList bool `json:"isOverrideList"`

	// Is it the default value?
	IsDefault bool `json:"isDefault"`

	// Configure the bound experiment
	Experiment *Group `json:"experiment"`

	// Remote configuration information is not disclosed to prevent concurrency problems
	// and can provide read-only operations through the API
	remoteConfig *protoccacheserver.RemoteConfig `json:"-"`

	// Account system
	unitIDType protoccacheserver.UnitIDType `json:"-"`
}

// Byte gets the specific configuration data. The original data is a snapshot of the local cache.
// To avoid tampering, a new data will be copied here every time.
func (c *Config) Byte() []byte {
	var result = make([]byte, len(c.data))
	copy(result, c.data)
	return result
}

// String gets specific configuration data, character value
func (c *Config) String() string {
	return string(c.data)
}

// GetInt64 gets the int64 type value
func (c *Config) GetInt64() (int64, error) {
	return strconv.ParseInt(c.String(), 10, 64)
}

// GetInt64WithDefault gets the int type value, if it fails, returns the default value
func (c *Config) GetInt64WithDefault(defaultValue int64) int64 {
	result, err := c.GetInt64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetInt64 gets the specific configuration data and force-converts it to int64 type.
// If the forced conversion fails, it will be ignored.
func (c *Config) MustGetInt64() int64 {
	result, _ := strconv.ParseInt(c.String(), 10, 64)
	return result
}

// MustGetBool obtains specific configuration data, the forced transfer layer is of boolean type,
// and the forced transfer failure is ignored.
func (c *Config) MustGetBool() bool {
	flag, _ := strconv.ParseBool(c.String())
	return flag
}

// GetBool gets the Boolean type
func (c *Config) GetBool() (bool, error) {
	return strconv.ParseBool(c.String())
}

// GetBoolWithDefault gets the Boolean value and returns the default value if it fails
func (c *Config) GetBoolWithDefault(defaultValue bool) bool {
	result, err := c.GetBool()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64 gets float64 type data, and returns defaultValue if it fails.
func (c *Config) GetFloat64() (float64, error) {
	return strconv.ParseFloat(c.String(), 64)
}

// GetFloat64WithDefault gets float64 type data, and returns defaultValue if it fails.
func (c *Config) GetFloat64WithDefault(defaultValue float64) float64 {
	result, err := c.GetFloat64()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetFloat64 obtains float64 type data and returns zero value when reporting an error.
func (c *Config) MustGetFloat64() float64 {
	result, _ := c.GetFloat64()
	return result
}

// GetJSONMap gets json map type data
func (c *Config) GetJSONMap() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	err := json.Unmarshal(c.data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJSONMapWithDefault gets json map type data, and returns defaultValue if it fails.
func (c *Config) GetJSONMapWithDefault(defaultValue map[string]interface{}) map[string]interface{} {
	result, err := c.GetJSONMap()
	if err != nil {
		return defaultValue
	}
	return result
}

// MustGetJSONMap obtains json map type data and returns zero value when reporting an error.
func (c *Config) MustGetJSONMap() map[string]interface{} {
	result, _ := c.GetJSONMap()
	return result
}

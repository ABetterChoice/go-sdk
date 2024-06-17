// Package internal sdk
package internal

import (
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/protoc_cache_server"
)

// GlobalConfig Global configuration, including global configuration after storage Init,
// including whether to enable reporting, reporting component selection, whether to pre-process DMP, etc.
type GlobalConfig struct {
	// Initialize the incoming projectID list and cache all configurations of the business locally
	ProjectIDList []string `json:"projectIdList"`
	// Whether to disable reporting, default is false
	IsDisableReport bool `json:"isEnableReport"`
	// Environment type, default is prd
	EnvType env.Type `json:"envType"`
	// Whether to disable dmp, global configuration; each diversion inherits the global configuration by default,
	// can be set independently, default false
	IsDisableDMP bool `json:"isDisableDmp"`
	// Monitoring reporting component initialization parameters, key is the plugin name,
	// value is the initialization parameter
	MetricsPluginInitConfig map[string]*protoc_cache_server.MetricsInitConfig `json:"metricsPluginInitConfig"`
	// Whether to customize the cache service plug-in. If not, the default is to use the TAB cache service
	IsCustomCacheClient bool `json:"isCustomCacheClient"`
	// Whether to customize the DMP user portrait service plug-in. If not,
	// the default is to use the TAB DMP user portrait service
	IsCustomDMPClient bool `json:"isCustomDmpClient"`
	// Region information, supports sending different configurations to different regions,
	// such as different reporting addresses for different regions
	RegionCode string `json:"regionCode"`
	// secretKey, used for authentication
	SecretKey string `json:"secretKey"`
}

// C global configuration related instances, no need to lock,
// the instance will only be modified during Init/Release,
// Init is protected by Once, Release is not concurrently safe, user notice
var C = &GlobalConfig{}

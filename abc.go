// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval,
// user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/client"
	mp "github.com/abetterchoice/go-sdk/plugin/metrics"
	_ "github.com/abetterchoice/metrics-pubsub" // metrics-pubsub TODO
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/pkg/errors"
)

// Init initializes the ABC SDK.
// This process retrieves configuration data from remote servers, storing it in a local cache,
// and carries out operations like initializing event logging components.
// The primary logic for splitting experiment traffic is computed locally.
// The Init() function only needs to be called once, with subsequent calls returning nil directly.
// If re-initialization is required, please call Release() first.
// Timeout control is facilitated through the provided ctx parameter.
// Additional control configurations can be supplied as needed via InitOption
func Init(ctx context.Context, projectIDList []string, opts ...InitOption) (err error) {
	defer func(start time.Time) {
		manualInitEvent(projectIDList, time.Since(start), err)
		if err != nil {
			Release()
		}
	}(time.Now())
	if len(projectIDList) == 0 {
		err = fmt.Errorf("projectIDList is required")
		return
	}
	once.Do(func() {
		var c = &internal.GlobalConfig{
			ProjectIDList: projectIDList, EnvType: env.TypePrd,
			MetricsPluginInitConfig: map[string]*protoccacheserver.MetricsInitConfig{}}
		for _, opt := range opts {
			err = opt(c)
			if err != nil {
				return
			}
		}
		internal.C = c
		if !c.IsCustomCacheClient {
			client.RegisterCacheClient(client.NewTABCacheClient(client.WithEnvType(c.EnvType)))
		}
		if !c.IsCustomDMPClient {
			client.RegisterDMPClient(client.NewDMPClient(client.WithEnvTypeOption(c.EnvType)))
		}
		initExposureConsumer()
		err = initCustomMetricsPlugin(ctx, c)
		if err != nil {
			return
		}
		err = cache.InitLocalCache(ctx, projectIDList)
		if err != nil {
			return
		}
		err = initMetricsPlugin(ctx, c)
		if err != nil {
			return
		}
		return
	})
	return err
}

// Release local cache, concurrency is not safe
func Release() {
	cache.Release()
	once = sync.Once{}
	internal.C = &internal.GlobalConfig{}
}

var (
	once sync.Once
)

// initCustomMetricsPlugin TODO
// initialize the registered monitoring and reporting plug-ins one by one
func initCustomMetricsPlugin(ctx context.Context, config *internal.GlobalConfig) error {
	// traverse all registered monitoring and reporting plug-ins
	return mp.WalkFunc(func(name string, client mp.Client) error {
		initConfig, ok := config.MetricsPluginInitConfig[name]
		if !ok {
			return nil
		}
		// If the user explicitly passes in the initialization parameters,
		// the user-defined passed parameters will be used directly.
		err := client.Init(ctx, initConfig)
		if err != nil {
			return errors.Wrapf(err, "init metrics plugin [%v]", name)
		}
		return nil
	})
}

// initMetricsPlugin TODO
// Initialize monitoring plugins provided by remote configuration.
func initMetricsPlugin(ctx context.Context, config *internal.GlobalConfig) error {
	return mp.WalkFunc(func(name string, client mp.Client) error {
		// traverse all projectIDs and initialize related metrics plugin
		// different projectIDs may have the same metrics plugin, and the same plugin may be initialized multiple times.
		// it is necessary to ensure that the monitoring and reporting initConfig of projectIDList is consistent.
		// if inconsistent, the initConfig will be randomly initialized.
		for _, projectID := range internal.C.ProjectIDList {
			application := cache.GetApplication(projectID)
			if application == nil {
				continue
			}
			for pluginName, initConfig := range application.MetricsPluginInitConfigIndex {
				if name != pluginName {
					continue
				}
				if _, ok := config.MetricsPluginInitConfig[pluginName]; ok {
					continue // The custom plug-in has been initialized and does not need to be initialized again.
				}
				initConfigNew, err := deepCopyMetricsInitConfig(initConfig)
				if err != nil {
					return err
				}
				if initConfigNew.Kv == nil {
					initConfigNew.Kv = make(map[string]string)
				}
				if _, ok := initConfigNew.Kv[mp.InitConfigKvToken]; !ok {
					initConfigNew.Kv[mp.InitConfigKvToken] = config.SecretKey
				}
				// Initialize the client according to the initialization parameters of the remote configuration
				return client.Init(ctx, initConfigNew)
			}
		}
		return nil
	})
}

// InitOption Initialization Option is used to customize and control more fine-grained configurations,
// such as whether to enable reporting, setting environment, RPC protocol, back-end cache service address,
// socket5 proxy, etc.
type InitOption func(config *internal.GlobalConfig) error

// WithRegisterMetricsPlugin register monitoring and reporting plug-in
func WithRegisterMetricsPlugin(client mp.Client, initConfig *protoccacheserver.MetricsInitConfig) InitOption {
	return func(config *internal.GlobalConfig) error {
		if client == nil {
			return fmt.Errorf("client should not be nil")
		}
		mp.RegisterClient(client)
		// it is used for subsequent initialization.
		// some initialization may rely on cached data,
		// so the monitoring and reporting component is initialized after the cache is successfully pulled.
		config.MetricsPluginInitConfig[client.Name()] = initConfig
		return nil
	}
}

// WithRegisterCacheClient register the background cache service interface implementation,
// which can replace the default TAB background cache service
func WithRegisterCacheClient(c client.Client) InitOption {
	return func(config *internal.GlobalConfig) error {
		if c == nil {
			return errors.Errorf("client is required")
		}
		client.CacheClient = c
		config.IsCustomCacheClient = true
		return nil
	}
}

// WithRegisterDMPClient register the dmp user portrait service interface implementation,
// which can replace the default TAB user portrait service
func WithRegisterDMPClient(dmpClient client.DMPClient) InitOption {
	return func(config *internal.GlobalConfig) error {
		if dmpClient == nil {
			return errors.Errorf("dmpClient is required")
		}
		client.DC = dmpClient
		config.IsCustomDMPClient = true
		return nil
	}
}

// WithSecretKey pass in secretKey for backend authentication use
func WithSecretKey(secretKey string) InitOption {
	return func(config *internal.GlobalConfig) error {
		config.SecretKey = secretKey
		return nil
	}
}

// WithEnvType set environment, default official environment
func WithEnvType(envType env.Type) InitOption {
	return func(config *internal.GlobalConfig) error {
		config.EnvType = envType
		return nil
	}
}

// WithDisableReport disable reporting. even if the plug-in is registered, reporting will not be executed.
func WithDisableReport(isDisable bool) InitOption {
	return func(config *internal.GlobalConfig) error {
		config.IsDisableReport = isDisable
		return nil
	}
}

// WithRegionCode set region code to support different regions delivering different configurations
func WithRegionCode(regionCode string) InitOption {
	return func(config *internal.GlobalConfig) error {
		config.RegionCode = regionCode
		return nil
	}
}

// GetGlobalConfig returns the global configuration object,
// including the projectID passed in Init, whether to enable exposure reporting, etc., deep copy
// modifying the returned globalConfig will not update the global configuration, it is only used as a data query
func GetGlobalConfig() (*internal.GlobalConfig, error) {
	return deepCopyGlobalConfig(internal.C)
}

func deepCopyGlobalConfig(source *internal.GlobalConfig) (*internal.GlobalConfig, error) {
	var result = &internal.GlobalConfig{}
	data, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func deepCopyMetricsInitConfig(source *protoccacheserver.MetricsInitConfig) (*protoccacheserver.MetricsInitConfig,
	error) {
	var result = &protoccacheserver.MetricsInitConfig{}
	data, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RegisterProjectIDs register a new projectID and support multiple registrations
func RegisterProjectIDs(ctx context.Context, projectIDList []string) error {
	return cache.InitLocalCache(ctx, projectIDList)
}

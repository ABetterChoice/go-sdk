// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval, user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/abetterchoice/go-sdk/plugin/log"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
)

var (
	// defaultExperimentOptions global defaults, as templates remain unchanged,
	// abtest shunting process is passed by value without using references
	// the objects here cannot have pointer instances, otherwise there will be concurrency problems.
	// please pay attention to modifications.
	defaultExperimentOptions = experiment.Options{
		IsExposureLoggingAutomatic: true,
		IsPreparedDMPTag:           false,
		IsDisableDMP:               false,
	}
)

// GetExperiment Assignment is the recommended method for retrieving experiment assignments.
// This API checks which experiment and group the unit ID from the provided context belongs to,
// under a specified layer within a given projectID. If the unit ID is part of a specific experiment group,
// the API returns the corresponding experiment and group information.
// If not, the API returns an empty experiment. Subsequent calls to the parameter retrieval API
// with this empty experiment will return the layer's default parameter value.
//
// For additional information, visit https://abetterchoice.ai/docs?wj-docs=%2Fprod%2Fdocs%2Fgo_sdk.
//
// Note: You can also specify the layerKey in the input parameter opts. However, if layerKey is directly inputted,
// it will take precedence over the one defined in opts.
//
// For more filter conditions, refer to the ExperimentOption definition.
func (c *userContext) GetExperiment(ctx context.Context, projectID string, layerKey string,
	opts ...ExperimentOption) (*ExperimentResult, error) {
	// if the layerKey is specified, the layerKeys in options will also be integrated.
	// this will integrate layerKeys, sceneIDs, experimentKeys in options, and relationships
	opts = append(opts, WithLayerKey(layerKey))
	// the underlying implementation is based on GetExperiments
	experimentList, err := c.GetExperiments(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}
	e, ok := experimentList.Data[layerKey]
	if ok && e != nil {
		return &ExperimentResult{
			userCtx: c,
			Group:   e,
		}, nil
	}
	return nil, nil
}

// GetExperiments is a batch version of GetExperiment(). It returns the experiment assignments
// for all experiment layers under the specified project ID. This API is typically used when users
// need to fetch all experiment results and pass them to downstream services.
//
// However, using this API is not recommended as it may lead to logging a large number of unexpected exposures,
// which could dilute your results. To mitigate this issue, exposures are not logged automatically by default.
// Instead, you may need to use the exposure logging API to manually log the exposures.
func (c *userContext) GetExperiments(ctx context.Context, projectID string,
	opts ...ExperimentOption) (result *ExperimentList, err error) {
	options := defaultExperimentOptions // copy, defaultExperimentOptions as template remains unchanged
	defer func(startTime time.Time) {
		latency := time.Since(startTime)
		if options.IsExposureLoggingAutomatic && !internal.C.IsDisableReport {
			exposureErr := asyncExposureExperiments(projectID, result, protoc_event_server.ExposureType_EXPOSURE_TYPE_AUTOMATIC)
			if exposureErr != nil {
				log.Errorf("[projectID=%v]asyncExposureExperiments fail:%v", projectID, exposureErr)
			}
		}
		exposureErr := asyncExposureExperimentEvent(projectID, result, latency, env.JSONString(&options), err)
		if exposureErr != nil {
			log.Errorf("[projectID=%v]asyncExposureExperimentEvent fail:%v", projectID, exposureErr)
		}
	}(time.Now())
	if c.err != nil {
		return nil, c.err
	}
	c.fillOption(&options)
	for _, opt := range opts {
		err = opt(&options)
		if err != nil {
			return nil, errors.Wrap(err, "opt")
		}
	}
	experimentList, err := experiment.Executor.GetExperiments(ctx, projectID, &options)
	if err != nil {
		return nil, err // the error here does not need to be wrapped, it is all GetExperiments
	}
	result = &ExperimentList{
		Data: make(map[string]*Group, len(experimentList)),
	}
	for layerKey, group := range experimentList {
		if group == nil {
			continue
		}
		result.Data[layerKey] = convertGroup2Experiment(group)
	}
	for layerKey, holdoutGroup := range options.HoldoutLayerResult {
		if holdoutGroup == nil {
			continue
		}
		result.Data[layerKey] = convertGroup2Experiment(holdoutGroup)
	}
	result.userCtx = c
	return result, nil
}

// GetDefaultExperiments The default experiment on the acquisition layer does not involve any diversion process,
// and is mainly aimed at obtaining a bottom-up strategy.
func GetDefaultExperiments(ctx context.Context, projectID string, opts ...ExperimentOption) (*ExperimentList, error) {
	options := defaultExperimentOptions // copy, defaultExperimentOptions as template remains unchanged
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, errors.Wrap(err, "opt")
		}
	}
	experimentList, err := experiment.Executor.GetDefaultExperiments(ctx, projectID, &options)
	if err != nil {
		return nil, err
	}
	result := &ExperimentList{
		Data: make(map[string]*Group, len(experimentList)),
	}
	for layerKey, group := range experimentList {
		if group == nil {
			continue
		}
		result.Data[layerKey] = convertGroup2Experiment(group)
	}
	return result, nil
}

// fillOption the attributes with userContext into options for subsequent use of abtest offloading.
// fill in the default value. The opt function will overwrite the following default value.
func (c *userContext) fillOption(options *experiment.Options) {
	options.AttributeTag = c.tags
	options.UnitID = c.unitID
	options.DecisionID = c.decisionID
	options.NewUnitID = c.newUnitID
	options.NewDecisionID = c.newDecisionID
	options.DMPTagResult = make(map[string]bool)
	options.HoldoutLayerResult = make(map[string]*experiment.Experiment)
	options.IsDisableDMP = internal.C.IsDisableDMP
}

func convertGroup2Experiment(group *experiment.Experiment) *Group {
	if group == nil {
		return nil
	}
	result := convertGroup2ExperimentWithoutHoldout(group)
	if len(group.HoldoutData) != 0 {
		if result.holdoutData == nil {
			result.holdoutData = make(map[string]*Group)
		}
		for key, holdoutExperiment := range group.HoldoutData {
			g := convertGroup2ExperimentWithoutHoldout(holdoutExperiment)
			result.holdoutData[key] = g
		}
	}
	return result
}

func convertGroup2ExperimentWithoutHoldout(group *experiment.Experiment) *Group {
	return &Group{
		ID:             group.Id,
		Key:            group.GroupKey,
		ExperimentKey:  group.ExperimentKey,
		LayerKey:       group.LayerKey,
		IsDefault:      group.IsDefault,
		IsControl:      group.IsControl,
		IsOverrideList: group.IsOverrideList,
		params:         group.Params,
		UnitIDType:     group.UnitIdType,
		sceneIDList:    group.SceneIdList,
	}
}

// ExperimentOption experimental diversion Options, providing extended control information, including unlimited
// Whether to disable DMP crowd selection
// Whether to turn on TAB automatic record exposure automatic, which is managed by the web backend system by default.
// You can specify whether to turn on automatic for a single abTest diversion, provided that it is during Init
// The reporting component is enabled, please view the relevant code for specific Options
// Each option is an AND relationship, such as passing in layerKeys, sceneIDs,
// and only returning experimental groups that meet the sceneID and layerKey conditions.
// If multiple layers are passed in layerKeys, the experimental group hit on each layerKey layer will be returned.
type ExperimentOption func(options *experiment.Options) error

// WithSceneIDList Set filtering by scene ID, which is equivalent to coloring.
// It only focuses on the hits of experiments in certain scenes.
// The specific scene ID can be specified when creating a layer.
func WithSceneIDList(sceneIDList []int64) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.SceneIDs == nil {
			options.SceneIDs = make(map[int64]bool, len(sceneIDList))
		}
		for i := range sceneIDList {
			options.SceneIDs[sceneIDList[i]] = true
		}
		return nil
	}
}

// WithExperimentKey to filter experiments,
// it is generally recommended to use GetExperiment api and pass in layerKey to filter the layer
// There may be multiple experiments on the layer. If filtered by experiment,
// assuming that the user hits other experiments on the same layer,
// then other experiments hit on the same layer will be returned.
// Deprecated: In short, it is recommended to use WithLayerKey instead of WithExperimentKey
func WithExperimentKey(experimentKey string) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.ExperimentKeys == nil {
			options.ExperimentKeys = make(map[string]bool, 1) // 长度 1
		}
		options.ExperimentKeys[experimentKey] = true
		return nil
	}
}

// WithExperimentKeys to filter experiments,
// it is generally recommended to use GetExperiment api and pass in layerKey to filter the layer
// There may be multiple experiments on the layer. If filtered by experiment,
// assuming that the user hits other experiments on the same layer,
// then other experiments hit on the same layer will be returned.
// Deprecated: In short, it is recommended to use WithLayerKey instead of WithExperimentKey
func WithExperimentKeys(experimentKeys []string) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.ExperimentKeys == nil {
			options.ExperimentKeys = make(map[string]bool, 1)
		}
		for _, experimentKey := range experimentKeys {
			options.ExperimentKeys[experimentKey] = true
		}
		return nil
	}
}

// WithLayerKeyList sets the layer key list,
// and uses filtering to only focus on the hits of experiments under these layers.
// For specific layer keys, you can view the layer key where the experiment is located on the web interface.
func WithLayerKeyList(layerKeyList []string) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.LayerKeys == nil {
			options.LayerKeys = make(map[string]bool, len(layerKeyList))
		}
		for i := range layerKeyList {
			options.LayerKeys[layerKeyList[i]] = true
		}
		return nil
	}
}

// WithSceneID is the same as WithSceneIDList. This is convenient for use with only one scene.
func WithSceneID(sceneID int64) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.SceneIDs == nil {
			options.SceneIDs = make(map[int64]bool, 1)
		}
		options.SceneIDs[sceneID] = true
		return nil
	}
}

// WithLayerKey is the same as WithLayerKeyList, which is convenient for use with only one layer.
// If there is only one layer, it is recommended to use GetExperiment API instead of GetExperiments
func WithLayerKey(layerKey string) ExperimentOption {
	return func(options *experiment.Options) error {
		if options.LayerKeys == nil {
			options.LayerKeys = make(map[string]bool, 1)
		}
		options.LayerKeys[layerKey] = true
		return nil
	}
}

// WithAutomatic sets whether TAB automatically records exposure
func WithAutomatic(isAutomatic bool) ExperimentOption {
	return func(options *experiment.Options) error {
		options.IsExposureLoggingAutomatic = isAutomatic
		return nil
	}
}

// WithIsPreparedDMPTag sets whether to preprocess DMP tags
// If there is no dmp tag configured, or there is only one, preprocessing will not be enabled.
// Closed by default, developers can use this option to manage manually
func WithIsPreparedDMPTag(isPreparedDMPTag bool) ExperimentOption {
	return func(options *experiment.Options) error {
		options.IsPreparedDMPTag = isPreparedDMPTag
		return nil
	}
}

// WithIsDisableDMP sets whether to turn off DMP to prevent rpc operations.
// If the DMP label is clearly not needed, it can be turned off, such as local testing, etc.
func WithIsDisableDMP(isDisableDMP bool) ExperimentOption {
	return func(options *experiment.Options) error {
		options.IsDisableDMP = isDisableDMP
		return nil
	}
}

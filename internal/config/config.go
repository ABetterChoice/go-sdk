// Package config Remote configuration acquisition related implementation,
// the underlying logic part depends on the experiment package
package config

import (
	"context"

	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/experiment"
	"github.com/abetterchoice/go-sdk/plugin/log"
	"github.com/abetterchoice/hashutil"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/pkg/errors"
)

type executor struct{}

var (
	// Executor Remote configuration of logic actuators
	Executor = &executor{}
)

// Value Remote Configuration Values
type Value struct {
	Data           []byte
	IsOverrideList bool
	IsDefault      bool
	IsHoldout      bool
	Experiment     *experiment.Experiment            // Configuration Binding Experiment
	RemoteConfig   *protoc_cache_server.RemoteConfig // Remote configuration details
	UnitIDType     protoc_cache_server.UnitIDType    // ID Account System
}

// GetRemoteConfig Get the specific implementation of remote configuration
func (e *executor) GetRemoteConfig(ctx context.Context, projectID string, key string,
	options *experiment.Options) (*Value,
	error) {
	application := cache.GetApplication(projectID)
	if application == nil {
		return nil, errors.Errorf("projectID [%s] not found", projectID)
	}
	options.Application = application
	remoteConfig, ok := application.TabConfig.ConfigData.RemoteConfigIndex[key]
	if !ok || remoteConfig == nil {
		return nil, errors.Errorf("remoteConfig[%s] not found", key)
	}
	return e.getRemoteConfigValue(ctx, remoteConfig, options)
}

func (e *executor) getRemoteConfigValue(ctx context.Context, config *protoc_cache_server.RemoteConfig,
	options *experiment.Options) (*Value, error) {
	data, unitIDType, ok := e.processOverrideList(config, options)
	if ok {
		return &Value{Data: data, IsOverrideList: true, RemoteConfig: config,
			UnitIDType: unitIDType}, nil
	}
	holdoutExp, err := e.checkCaughtByHoldout(ctx, config.HoldoutLayerKeys, options)
	if err != nil {
		return nil, errors.Wrap(err, "checkCaughtByHoldout")
	}
	if holdoutExp != nil {
		value, _ := holdoutExp.Params[config.Key]
		return &Value{
			Data:           []byte(value),
			IsOverrideList: false,
			IsDefault:      false,
			IsHoldout:      true,
			Experiment:     holdoutExp,
			RemoteConfig:   config,
			UnitIDType:     holdoutExp.UnitIdType,
		}, nil
	}
	unitIDType = protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT
	for _, condition := range config.ConditionList {
		unitIDType = condition.UnitIdType
		value, hit, err := e.processCondition(ctx, condition, options)
		if err != nil {
			return nil, err
		}
		if hit {
			value.RemoteConfig = config
			value.UnitIDType = condition.UnitIdType
			return value, err
		}
	}
	return &Value{Data: config.DefaultValue, IsDefault: true, RemoteConfig: config, UnitIDType: unitIDType}, nil
}

func (e *executor) checkCaughtByHoldout(ctx context.Context, holdoutLayerKeys []string,
	options *experiment.Options) (*experiment.Experiment, error) {
	holdoutData := options.Application.TabConfig.ExperimentData.HoldoutData
	if holdoutData == nil {
		return nil, nil
	}
	for _, holdoutLayerKey := range holdoutLayerKeys {
		if holdoutExp, ok := options.HoldoutLayerResult[holdoutLayerKey]; !ok {
			holdoutLayer, isExist := holdoutData.HoldoutLayerIndex[holdoutLayerKey]
			if !isExist {
				return nil, errors.Errorf("invalid holdout layerKey=%s", holdoutLayerKey)
			}
			options.HoldoutLayerResult[holdoutLayerKey] = nil
			newHoldoutExp, err := experiment.Executor.GetLayerExperiment(ctx, holdoutLayer, options)
			if err != nil {
				return nil, err
			}
			options.HoldoutLayerResult[holdoutLayerKey] = newHoldoutExp
			if newHoldoutExp != nil {
				return newHoldoutExp, nil
			}
		} else if holdoutExp != nil {
			return holdoutExp, nil
		}
	}
	return nil, nil
}

func (e *executor) processOverrideList(config *protoc_cache_server.RemoteConfig, options *experiment.Options) ([]byte,
	protoc_cache_server.UnitIDType, bool) {
	data, ok := config.OverrideList[options.UnitID]
	if ok {
		return data, protoc_cache_server.UnitIDType_UNIT_ID_TYPE_DEFAULT, ok
	}
	data, ok = config.OverrideList[options.NewUnitID]
	return data, protoc_cache_server.UnitIDType_UNIT_ID_TYPE_NEW_ID, ok
}

func (e *executor) processCondition(ctx context.Context, condition *protoc_cache_server.Condition,
	options *experiment.Options) (*Value, bool, error) {
	bucketNum := hashutil.GetBucketNum(condition.HashMethod, e.getHashSource(condition.UnitIdType, options),
		condition.HashSeed, condition.BucketSize)
	if !e.isHitConditionBucketInfo(bucketNum, condition.BucketInfo) {
		return nil, false, nil
	}
	if condition.IssueInfo == nil {
		return nil, false, nil
	}
	switch condition.IssueInfo.IssueType {
	case protoc_cache_server.IssueType_ISSUE_TYPE_PERCENTAGE:
		value, err := processConditionExperiment(ctx, condition, options)
		return value, true, err
	case protoc_cache_server.IssueType_ISSUE_TYPE_TAG, protoc_cache_server.IssueType_ISSUE_TYPE_CITY_TAG:
		hit, err := experiment.IsHitTag(ctx, condition.IssueInfo.TagListGroup, options)
		if err != nil {
			return nil, false, errors.Wrapf(err, "isHitTag")
		}
		if !hit {
			return nil, false, nil
		}
		value, err := processConditionExperiment(ctx, condition, options)
		return value, true, err
	}
	return nil, false, nil
}

func processConditionExperiment(ctx context.Context, condition *protoc_cache_server.Condition,
	options *experiment.Options) (*Value, error) {
	value := &Value{}
	value.Data = condition.Value
	if condition.ExperimentKey == "" {
		return value, nil
	}
	if options.ExperimentKeys == nil {
		options.ExperimentKeys = make(map[string]bool, 1)
	}
	options.ExperimentKeys[condition.ExperimentKey] = true
	experimentList, err := experiment.Executor.GetExperiments(ctx, options.Application.ProjectID, options)
	if err != nil {
		return nil, errors.Wrapf(err, "getExperiments with experimentKey[%s]", condition.ExperimentKey)
	}
	for _, e := range experimentList {
		if e.Group.ExperimentKey != condition.ExperimentKey {
			continue
		}
		result, ok := e.Group.Params[condition.ConfigKey]
		if !ok {
			return value, nil
		}
		value.Experiment = e
		value.Data = []byte(result)
		return value, nil
	}
	return value, nil
}

func (e *executor) isHitConditionBucketInfo(bucketNum int64, bucketInfo *protoc_cache_server.BucketInfo) bool {
	switch bucketInfo.BucketType {
	case protoc_cache_server.BucketType_BUCKET_TYPE_RANGE:
		return bucketInfo.TrafficRange.Left <= bucketNum && bucketNum <= bucketInfo.TrafficRange.Right
	default:
		log.Errorf("invalid bucketType=%v", bucketInfo.BucketType)
		return false
	}
}

func (e *executor) getHashSource(unitType protoc_cache_server.UnitIDType, options *experiment.Options) string {
	if unitType == protoc_cache_server.UnitIDType_UNIT_ID_TYPE_NEW_ID {
		return options.NewDecisionID
	}
	return options.DecisionID
}

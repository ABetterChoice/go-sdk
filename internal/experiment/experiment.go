// Package experiment abtest Experimental diversion related implementation
package experiment

import (
	"context"
	"strconv"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/plugin/log"
	"github.com/abetterchoice/hashutil"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_dmp_proxy_server"
	"github.com/abetterchoice/tagutil"
	"github.com/pkg/errors"
)

type executor struct{}

var (
	// Executor abtest Experimental Shunt Logic Actuator
	Executor = &executor{}
)

// Experiment Experimental information, encapsulates the group in the protocol,
// and adds status attributes during the diversion process
// IsOverrideList Whether it is the experimental group hit by the whitelist
type Experiment struct {
	*protoccacheserver.Group // The hit experiment group, Group must be not empty
	IsOverrideList           bool
	HoldoutData              map[string]*Experiment
}

// VariantKey2LayerKey Get the layer where the parameter key is located according to the parameter key
func (e *executor) VariantKey2LayerKey(projectID, variantKey string) ([]string, error) {
	application := cache.GetApplication(projectID)
	if application == nil {
		return nil, errors.Errorf("projectID [%s] not found", projectID)
	}
	return application.VariantKeyLayerMap[variantKey], nil
}

// GetVariantValue Get the parameter value of the layer default parameter
func (e *executor) GetVariantValue(projectID, layerKey, variantKey string) ([]byte, error) {
	application := cache.GetApplication(projectID)
	if application == nil {
		return nil, errors.Errorf("projectID [%s] not found", projectID)
	}
	layer, ok := application.LayerIndex[layerKey]
	if !ok {
		return nil, errors.Errorf("invalid layerKey:%v", layerKey)
	}
	if layer.Metadata == nil || layer.Metadata.DefaultGroup == nil {
		return nil, errors.Errorf("invalid layer")
	}
	defaultGroup := layer.Metadata.DefaultGroup
	data, ok := defaultGroup.Params[variantKey]
	if ok {
		return []byte(data), nil
	}
	return nil, errors.Errorf("default parameter of layer[%s] does not exist parameter %s", layerKey, variantKey)
}

// GetExperiments Get the set of experiment information that the user hits under the conditions specified by options
func (e *executor) GetExperiments(ctx context.Context, projectID string, options *Options) (map[string]*Experiment,
	error) {
	application := cache.GetApplication(projectID)
	if application == nil {
		return nil, errors.Errorf("projectID [%s] not found", projectID)
	}
	err := e.fillOptions(ctx, application, options)
	if err != nil {
		return nil, errors.Wrap(err, "fillOptions")
	}
	layers, flag, err := e.layersCanBeHit(ctx, application, options)
	if err != nil {
		return nil, errors.Wrap(err, "layersCanBeHit")
	}
	if flag {
		return e.getMultiLayerExperiments(ctx, layers, options)
	}
	result, err := e.getDomainExperiments(ctx, application.TabConfig.ExperimentData.GlobalDomain, options)
	if err != nil {
		return nil, errors.Wrapf(err, "getDomainExperiments")
	}
	if len(options.OverrideList) > 0 {
		e.setOverrideGroup(ctx, application, result, options)
	}
	return result, nil
}

func (e *executor) setOverrideGroup(ctx context.Context, application *cache.Application, result map[string]*Experiment,
	options *Options) {
	for layerKey, groupID := range options.OverrideList {
		layer, ok := application.LayerIndex[layerKey]
		if !ok {
			log.Error("layer[%s] not found", layerKey)
			continue
		}
		if !e.isLayerFilterPass(ctx, layer, options) {
			continue
		}
		if _, ok := result[layerKey]; !ok {
			group, exist := layer.GroupIndex[groupID]
			if exist && group != nil {
				result[layerKey] = &Experiment{
					Group:          group,
					IsOverrideList: true,
				}
			}
		}
	}
}

func (e *executor) layersCanBeHit(ctx context.Context, application *cache.Application,
	options *Options) ([]*protoccacheserver.Layer, bool, error) {
	if len(options.LayerKeys) == 0 {
		return nil, false, nil
	}
	var result = make([]*protoccacheserver.Layer, 0, len(options.LayerKeys))
	for layerKey := range options.LayerKeys {
		layer, err := e.checkLayerCanBeHit(ctx, application, layerKey, options)
		if err != nil {
			return nil, false, errors.Wrap(err, "checkLayerCanBeHit")
		}
		if layer != nil {
			result = append(result, layer)
		}
	}
	return result, true, nil
}

func (e *executor) checkLayerCanBeHit(ctx context.Context, application *cache.Application, layerKey string,
	options *Options) (*protoccacheserver.Layer,
	error) {
	layer, ok := application.FullFlowLayerIndex[layerKey]
	if ok {
		return layer, nil
	}
	layer, ok = application.LayerIndex[layerKey]
	if !ok || layer == nil {
		return nil, errors.Errorf("invalid layerKey=%s", layerKey)
	}
	holdoutExp, err := e.checkCaughtByHoldout(ctx, application, layer, options)
	if err != nil {
		return nil, errors.Wrap(err, "checkCaughtByHoldout")
	}
	if holdoutExp != nil {
		return application.LayerIndex[layerKey], nil
	}
	flag, err := isHitLayer(application, layerKey, options)
	if err != nil {
		return nil, errors.Wrap(err, "isHitLayer")
	}
	if flag {
		return application.LayerIndex[layerKey], nil
	}
	return nil, nil
}

func (e *executor) checkCaughtByHoldout(ctx context.Context, application *cache.Application,
	layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	holdoutData := application.TabConfig.ExperimentData.HoldoutData
	if holdoutData == nil {
		return nil, nil
	}
	for _, holdoutLayerKey := range layer.Metadata.HoldoutLayerKeys {
		if holdoutExp, ok := options.HoldoutLayerResult[holdoutLayerKey]; !ok {
			holdoutLayer, isExist := holdoutData.HoldoutLayerIndex[holdoutLayerKey]
			if !isExist {
				return nil, errors.Errorf("invalid holdout layerKey=%s", holdoutLayerKey)
			}
			options.HoldoutLayerResult[holdoutLayerKey] = nil
			experiment, err := e.GetLayerExperiment(ctx, holdoutLayer, options)
			if err != nil {
				return nil, err
			}
			options.HoldoutLayerResult[holdoutLayerKey] = experiment
			if experiment != nil && !experiment.IsDefault && experiment.IsControl {
				return experiment, nil
			}
		} else if holdoutExp != nil && !holdoutExp.IsDefault && holdoutExp.IsControl {
			return holdoutExp, nil
		}
	}
	return nil, nil
}

func isHitLayer(application *cache.Application, layerKey string, options *Options) (bool, error) {
	domainMetadataList, ok := application.LayerDomainMetadataListIndex[layerKey]
	if !ok {
		log.Warnf("domainMetadata not found[invalid layerKey=%v]", layerKey)
		return false, nil
	}
	bucketNum := int64(0)
	for i, domainMetadata := range domainMetadataList {
		if i != 0 && !isHitTraffic(bucketNum, domainMetadata) {
			return false, nil
		}
		bucketNum = hashutil.GetBucketNum(domainMetadata.HashMethod, getHashSource(domainMetadata.UnitIdType, options),
			domainMetadata.HashSeed, domainMetadata.BucketSize)
	}
	return true, nil
}

func (e *executor) fillOptions(ctx context.Context, application *cache.Application, options *Options) error {
	e.setOverrideList(application, options)
	options.Application = application
	if !options.IsDisableDMP && options.IsPreparedDMPTag {
		err := e.preparedDMP(ctx, application, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *executor) preparedDMP(ctx context.Context, application *cache.Application, options *Options) error {
	if len(application.DMPTagInfo) == 0 {
		return nil
	}
	for unitIDType, platformDMPTagIndex := range application.DMPTagInfo {
		for platformCode, tagSet := range platformDMPTagIndex {
			if len(tagSet) <= 1 {
				continue
			}
			req := &protoc_dmp_proxy_server.BatchGetTagValueReq{
				ProjectId:       application.ProjectID,
				UnitId:          getUnitID(unitIDType, options),
				UnitType:        int64(unitIDType),
				SdkVersion:      env.SDKVersion,
				DmpPlatformCode: protoc_dmp_proxy_server.DMPPlatform(platformCode),
				TagList:         convertMap2Array(tagSet),
			}
			resp, err := client.DC.BatchGetTagValue(ctx, req)
			if err != nil {
				log.Errorf("[req=%+v]BatchGetTagValue fail:%v", req, err)
				continue
			}
			if resp.RetCode != protoc_dmp_proxy_server.RetCode_RET_CODE_SUCCESS {
				log.Errorf("[req=%+v]BatchGetTagValue invalid code=%v, message=%v", req, resp.RetCode, resp.Message)
				continue
			}
			if options.DMPTagValueResult == nil {
				options.DMPTagValueResult = make(map[string]string)
			}
			for key, value := range resp.TagResult {
				options.DMPTagValueResult[dmpTagResultKeyFormat(protoccacheserver.UnitIDType(req.UnitType),
					int64(req.DmpPlatformCode), key, options)] = value
			}
		}
	}
	return nil
}

func convertMap2Array(source map[string]interface{}) []string {
	var result = make([]string, len(source))
	var i int
	for key := range source {
		result[i] = key
		i++
	}
	return result
}

func (e *executor) setOverrideList(application *cache.Application, options *Options) {
	if application.TabConfig.ExperimentData.OverrideList != nil {
		overrideList := application.TabConfig.ExperimentData.OverrideList[options.UnitID]
		if overrideList != nil && len(overrideList.LayerToGroupId) > 0 {
			options.OverrideList = overrideList.LayerToGroupId
		}
		if len(options.NewUnitID) > 0 {
			newOverrideList := application.TabConfig.ExperimentData.OverrideList[options.NewUnitID]
			if newOverrideList != nil && len(newOverrideList.LayerToGroupId) >= 0 {
				var result = make(map[string]int64, len(options.OverrideList)+len(newOverrideList.LayerToGroupId))
				for key, value := range options.OverrideList {
					result[key] = value
				}
				for key, value := range newOverrideList.LayerToGroupId {
					result[key] = value
				}
				options.OverrideList = result
			}
		}
	}
}

func (e *executor) getDomainExperiments(ctx context.Context, domain *protoccacheserver.Domain,
	options *Options) (map[string]*Experiment, error) {
	bucketNum := hashutil.GetBucketNum(domain.Metadata.HashMethod,
		getHashSource(domain.Metadata.UnitIdType, options),
		domain.Metadata.HashSeed, domain.Metadata.BucketSize)
	for _, holdoutDomain := range domain.HoldoutDomainList {
		if isHitTraffic(bucketNum, holdoutDomain.Metadata) {
			return e.getHoldoutDomainExperiments(ctx, holdoutDomain, options)
		}
	}
	var result = make(map[string]*Experiment)
	for _, multiLayerDomain := range domain.MultiLayerDomainList {
		if isHitTraffic(bucketNum, multiLayerDomain.Metadata) {
			multiLayerDomainResult, err := e.getMultiLayerDomainExperiments(ctx, multiLayerDomain, options)
			if err != nil {
				return nil, errors.Wrap(err, "getMultiLayerDomainExperiments")
			}
			result = multiLayerDomainResult
			break
		}
	}
	for _, subdomain := range domain.DomainList {
		if !isHitTraffic(bucketNum, subdomain.Metadata) {
			continue
		}
		subdomainResult, err := e.getDomainExperiments(ctx, subdomain, options)
		if err != nil {
			return nil, errors.Wrapf(err, "getDomainExperiment[%s]", subdomain.Metadata.Key)
		}
		for layerKey, group := range subdomainResult {
			result[layerKey] = group
		}
	}
	return result, nil
}

func (e *executor) getHoldoutDomainExperiments(ctx context.Context, holdoutDomain *protoccacheserver.HoldoutDomain,
	options *Options) (map[string]*Experiment, error) {
	return e.getMultiLayerExperiments(ctx, holdoutDomain.LayerList, options)
}

func (e *executor) getMultiLayerDomainExperiments(ctx context.Context,
	multiLayerDomain *protoccacheserver.MultiLayerDomain, options *Options) (map[string]*Experiment, error) {
	return e.getMultiLayerExperiments(ctx, multiLayerDomain.LayerList, options)
}

func (e *executor) getMultiLayerExperiments(ctx context.Context, layerList []*protoccacheserver.Layer,
	options *Options) (map[string]*Experiment, error) {
	var result = make(map[string]*Experiment)
	for _, layer := range layerList {
		g, err := e.getLayerExperimentWithDefault(ctx, layer, options)
		if err != nil {
			return nil, errors.Wrap(err, "GetLayerExperiment")
		}
		if g == nil {
			continue
		}
		result[layer.Metadata.Key] = g
		// set holdout data
		e.setHoldout2Experiment(g, layer, options)
	}
	return result, nil
}

func (e *executor) setHoldout2Experiment(experiment *Experiment, layer *protoccacheserver.Layer, options *Options) {
	if len(layer.Metadata.HoldoutLayerKeys) == 0 {
		return
	}
	var holdoutLayerKeys = make(map[string]interface{})
	for _, holdoutKey := range layer.Metadata.HoldoutLayerKeys {
		var hs = make(map[string]interface{})
		e.getHoldoutLayerKey(holdoutKey, options, hs)
		for key := range hs {
			holdoutLayerKeys[key] = nil
		}
	}
	for key := range holdoutLayerKeys {
		holdoutExperiment, ok := options.HoldoutLayerResult[key]
		if !ok {
			continue
		}
		if holdoutExperiment == nil {
			continue
		}
		if experiment.HoldoutData == nil {
			experiment.HoldoutData = make(map[string]*Experiment)
		}
		experiment.HoldoutData[key] = holdoutExperiment
	}
}

func (e *executor) getHoldoutLayerKey(holdoutLayerKey string, options *Options, result map[string]interface{}) {
	if options.Application == nil || options.Application.TabConfig == nil ||
		options.Application.TabConfig.ExperimentData == nil ||
		options.Application.TabConfig.ExperimentData.HoldoutData == nil {
		return
	}
	holdoutData := options.Application.TabConfig.ExperimentData.HoldoutData
	holdoutLayer, _ := holdoutData.HoldoutLayerIndex[holdoutLayerKey]
	if holdoutLayer == nil || holdoutLayer.Metadata == nil {
		return
	}
	result[holdoutLayerKey] = nil
	for _, key := range holdoutLayer.Metadata.HoldoutLayerKeys {
		if _, ok := result[key]; !ok {
			e.getHoldoutLayerKey(key, options, result)
		}
	}
}

func (e *executor) isLayerFilterPass(ctx context.Context, layer *protoccacheserver.Layer, options *Options) bool {
	return layerKeyFilter(ctx, layer, options) &&
		sceneIDListFilter(ctx, layer, options) &&
		experimentKeyFilter(ctx, layer, options)
}

func (e *executor) getLayerExperimentWithDefault(ctx context.Context, layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	if len(layer.GroupIndex) == 0 {
		return nil, nil
	}
	if !e.isLayerFilterPass(ctx, layer, options) {
		return nil, nil
	}
	experiment, err := e.GetLayerExperiment(ctx, layer, options)
	if err != nil {
		return nil, err
	}
	if experiment != nil {
		return experiment, nil
	}
	if layer.Metadata.DefaultGroup != nil {
		return &Experiment{Group: layer.Metadata.DefaultGroup}, nil
	}
	if len(layer.GroupIndex) == 0 {
		return nil, nil
	}
	return &Experiment{
		Group: &protoccacheserver.Group{
			Id:        e.defaultSystemGlobalGroupID(options),
			GroupKey:  env.DefaultGlobalGroupKey,
			Params:    nil,
			IsDefault: true,
			LayerKey:  layer.Metadata.Key,
		},
		IsOverrideList: false,
	}, nil
}

func (e *executor) defaultSystemGlobalGroupID(options *Options) int64 {
	if options.Application.TabConfig.ExperimentData.DefaultGroupId != 0 {
		return options.Application.TabConfig.ExperimentData.DefaultGroupId
	}
	return env.DefaultGlobalGroupID
}

// GetLayerExperiment Get the set of experiment information that the user hits under the conditions specified by options
func (e *executor) GetLayerExperiment(ctx context.Context, layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	overrideGroup := e.getLayerOverrideExperiment(layer, options)
	if overrideGroup != nil {
		return overrideGroup, nil
	}
	holdoutExp, err := e.checkCaughtByHoldout(ctx, options.Application, layer, options)
	if err != nil {
		return nil, errors.Wrap(err, "checkCaughtByHoldout")
	}
	if holdoutExp != nil && !holdoutExp.IsDefault && holdoutExp.IsControl {
		return holdoutExp, nil
	}
	switch layer.Metadata.HashType {
	case protoccacheserver.HashType_HASH_TYPE_DOUBLE:
		return e.getDoubleHashLayerExperiment(ctx, layer, options)
	case protoccacheserver.HashType_HASH_TYPE_SINGLE:
		return e.getSingleHashLayerExperiment(ctx, layer, options)
	default:
		return nil, nil
	}
}

func (e *executor) getLayerOverrideExperiment(layer *protoccacheserver.Layer, options *Options) *Experiment {
	groupID, ok := options.OverrideList[layer.Metadata.Key]
	if !ok {
		return nil
	}
	group, ok := layer.GroupIndex[groupID]
	if ok {
		return &Experiment{
			Group:          group,
			IsOverrideList: true,
		}
	}
	return nil
}

func (e *executor) getSingleHashLayerExperiment(ctx context.Context, layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	bucketNum := hashutil.GetBucketNum(layer.Metadata.HashMethod,
		getHashSource(layer.Metadata.UnitIdType, options),
		layer.Metadata.HashSeed, layer.Metadata.BucketSize)
	for _, group := range layer.GroupIndex {
		if group.IsDefault {
			continue
		}
		if !e.isHitGroupBucketInfo(group, bucketNum, options) {
			continue
		}
		switch group.IssueInfo.IssueType {
		case protoccacheserver.IssueType_ISSUE_TYPE_TAG:
			tagFlag, err := IsHitTag(ctx, group.IssueInfo.TagListGroup, options)
			if err != nil {
				return nil, errors.Wrap(err, "isHitTag")
			}
			if tagFlag {
				return &Experiment{Group: group}, nil
			}
		case protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE:
			return &Experiment{Group: group}, nil
		}
	}
	return nil, nil
}

// IsHitTag Whether the tag is hit
func IsHitTag(ctx context.Context, tagListGroup []*protoccacheserver.TagList, options *Options) (bool, error) {
	if len(tagListGroup) == 0 {
		return true, nil
	}
	for _, tagList := range tagListGroup {
		isHit := true
		for _, tag := range tagList.TagList {
			if tag.TagType == protoccacheserver.TagType_TAG_TYPE_DMP {
				if options.IsDisableDMP { // 禁用 结果都为 false
					return false, nil
				}
				dmpFlagStr, err := getTagValue(ctx, tag, options)
				if err != nil {
					log.Errorf("[tag=%v]getTagValue fail:%v", err)
					return false, nil
				}
				dmpFlag, err := strconv.ParseBool(dmpFlagStr)
				if err != nil {
					log.Errorf("dmpFlag=%v needs to be convertible to bool type:%v", dmpFlagStr, err)
					return false, nil
				}
				if dmpFlag && tag.Operator == protoccacheserver.Operator_OPERATOR_FALSE ||
					!dmpFlag && tag.Operator == protoccacheserver.Operator_OPERATOR_TRUE {
					isHit = false
					break
				}
				continue
			}
			if tag.TagOrigin == protoccacheserver.TagOrigin_TAG_ORIGIN_DMP {
				if options.IsDisableDMP { // 禁用 结果都为 false
					return false, nil
				}
				value, err := getTagValue(ctx, tag, options)
				if err != nil {
					log.Errorf("[tag=%v]getTagValue fail:%v", tag.Key, err)
					return false, nil
				}
				if v, ok := options.AttributeTag[tag.Key]; ok {
					log.Warnf("replace tagValue[key=%s, value=%s]", tag.Key, v)
				}
				if options.AttributeTag == nil {
					options.AttributeTag = make(map[string][]string)
				}
				options.AttributeTag[tag.Key] = []string{value}
			}
			if !tagutil.IsHit(tag.TagType, tag.Operator, options.AttributeTag[tag.Key], tag.Value) {
				isHit = false
				break
			}
		}
		if isHit {
			return true, nil
		}
	}
	return false, nil
}

func getTagValue(ctx context.Context, tag *protoccacheserver.Tag, options *Options) (string, error) {
	key := dmpTagResultKeyFormat(tag.UnitIdType, tag.DmpPlatform, tag.Key, options)
	value, ok := options.DMPTagValueResult[key]
	if ok {
		return value, nil
	}
	resp, err := client.DC.BatchGetTagValue(ctx, &protoc_dmp_proxy_server.BatchGetTagValueReq{
		ProjectId:       options.Application.ProjectID,
		UnitId:          options.UnitID,
		UnitType:        0,
		SdkVersion:      env.SDKVersion,
		DmpPlatformCode: protoc_dmp_proxy_server.DMPPlatform(tag.DmpPlatform),
		TagList:         []string{tag.Key},
	})
	if err != nil {
		return "", err
	}
	if resp.RetCode != protoc_dmp_proxy_server.RetCode_RET_CODE_SUCCESS {
		return "", errors.Errorf("invalid result, code=%v, message=%v", resp.RetCode, resp.Message)
	}
	value, ok = resp.TagResult[tag.Key]
	if !ok {
		return "", errors.Errorf("value is empty")
	}
	if options.DMPTagValueResult == nil {
		options.DMPTagValueResult = make(map[string]string)
	}
	options.DMPTagValueResult[key] = value
	return value, nil
}

func dmpTagResultKeyFormat(unitIDType protoccacheserver.UnitIDType, dmpPlatformCode int64,
	dmpTagKey string, options *Options) string {
	return getUnitID(unitIDType, options) + "-" + strconv.FormatInt(dmpPlatformCode, 10) + "-" + dmpTagKey
}

func (e *executor) isHitGroupBucketInfo(group *protoccacheserver.Group, bucketNum int64, options *Options) bool {
	bucketInfo, ok := options.Application.GroupIDBucketInfoIndex[group.Id]
	if !ok {
		return false
	}
	switch bucketInfo.BucketType {
	case protoccacheserver.BucketType_BUCKET_TYPE_RANGE:
		return bucketInfo.TrafficRange.Left <= bucketNum && bucketNum <= bucketInfo.TrafficRange.Right
	case protoccacheserver.BucketType_BUCKET_TYPE_BITMAP:
		bitmap, ok := options.Application.GroupIDRoaringBitmapIndex[group.Id]
		if !ok {
			log.Warnf("groupID[%d] bitmap not found", group.Id)
			return false
		}
		return bitmap.ContainsInt(int(bucketNum))
	default:
		return false
	}
}

func (e *executor) isHitExperimentBucketInfo(experiment *protoccacheserver.Experiment, bucketNum int64,
	options *Options) bool {
	if experiment.Id == 0 {
		return false
	}
	bucketInfo, ok := options.Application.ExperimentIDBucketInfoIndex[experiment.Id]
	if !ok {
		return false
	}
	switch bucketInfo.BucketType {
	case protoccacheserver.BucketType_BUCKET_TYPE_RANGE:
		return bucketInfo.TrafficRange.Left <= bucketNum && bucketNum <= bucketInfo.TrafficRange.Right
	case protoccacheserver.BucketType_BUCKET_TYPE_BITMAP:
		bitmap, ok := options.Application.ExperimentIDRoaringBitmapIndex[experiment.Id]
		if !ok {
			log.Warnf("experimentID[%d] bitmap not found", experiment.Id)
			return false
		}
		return bitmap.ContainsInt(int(bucketNum))
	default:
		return false
	}
}

func (e *executor) getDoubleHashLayerExperiment(ctx context.Context, layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	bucketNum := hashutil.GetBucketNum(layer.Metadata.HashMethod,
		getHashSource(layer.Metadata.UnitIdType, options),
		layer.Metadata.HashSeed, layer.Metadata.BucketSize)
	for _, experiment := range layer.ExperimentIndex {
		if !e.isHitExperimentBucketInfo(experiment, bucketNum, options) {
			continue
		}
		return e.getExperimentGroup(ctx, experiment, layer, options)
	}
	return nil, nil
}

func (e *executor) getExperimentGroup(ctx context.Context, experiment *protoccacheserver.Experiment,
	layer *protoccacheserver.Layer, options *Options) (*Experiment, error) {
	expBucketNum := hashutil.GetBucketNum(experiment.HashMethod,
		getHashSource(layer.Metadata.UnitIdType, options), experiment.HashSeed, experiment.BucketSize)
	switch experiment.IssueType {
	case protoccacheserver.IssueType_ISSUE_TYPE_PERCENTAGE:
		return e.getPercentageExperimentGroup(experiment, expBucketNum, layer, options)
	case protoccacheserver.IssueType_ISSUE_TYPE_TAG:
		return e.getTagExperimentGroup(ctx, experiment, expBucketNum, layer, options)
	case protoccacheserver.IssueType_ISSUE_TYPE_CITY_TAG:
		return e.getCityTagExperimentGroup(ctx, experiment, expBucketNum, layer, options)
	default:
		return nil, nil
	}
}

func (e *executor) getPercentageExperimentGroup(experiment *protoccacheserver.Experiment, expBucketNum int64,
	layer *protoccacheserver.Layer,
	options *Options) (*Experiment, error) {
	for groupID := range experiment.GroupIdIndex {
		group, ok := layer.GroupIndex[groupID]
		if !ok || group == nil {
			return nil, errors.Errorf("invalid groupID[%v]", groupID)
		}
		if !e.isHitGroupBucketInfo(group, expBucketNum, options) {
			continue
		}
		return &Experiment{Group: group}, nil
	}
	return nil, nil
}

func (e *executor) getTagExperimentGroup(ctx context.Context, experiment *protoccacheserver.Experiment,
	expBucketNum int64, layer *protoccacheserver.Layer, options *Options) (*Experiment, error) {
	var needToTestHitTag = true
	for groupID := range experiment.GroupIdIndex {
		group, ok := layer.GroupIndex[groupID]
		if !ok || group == nil {
			return nil, errors.Errorf("invalid groupID[%v]", groupID)
		}
		if needToTestHitTag {
			tagFlag, err := IsHitTag(ctx, group.IssueInfo.TagListGroup, options)
			if err != nil {
				return nil, errors.Wrap(err, "isHitTag")
			}
			if !tagFlag {
				return nil, nil
			}
			needToTestHitTag = false
		}
		if !e.isHitGroupBucketInfo(group, expBucketNum, options) {
			continue
		}
		return &Experiment{Group: group}, nil
	}
	return nil, nil
}

func (e *executor) getCityTagExperimentGroup(ctx context.Context, experiment *protoccacheserver.Experiment,
	expBucketNum int64, layer *protoccacheserver.Layer, options *Options) (*Experiment, error) {
	for groupID := range experiment.GroupIdIndex {
		group, ok := layer.GroupIndex[groupID]
		if !ok || group == nil {
			return nil, errors.Errorf("invalid groupID=%v", groupID)
		}
		tagFlag, err := IsHitTag(ctx, group.IssueInfo.TagListGroup, options)
		if err != nil {
			return nil, errors.Wrap(err, "isHitTag")
		}
		if !tagFlag {
			continue
		}
		if !e.isHitGroupBucketInfo(group, expBucketNum, options) {
			return nil, nil
		}
		return &Experiment{Group: group}, nil
	}
	return nil, nil
}

func getHashSource(unitType protoccacheserver.UnitIDType, options *Options) string {
	if unitType == protoccacheserver.UnitIDType_UNIT_ID_TYPE_NEW_ID {
		return options.NewDecisionID
	}
	return options.DecisionID
}

func getUnitID(unitType protoccacheserver.UnitIDType, options *Options) string {
	if unitType == protoccacheserver.UnitIDType_UNIT_ID_TYPE_NEW_ID && options.NewUnitID != "" {
		return options.NewUnitID
	}
	return options.UnitID
}

func isHitTraffic(bucketNum int64, metadata *protoccacheserver.DomainMetadata) bool {
	if metadata == nil {
		return false
	}
	for _, traffic := range metadata.TrafficRangeList {
		if traffic.Left <= bucketNum && bucketNum <= traffic.Right {
			return true
		}
	}
	return false
}

// Package cache 本地缓存具体实现
package cache

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/plugin/log"
	metrics2 "github.com/abetterchoice/go-sdk/plugin/metrics"
	protoctabcacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Application Locally cached entities
type Application struct {
	// Project Code
	ProjectID string
	Version   string
	// Experiment, configuration, switch information
	TabConfig *protoctabcacheserver.TabConfig
	// Experimental barrel information
	ExperimentIDBucketInfoIndex map[int64]*protoctabcacheserver.BucketInfo
	// roaring bitmap index, bucketInfo type is bitmap, roaring bitmap object will be pre-built
	ExperimentIDRoaringBitmapIndex map[int64]*roaring.Bitmap
	// Experimental group bucket information
	GroupIDBucketInfoIndex map[int64]*protoctabcacheserver.BucketInfo
	// roaring bitmap index, bucketInfo type is bitmap, roaring bitmap object will be pre-built
	GroupIDRoaringBitmapIndex map[int64]*roaring.Bitmap
	// The layer that accounts for 100% of the overall traffic, key is layerKey, value = specific layer information
	FullFlowLayerIndex map[string]*protoctabcacheserver.Layer
	// Experiment layer index, key is layerKey, value is layer
	LayerIndex map[string]*protoctabcacheserver.Layer
	// The metadata information list of each level of the layer, from left to right,
	// corresponds to the layer domain structure from top to bottom
	LayerDomainMetadataListIndex map[string][]*protoctabcacheserver.DomainMetadata
	// Monitor and report component initialization parameters, key is the plugin name,
	// value is the initialization parameter, parse remote cache data
	MetricsPluginInitConfigIndex map[string]*protoctabcacheserver.MetricsInitConfig
	// dmp tag information, from left to right, the keys are unitIDType-dmp platform enumeration ID-dmp tag
	DMPTagInfo map[protoctabcacheserver.UnitIDType]map[int64]map[string]interface{}
	// Mapping of parameters to experimental layers
	VariantKeyLayerMap map[string][]string
	// Whether to preprocess dmp tags
	PreparedDMPTag bool
	// Whether to disable the dmp tag, then the abtest traffic will be completely diverted to the local cache,
	// and there will be no rpc. If disabled, the dmp tag will not be hit by default
	DisableDMPTag bool
	// Current retry count. When TabConfig changes,
	// the bucket information will request the background cache service within the next n times.
	// Avoid local cache not updating when there is a problem with backend consistency.
	// It will be silent after maxRetryTime times to save bandwidth
	retryTime int
}

const (
	// The number of retries to obtain bucket information. If the cached data version is not updated,
	// it will enter the silent period after 10 cumulative retries in each subsequent asynchronous pull.
	// During the retry period, if there is data update, the number of retries will be reset and counted again
	maxRetryTime = 10
	// The local cache refresh interval can be dynamically replaced by the control field in the cache data
	defaultRefreshInterval = 3
)

var localApplicationCache sync.Map

// InitLocalCache Initialize the local cache, and start an independent asynchronous refresh coroutine for
// each projectID, and regularly pull the latest data from the remote background cache service to the local
// Can be initialized multiple times, concurrent and safe
func InitLocalCache(ctx context.Context, projectIDList []string) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, projectID := range projectIDList {
		if _, ok := localApplicationCache.Load(projectID); ok { // If it exists, it will not be refreshed again
			continue
		}
		bc := projectID
		g.Go(func() error {
			_, err := NewAndSetApplication(ctx, bc)
			if err != nil {
				return err
			}
			go continuousFetch(bc)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return err
	}
	return nil
}

// asyncRefreshLocalCache Asynchronously refresh each projectID local cache
func asyncRefreshLocalCache(projectIDList []string) {
	for _, projectID := range projectIDList {
		if _, ok := localApplicationCache.Load(projectID); ok { // If the cache exists, start the refresh coroutine
			go continuousFetch(projectID)
		}
	}
}

// continuousFetch Infinite loop refresh local cache
func continuousFetch(projectID string) {
	for {
		application := GetApplication(projectID)
		if application == nil { // The local cache does not exist, exit the refresh coroutine
			log.Warnf("stop refresh %v", projectID)
			return
		}
		log.Debugf("[projectID=%v] alive", projectID)
		start := time.Now()
		_, err := NewAndSetApplication(context.Background(), projectID)
		latency := time.Since(start)
		if err != nil {
			log.Errorf("[projectID=%v,latency=%s]newApplication fail:%v", projectID, latency.String(), err)
		}
		manualFetchEvent(projectID, latency, err)
		time.Sleep(time.Duration(refreshInterval(projectID)) * time.Second)
	}
}

// manualFetchEvent Log local cache refresh events
func manualFetchEvent(projectID string, latency time.Duration, err error) {
	application := GetApplication(projectID)
	if application == nil {
		return
	}
	metricsConfig := application.TabConfig.ControlData.EventMetricsConfig
	if metricsConfig == nil || !metricsConfig.IsEnable {
		return
	}
	sendDataErr := metrics2.LogMonitorEvent(context.Background(), &metrics2.Metadata{
		MetricsPluginName: metricsConfig.PluginName,
		TableName:         metricsConfig.Metadata.Name,
		TableID:           metricsConfig.Metadata.Id,
		Token:             metricsConfig.Metadata.Token,
		SamplingInterval:  metricsConfig.ErrSamplingInterval,
	}, &protoc_event_server.MonitorEventGroup{Events: []*protoc_event_server.MonitorEvent{
		{
			Time:       time.Now().Unix(),
			Ip:         "",
			ProjectId:  projectID,
			EventName:  "refresh",
			Latency:    float32(latency.Microseconds()), // us
			StatusCode: env.EventStatus(err),
			Message:    env.ErrMsg(err),
			SdkType:    env.SDKType,
			SdkVersion: env.Version,
			InvokePath: env.InvokePath(4), // Skip 4 levels of the call stack
			InputData:  "",
			OutputData: "",
			ExtInfo:    nil,
		},
	}})
	if sendDataErr != nil {
		log.Errorf("logMonitorEvent fail:%v", sendDataErr)
	}
}

// refreshInterval Local cache refresh interval
func refreshInterval(projectID string) uint32 {
	application := GetApplication(projectID)
	if application == nil || application.TabConfig == nil || application.TabConfig.ControlData == nil {
		return defaultRefreshInterval
	}
	controlData := application.TabConfig.ControlData
	if controlData.RefreshInterval <= 0 {
		return defaultRefreshInterval
	}
	return controlData.RefreshInterval
}

// NewAndSetApplication builds a new application, obtains the latest data from the background cache service,
// and automatically updates the local cache after the data is obtained normally
// If the returned application is not empty, it can ensure that ExperimentData, ConfigData,
// and ControlData under TabConfig are not nil,
// Avoid multiple empty judgments in the place where it is used
func NewAndSetApplication(ctx context.Context, projectID string) (application *Application, err error) {
	defer func() {
		recoverErr := recover()
		if recoverErr != nil {
			body := make([]byte, 1<<10)
			runtime.Stack(body, false)
			log.Errorf("recoverErr:%v\n%s", recoverErr, body)
			err = fmt.Errorf("recoverErr:%v\n%s", recoverErr, body)
			return
		}
	}()
	var modified = true
	application, modified, err = refreshApplication(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "refreshApplication")
	}
	if modified { // The local cache needs to be updated only when data changes
		log.Infof("[projectID=%v] version=%v", application.ProjectID, application.Version)
		setApplication(application)
	}
	return application, nil
}

// refreshApplication Refresh the local cache data application, return cache data,
// whether the data is updated, error information
func refreshApplication(ctx context.Context, projectID string) (*Application, bool, error) {
	application := getLocalCacheWithDefault(projectID)
	err := setupTabConfig(ctx, application) // Pull cache data
	if err != nil {
		return nil, false, errors.Wrap(err, "setupTabConfig")
	}
	if application.retryTime > maxRetryTime { // If there is no data change,
		// the silent period will begin when maxRetryTime is reached.
		return application, false, nil
	}
	err = setupLayerIndex(application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupLayerIndex")
	}
	err = setupFullFlowLayerIndex(application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupFullFlowLayerIndex")
	}
	err = setupLayerDomainMetadataListIndex(application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupLayerDomainMetadataListIndex")
	}
	err = setupExperimentBucketInfo(ctx, application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupExperimentBucketInfo")
	}
	err = setupGroupBucketInfo(ctx, application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupGroupBucketInfo")
	}
	err = setupDMPTagInfo(application)
	if err != nil {
		return nil, false, errors.Wrap(err, "setupDMPTagInfo")
	}
	setupMetricsInitConfigIndex(application)
	setupVariantKeyLayerKeyMap(application)
	return application, true, nil
}

func setupVariantKeyLayerKeyMap(application *Application) {
	var variantKeyLayerKeyMap = make(map[string][]string)
	for layerKey, layer := range application.LayerIndex {
		for _, group := range layer.GroupIndex {
			if group.IsDefault {
				for key := range group.Params {
					variantKeyLayerKeyMap[key] = append(variantKeyLayerKeyMap[key], layerKey)
				}
				continue
			}

		}
	}
	application.VariantKeyLayerMap = variantKeyLayerKeyMap
}

func setupMetricsInitConfigIndex(application *Application) {
	application.MetricsPluginInitConfigIndex = application.TabConfig.ControlData.MetricsInitConfigIndex
}

func setupDMPTagInfo(application *Application) error {
	var result = make(map[protoctabcacheserver.UnitIDType]map[int64]map[string]interface{})
	for _, layer := range application.LayerIndex {
		for _, group := range layer.GroupIndex {
			if group == nil {
				return errors.Errorf("[layerKey=%v]group should not be nil", layer.Metadata.Key)
			}
			if group.IssueInfo == nil {
				continue
			}
			for _, tagList := range group.IssueInfo.TagListGroup {
				if tagList == nil {
					continue
				}
				for _, tag := range tagList.TagList {
					if tag == nil {
						return errors.Errorf("[groupID=%d]invalid tag", group.Id)
					}
					if tag.TagType != protoctabcacheserver.TagType_TAG_TYPE_DMP {
						continue
					}
					if result[layer.Metadata.UnitIdType] == nil {
						result[layer.Metadata.UnitIdType] = make(map[int64]map[string]interface{})
					}
					if result[layer.Metadata.UnitIdType][tag.DmpPlatform] == nil {
						result[layer.Metadata.UnitIdType][tag.DmpPlatform] = make(map[string]interface{})
					}
					result[layer.Metadata.UnitIdType][tag.DmpPlatform][tag.Value] = nil
				}
			}
		}
	}
	application.DMPTagInfo = result
	return nil
}

func setupLayerIndex(application *Application) error {
	layerIndex, err := setupLayerIndexDomain(application.TabConfig.ExperimentData.GlobalDomain)
	if err != nil {
		return err
	}
	application.LayerIndex = layerIndex
	return nil
}

func setupFullFlowLayerIndex(application *Application) error {
	fullFlowLayerIndex, err := setupFillFlowLayerIndexDomain(application.TabConfig.ExperimentData.GlobalDomain)
	if err != nil {
		return err
	}
	application.FullFlowLayerIndex = fullFlowLayerIndex
	return nil
}

func setupLayerDomainMetadataListIndex(application *Application) error {
	var result = make(map[string][]*protoctabcacheserver.DomainMetadata)
	err := setupDomainMetadataListIndex(application.TabConfig.ExperimentData.GlobalDomain,
		nil, result)
	if err != nil {
		return err
	}
	application.LayerDomainMetadataListIndex = result
	return nil
}

// setupDomainMetadataListIndex Set the domain metadata of each layer's parent domain, from left to right,
// corresponding to the layer domain structure from top to bottom
func setupDomainMetadataListIndex(domain *protoctabcacheserver.Domain,
	rootMetadataList []*protoctabcacheserver.DomainMetadata,
	result map[string][]*protoctabcacheserver.DomainMetadata) error {
	if domain == nil || domain.Metadata == nil {
		return nil
	}
	for _, holdoutDomain := range domain.HoldoutDomainList {
		if holdoutDomain == nil || holdoutDomain.Metadata == nil {
			return errors.Errorf("invalid domain metadata")
		}
		tmpRootMetadataList := make([]*protoctabcacheserver.DomainMetadata, len(rootMetadataList))
		copy(tmpRootMetadataList, rootMetadataList)
		holdoutLayerMetadataListIndex := append(tmpRootMetadataList, domain.Metadata, holdoutDomain.Metadata)
		for _, layer := range holdoutDomain.LayerList {
			if layer == nil || layer.Metadata == nil {
				return errors.Errorf("invalid layer:%+v", layer)
			}
			result[layer.Metadata.Key] = holdoutLayerMetadataListIndex
		}
	}
	for _, multiLayerDomain := range domain.MultiLayerDomainList {
		if multiLayerDomain == nil || multiLayerDomain.Metadata == nil {
			return errors.Errorf("invalid domain metadata")
		}
		tmpRootMetadataList := make([]*protoctabcacheserver.DomainMetadata, len(rootMetadataList))
		copy(tmpRootMetadataList, rootMetadataList)
		multiLayerMetadataListIndex := append(tmpRootMetadataList, domain.Metadata, multiLayerDomain.Metadata)
		for _, layer := range multiLayerDomain.LayerList {
			if layer == nil || layer.Metadata == nil {
				return errors.Errorf("invalid layer:%+v", layer)
			}
			result[layer.Metadata.Key] = multiLayerMetadataListIndex
		}
	}
	for _, subDomain := range domain.DomainList {
		tmpRootMetadataList := make([]*protoctabcacheserver.DomainMetadata, len(rootMetadataList))
		copy(tmpRootMetadataList, rootMetadataList)
		err := setupDomainMetadataListIndex(subDomain, append(tmpRootMetadataList, domain.Metadata), result)
		if err != nil {
			return errors.Wrap(err, "setupLayerDomainMetadataListIndex")
		}
	}
	return nil
}

// setupFillFlowLayerIndexDomain Get the layer index of the full-flow layer in the layer domain structure.
// If the layer accounts for 100% of the disk traffic, it is a full-flow layer.
// If the current domain has a holdoutDomain and the traffic is not 0,
// then the multi-layer domains and sub-domains under this domain are not full-flow layers.
func setupFillFlowLayerIndexDomain(domain *protoctabcacheserver.Domain) (map[string]*protoctabcacheserver.Layer,
	error) {
	var result = make(map[string]*protoctabcacheserver.Layer)
	if domain == nil {
		return nil, nil
	}
	if !isFullFlowDomain(domain.Metadata) {
		return nil, nil
	}
	for _, holdoutDomain := range domain.HoldoutDomainList {
		if holdoutDomain == nil || holdoutDomain.Metadata == nil || holdoutDomain.Metadata.BucketSize <= 0 {
			return nil, errors.Errorf("invalid domain metadata")
		}
		if isFullFlowDomain(holdoutDomain.Metadata) {
			setLayerIndex(result, holdoutDomain.LayerList)
		}
		if hasTraffic(holdoutDomain.Metadata) {
			return result, nil
		}
	}
	for _, multiLayerDomain := range domain.MultiLayerDomainList {
		if multiLayerDomain == nil || multiLayerDomain.Metadata == nil || multiLayerDomain.Metadata.BucketSize <= 0 {
			return nil, errors.Errorf("invalid domain metadata")
		}
		if isFullFlowDomain(multiLayerDomain.Metadata) {
			setLayerIndex(result, multiLayerDomain.LayerList)
		}
	}
	for _, subdomain := range domain.DomainList {
		item, err := setupFillFlowLayerIndexDomain(subdomain)
		if err != nil {
			return nil, err
		}
		for key, value := range item {
			result[key] = value
		}
	}
	return result, nil
}

func hasTraffic(metadata *protoctabcacheserver.DomainMetadata) bool {
	for _, r := range metadata.TrafficRangeList {
		if r.Left > 0 && r.Left <= r.Right && r.Right <= metadata.BucketSize {
			return true
		}
	}
	return false
}

func setupLayerIndexDomain(domain *protoctabcacheserver.Domain) (map[string]*protoctabcacheserver.Layer, error) {
	var result = make(map[string]*protoctabcacheserver.Layer)
	if domain == nil {
		return nil, nil
	}
	for _, holdoutDomain := range domain.HoldoutDomainList {
		if holdoutDomain == nil || holdoutDomain.Metadata == nil || holdoutDomain.Metadata.BucketSize <= 0 {
			continue
		}
		setLayerIndex(result, holdoutDomain.LayerList)
	}
	for _, multiLayerDomain := range domain.MultiLayerDomainList {
		if multiLayerDomain == nil || multiLayerDomain.Metadata == nil || multiLayerDomain.Metadata.BucketSize <= 0 {
			continue
		}
		setLayerIndex(result, multiLayerDomain.LayerList)
	}
	for _, subdomain := range domain.DomainList {
		item, err := setupLayerIndexDomain(subdomain)
		if err != nil {
			return nil, err
		}
		for key, value := range item {
			result[key] = value
		}
	}
	return result, nil
}

func setLayerIndex(layerIndex map[string]*protoctabcacheserver.Layer, layerList []*protoctabcacheserver.Layer) {
	for _, layer := range layerList {
		if layer == nil || layer.Metadata == nil || layer.Metadata.BucketSize <= 0 {
			continue
		}
		layerIndex[layer.Metadata.Key] = layer
	}
}

func isFullFlowDomain(metadata *protoctabcacheserver.DomainMetadata) bool {
	trafficRangeList := metadata.TrafficRangeList
	length := len(trafficRangeList)
	if length == 0 {
		return false
	}
	sort.SliceStable(trafficRangeList, func(i, j int) bool {
		if trafficRangeList[i].Left < trafficRangeList[j].Left {
			return true
		}
		if trafficRangeList[i].Left == trafficRangeList[j].Left {
			return trafficRangeList[i].Right < trafficRangeList[j].Right
		}
		return false
	})
	if trafficRangeList[0].Left > 1 {
		return false
	}
	var right = trafficRangeList[0].Right
	for _, r := range trafficRangeList {
		if right+1 < r.Left { // 下个区间跟上个区间没有重叠 则判定 false
			return false
		}
		if r.Right > right {
			right = r.Right
		}
	}
	return right >= metadata.BucketSize
}

func setupExperimentBucketInfo(ctx context.Context, application *Application) error {
	if application.retryTime > maxRetryTime {
		return nil
	}
	experimentBucketInfo, err := client.CacheClient.BatchGetExperimentBucketInfo(ctx,
		&protoctabcacheserver.BatchGetExperimentBucketReq{
			ProjectId:          application.ProjectID,
			SdkVersion:         env.SDKVersion,
			BucketVersionIndex: genExperimentVersionIndex(application),
		})
	if err != nil {
		return errors.Wrap(err, "batchGetExperimentBucketInfo")
	}
	if experimentBucketInfo.Code != protoctabcacheserver.Code_CODE_SUCCESS {
		return errors.Errorf("invalid code:%v, message=%s", experimentBucketInfo.Code, experimentBucketInfo.Message)
	}
	for experimentID, bucketInfo := range experimentBucketInfo.BucketIndex {
		if bucketInfo.ModifyType == protoctabcacheserver.ModifyType_MODIFY_DELETE ||
			bucketInfo.ModifyType == protoctabcacheserver.ModifyType_MODIFY_UNKNOWN {
			delete(application.ExperimentIDBucketInfoIndex, experimentID) // this is safe
			delete(application.ExperimentIDRoaringBitmapIndex, experimentID)
			continue
		}
		application.ExperimentIDBucketInfoIndex[experimentID] = bucketInfo
		if bucketInfo.BucketType != protoctabcacheserver.BucketType_BUCKET_TYPE_BITMAP {
			continue
		}
		bitmap := roaring.New()
		_, err = bitmap.FromBuffer(bucketInfo.Bitmap)
		if err != nil {
			return errors.Wrapf(err, "[experimentID=%d]new bitmap fromBuffer", experimentID)
		}
		application.ExperimentIDRoaringBitmapIndex[experimentID] = bitmap
	}
	return nil
}

func setupGroupBucketInfo(ctx context.Context, application *Application) error {
	if application.retryTime > maxRetryTime { // 没有数据变化的情况下，达到 maxRetryTime 则进入静默期
		return nil
	}
	groupVersion := genGroupVersionIndex(application)
	if len(groupVersion) == 0 {
		return nil
	}
	groupBucketInfo, err := client.CacheClient.BatchGetGroupBucketInfo(ctx, &protoctabcacheserver.BatchGetGroupBucketReq{
		ProjectId:          application.ProjectID,
		SdkVersion:         env.SDKVersion,
		BucketVersionIndex: groupVersion,
	})
	if err != nil {
		return errors.Wrap(err, "batchGetGroupBucketInfo")
	}
	if groupBucketInfo.Code != protoctabcacheserver.Code_CODE_SUCCESS {
		return errors.Errorf("invalid code:%v, message=%s", groupBucketInfo.Code, groupBucketInfo.Message)
	}
	for groupID, bucketInfo := range groupBucketInfo.BucketIndex {
		if bucketInfo.ModifyType == protoctabcacheserver.ModifyType_MODIFY_DELETE ||
			bucketInfo.ModifyType == protoctabcacheserver.ModifyType_MODIFY_UNKNOWN {
			delete(application.GroupIDBucketInfoIndex, groupID) // this is safe
			delete(application.GroupIDRoaringBitmapIndex, groupID)
			continue
		}
		application.GroupIDBucketInfoIndex[groupID] = bucketInfo
		if bucketInfo.BucketType != protoctabcacheserver.BucketType_BUCKET_TYPE_BITMAP {
			continue
		}
		bitmap := roaring.New()
		_, err = bitmap.FromBuffer(bucketInfo.Bitmap)
		if err != nil {
			return errors.Wrapf(err, "[groupID=%d]new bitmap fromBuffer", groupID)
		}
		application.GroupIDRoaringBitmapIndex[groupID] = bitmap
	}
	return nil
}

func setupTabConfig(ctx context.Context, application *Application) error {
	tabConfigData, err := client.CacheClient.GetTabConfigData(ctx, &protoctabcacheserver.GetTabConfigReq{
		ProjectId:  application.ProjectID,
		Version:    application.Version,
		SdkVersion: env.SDKVersion,
		UpdateType: protoctabcacheserver.UpdateType_UPDATE_TYPE_COMPLETE,
	})
	if err != nil {
		return errors.Wrap(err, "getTabConfigData")
	}
	if tabConfigData == nil {
		return errors.Errorf("invalid tabConfigData")
	}
	if tabConfigData.Code != protoctabcacheserver.Code_CODE_SUCCESS && tabConfigData.Code !=
		protoctabcacheserver.Code_CODE_SAME_VERSION {
		return errors.Errorf("invalid code:%v, message=%s", tabConfigData.Code, tabConfigData.Message)
	}
	if tabConfigData.Code == protoctabcacheserver.Code_CODE_SUCCESS {
		if !validateTabConfig(tabConfigData) {
			return errors.Errorf("invalid tabConfig")
		}
		application.retryTime = 0
		application.TabConfig = tabConfigData.TabConfigManager.TabConfig
		application.Version = tabConfigData.TabConfigManager.Version
	}
	if tabConfigData.Code == protoctabcacheserver.Code_CODE_SAME_VERSION {
		if application.retryTime <= maxRetryTime {
			application.retryTime++
		}
	}
	return nil
}

func validateTabConfig(tabConfigData *protoctabcacheserver.GetTabConfigResp) bool {
	if tabConfigData.TabConfigManager == nil ||
		tabConfigData.TabConfigManager.TabConfig == nil {
		log.Errorf("invalid tabConfig111:%+v", tabConfigData)
		return false
	}
	tabConfig := tabConfigData.TabConfigManager.TabConfig
	if tabConfig.ExperimentData == nil || tabConfig.ConfigData == nil || tabConfig.ControlData == nil {
		log.Errorf("invalid data111")
		return false
	}
	return true
}

func genExperimentVersionIndex(application *Application) map[int64]string {
	if application == nil {
		return nil
	}
	var result = make(map[int64]string)
	for _, layer := range application.LayerIndex {
		if layer.Metadata.HashType != protoctabcacheserver.HashType_HASH_TYPE_DOUBLE {
			continue
		}
		for _, experiment := range layer.ExperimentIndex {
			if experiment == nil {
				continue
			}
			if experiment.Id == 0 {
				continue
			}
			result[experiment.Id] = ""
		}
	}
	if application.TabConfig != nil && application.TabConfig.ExperimentData != nil &&
		application.TabConfig.ExperimentData.HoldoutData != nil {
		for _, layer := range application.TabConfig.ExperimentData.HoldoutData.HoldoutLayerIndex {
			if layer == nil {
				continue
			}
			for _, experiment := range layer.ExperimentIndex {
				if experiment == nil {
					continue
				}
				if experiment.Id == 0 {
					continue
				}
				result[experiment.Id] = ""
			}
		}
	}
	for experimentID, bucketInfo := range application.ExperimentIDBucketInfoIndex {
		result[experimentID] = bucketInfo.Version
	}
	return result
}

func genGroupVersionIndex(application *Application) map[int64]string {
	var result = make(map[int64]string)
	for _, layer := range application.LayerIndex {
		if layer == nil {
			continue
		}
		for _, group := range layer.GroupIndex {
			if group == nil {
				continue
			}
			result[group.Id] = ""
		}
	}
	if application.TabConfig != nil && application.TabConfig.ExperimentData != nil &&
		application.TabConfig.ExperimentData.HoldoutData != nil {
		for _, layer := range application.TabConfig.ExperimentData.HoldoutData.HoldoutLayerIndex {
			if layer == nil {
				continue
			}
			for _, group := range layer.GroupIndex {
				if group == nil {
					continue
				}
				result[group.Id] = ""
			}
		}
	}
	for groupID, bucketInfo := range application.GroupIDBucketInfoIndex {
		result[groupID] = bucketInfo.Version
	}
	return result
}

func getLocalCacheWithDefault(projectID string) *Application {
	curApplication := GetApplication(projectID)
	if curApplication == nil {
		return &Application{
			ProjectID:                      projectID,
			Version:                        "",
			ExperimentIDBucketInfoIndex:    map[int64]*protoctabcacheserver.BucketInfo{},
			ExperimentIDRoaringBitmapIndex: map[int64]*roaring.Bitmap{},
			GroupIDBucketInfoIndex:         map[int64]*protoctabcacheserver.BucketInfo{},
			GroupIDRoaringBitmapIndex:      map[int64]*roaring.Bitmap{},
			FullFlowLayerIndex:             map[string]*protoctabcacheserver.Layer{},
			LayerIndex:                     map[string]*protoctabcacheserver.Layer{},
			VariantKeyLayerMap:             map[string][]string{},
			DMPTagInfo:                     map[protoctabcacheserver.UnitIDType]map[int64]map[string]interface{}{},
		}
	}
	// copy curApplication
	return getNewApplication(curApplication)
}

// getNewApplication gets a concurrently safe application, regenerates a new map object,
// but retains the original pointer of value
// Here you can directly use json marshal unmarshal to make a deep copy of the data to ensure concurrent safety,
// but the data volume is relatively large and the deep copy performance is poor
// So here we make a shallow copy for the fields with concurrency issues
func getNewApplication(curApplication *Application) *Application {
	return &Application{
		ProjectID:                      curApplication.ProjectID,
		Version:                        curApplication.Version,
		TabConfig:                      curApplication.TabConfig,
		ExperimentIDBucketInfoIndex:    getNewBucketInfoIndex(curApplication.ExperimentIDBucketInfoIndex),
		ExperimentIDRoaringBitmapIndex: getNewRoaringBitmapIndex(curApplication.ExperimentIDRoaringBitmapIndex),
		GroupIDBucketInfoIndex:         getNewBucketInfoIndex(curApplication.GroupIDBucketInfoIndex),
		GroupIDRoaringBitmapIndex:      getNewRoaringBitmapIndex(curApplication.GroupIDRoaringBitmapIndex),
		FullFlowLayerIndex:             curApplication.FullFlowLayerIndex,
		LayerIndex:                     curApplication.LayerIndex,
		DMPTagInfo:                     curApplication.DMPTagInfo,
		VariantKeyLayerMap:             curApplication.VariantKeyLayerMap,
		PreparedDMPTag:                 curApplication.PreparedDMPTag,
		DisableDMPTag:                  curApplication.DisableDMPTag,
		retryTime:                      curApplication.retryTime,
	}
}

func getNewBucketInfoIndex(
	curIndex map[int64]*protoctabcacheserver.BucketInfo) map[int64]*protoctabcacheserver.BucketInfo {
	var result = make(map[int64]*protoctabcacheserver.BucketInfo, len(curIndex))
	for k, v := range curIndex {
		result[k] = v
	}
	return result
}

func getNewRoaringBitmapIndex(curIndex map[int64]*roaring.Bitmap) map[int64]*roaring.Bitmap {
	var result = make(map[int64]*roaring.Bitmap, len(curIndex))
	for k, v := range curIndex {
		result[k] = v
	}
	return result
}

// GetApplication TODO
func GetApplication(projectID string) *Application {
	result, ok := localApplicationCache.Load(projectID)
	if !ok {
		return nil
	}
	application, ok := result.(*Application)
	if !ok {
		return nil
	}
	return application
}

func setApplication(application *Application) {
	if application == nil {
		return
	}
	localApplicationCache.Store(application.ProjectID, application)
}

// Release TODO
func Release() {
	localApplicationCache = sync.Map{}
}

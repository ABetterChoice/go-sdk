// Package experiment TODO
package experiment

import (
	"github.com/abetterchoice/go-sdk/internal/cache"
)

// Options abtest experiment diversion related options, life cycle for each abtest diversion session
type Options struct {
	// Whether TAB automatically records exposure
	IsExposureLoggingAutomatic bool `json:"ela,omitempty"`
	// Scene ID, used for filtering, returns the hit experiment groups in the lower layers of these scenes,
	// key is sceneID, value is whether it passes, if sceneIDs is empty, all passes
	SceneIDs map[int64]bool `json:"sceneIDs,omitempty"`
	// Layer key, used for filtering, returns the hit experiment groups under these layers,
	// key is layerKey, value is whether it passes, if layerKeys is empty, all passes
	LayerKeys map[string]bool `json:"layerKeys,omitempty"`
	// Experiment key, used for filtering, returns the experimental group hit in this experiment.
	// Since there are multiple experiments in the same layer, if the user hits other experiments in the same layer,
	// the experimental groups of other experiments will not be returned, but the default experiment will be returned.
	ExperimentKeys map[string]bool `json:"experimentKeys,omitempty"`
	// Whether to pre-process the dmp tag. If enabled, if there is a dmp tag under the business,
	// rpc will first access the dmp service to get the hit or not result. Enabled by default
	IsPreparedDMPTag bool `json:"isPreparedDmpTag,omitempty"`
	// The key format is {{reporting ID}}-{{DMP platform ID}}-{{crowd package key}},
	// and the value is whether it is a hit
	DMPTagResult map[string]bool `json:"dmpTagResult,omitempty"`
	// Whether to disable dmp. If disabled, no rpc will be initiated.
	// ab test will offload all local calculations and will not hit dmp by default
	IsDisableDMP bool `json:"isDisableDmp,omitempty"`
	// Whitelist, key is the layer, value is the experiment group specified by the user under the layer
	OverrideList map[string]int64 `json:"-"`
	// Attribute tag information owned by unitID
	AttributeTag map[string][]string `json:"-"`
	// unitID passed in by the user
	UnitID string `json:"unitId,omitempty"`
	// The ID used for hashing is usually consistent with the unitID.
	// The same UnitID will be mapped to the same DecisionID, that is,
	// the same UnitID will stably hit the same experimental group
	DecisionID string `json:"decisionId,omitempty"`
	// The newID passed in by the user is usually empty. If the user account system is switching to grayscale,
	// both unitID and newUnitID can be passed in. In abtest diversion,
	// according to the layer properties set by the web platform,
	// the specified unitID or newUnitID is selected for reporting
	NewUnitID string `json:"newUnitId,omitempty"`
	// ID used for hashing, usually consistent with newUnitID.
	// If abtest splits traffic and the layer specifies to use NewUnitID,
	// NewDecisionID will be used as the input of hashing. The same NewUnitID and the same NewDecisionID
	// will stably hit the same experimental group.
	NewDecisionID string `json:"newDecisionId,omitempty"`
	// Cache data snapshot
	Application *cache.Application `json:"-"`
	// The result of the holdout layer hit. If it is nil, it means that it is not held out.
	HoldoutLayerResult map[string]*Experiment `json:"-"`
}

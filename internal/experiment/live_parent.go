// Package experiment implements LIVE_PARENT traffic routing, i.e. parent-child experiments.
package experiment

import (
	"context"

	"github.com/abetterchoice/go-sdk/plugin/log"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/pkg/errors"
)

// liveParentCtx carries the recursion state of the parent-hit evaluation across a whole
// parent-child chain (A -> B -> C) within a single layer decision.
//
// It replaces the simpler "never route while evaluating a parent" guard: that guard is
// cycle-free by construction but also caps chains at depth two. Carrying an explicit stack
// detects cycles precisely (only a real return to a layer already on the stack is rejected)
// while letting acyclic multi-level chains resolve.
type liveParentCtx struct {
	// visited holds the layer keys currently on the recursion stack, pushed and popped around
	// each parent evaluation, so self references and mutual references are caught.
	visited map[string]struct{}
	// memo caches parent layer results already evaluated during this decision, keyed by layer key.
	// Sharing it across branches keeps the cost linear in the number of layers rather than
	// exponential: a diamond-shaped reference evaluates the shared parent only once, which also
	// avoids re-issuing DMP lookups when the parent experiment is tag-based.
	memo map[string]*Experiment
}

func newLiveParentCtx() *liveParentCtx {
	return &liveParentCtx{visited: make(map[string]struct{}, 2)}
}

// enter pushes a layer onto the recursion stack; false means a cycle, and the caller must skip
// that parent rather than descend into it.
func (c *liveParentCtx) enter(layerKey string) bool {
	if _, onStack := c.visited[layerKey]; onStack {
		return false
	}
	if c.visited == nil {
		c.visited = make(map[string]struct{}, 2)
	}
	c.visited[layerKey] = struct{}{}
	return true
}

func (c *liveParentCtx) leave(layerKey string) {
	delete(c.visited, layerKey)
}

// getLiveParentExperiment resolves the child experiment for LIVE_PARENT routes on this layer.
//
// Returns the hit child group, whether this layer is governed by a resolvable LIVE_PARENT route,
// and an error. "Governed" rather than "parent hit" is the condition the caller acts on: when the
// parent chain resolves but the user missed the parent, the layer's ordinary hash must be skipped,
// otherwise the user would be assigned into a child group anyway and the parent-child relation
// would mean nothing. An unresolvable chain (parent version not indexed, or a cycle) leaves the
// layer untouched, as if the route were absent.
func (e *executor) getLiveParentExperiment(ctx context.Context, layer *protoccacheserver.Layer,
	options *Options, parentCtx *liveParentCtx) (*Experiment, bool, error) {
	governed := false
	for _, routing := range layer.TrafficRouting {
		if routing == nil || routing.RoutingMode !=
			protoccacheserver.ExpGroupTrafficRoutingMode_EXP_GROUP_TRAFFIC_ROUTING_MODE_LIVE_PARENT {
			continue
		}
		childExperiment, ok := layer.ExperimentIndex[routing.TargetExpGroupId]
		if !ok || childExperiment == nil {
			// Child experiment missing from this layer is a config anomaly; treat the route as
			// absent so the layer keeps its ordinary behaviour instead of going dark.
			log.Warnf("live parent target experiment[%d] not found, layerKey=%s, skip this routing",
				routing.TargetExpGroupId, layer.Metadata.Key)
			continue
		}
		for _, route := range routing.Routes {
			if route == nil || route.SourceExperimentId == 0 {
				continue
			}
			parentHit, resolvable, err := e.isLiveParentHit(ctx, route.SourceExperimentId, options, parentCtx)
			if err != nil {
				return nil, false, err
			}
			if resolvable {
				governed = true
			}
			if !parentHit {
				continue
			}
			// Parent hit: assign within the child experiment using its own hash, so the child can
			// still randomise its own versions. The routing config carries no target version for
			// LIVE_PARENT (the platform leaves it empty by design), hence no direct bucket reuse.
			experiment, err := e.getExperimentGroup(ctx, childExperiment, layer, options)
			if err != nil {
				return nil, true, errors.Wrap(err, "getExperimentGroup")
			}
			if experiment == nil {
				continue
			}
			return experiment, true, nil
		}
	}
	return nil, governed, nil
}

// isLiveParentHit reports whether the user hit the given parent version.
//
// This is a shadow evaluation: the parent layer is re-decided, but its outcome is not this
// request's result, so it must stay off the output channels (see newLiveParentShadowOptions).
// Judging by bucket position alone would be wrong, because hitting a group also requires its
// audience conditions to match (see IsHitTag in the ordinary path).
//
// Multi-level chains: if the parent layer itself carries LIVE_PARENT routes, evaluating it
// descends further, forming A -> B -> C; only a return to a layer already on the stack is cut.
//
// Known boundary: the shadow evaluation starts at the parent layer and does not replay the
// domain/layer traffic admission above it. If the parent layer is not a full-traffic layer, a
// user who was never routed into it may still be judged as having hit the parent, which widens
// child admission. This mirrors the platform contract that a parent layer takes full traffic;
// supporting partial-traffic parent layers requires replaying admission top-down here.
func (e *executor) isLiveParentHit(ctx context.Context, sourceGroupID int64, options *Options,
	parentCtx *liveParentCtx) (hit bool, resolvable bool, err error) {
	parentLayer := options.Application.LiveParentSourceLayerIndex[sourceGroupID]
	if parentLayer == nil || parentLayer.Metadata == nil {
		return false, false, nil
	}
	if parentCtx == nil {
		parentCtx = newLiveParentCtx()
	}
	experiment, evaluated, err := e.shadowEvalParentLayer(ctx, parentLayer, options, parentCtx)
	if err != nil {
		return false, false, err
	}
	if !evaluated { // cycle: this chain is unresolvable and must not govern the layer
		return false, false, nil
	}
	if experiment == nil || experiment.Group == nil {
		return false, true, nil
	}
	return !experiment.Group.IsDefault && experiment.Group.Id == sourceGroupID, true, nil
}

// shadowEvalParentLayer evaluates a parent layer in an isolated context; cycle detection and
// result caching are both funnelled through here. The second return value reports whether the
// evaluation ran: false means the layer is already on the recursion stack, and descending would
// loop forever on what is in any case a misconfiguration.
func (e *executor) shadowEvalParentLayer(ctx context.Context, parentLayer *protoccacheserver.Layer,
	options *Options, parentCtx *liveParentCtx) (*Experiment, bool, error) {
	layerKey := parentLayer.Metadata.Key
	if cached, ok := parentCtx.memo[layerKey]; ok {
		return cached, true, nil
	}
	if !parentCtx.enter(layerKey) {
		return nil, false, nil
	}
	experiment, err := e.getLayerExperimentInternal(ctx, parentLayer,
		newLiveParentShadowOptions(options), parentCtx)
	parentCtx.leave(layerKey)
	if err != nil {
		return nil, true, errors.Wrap(err, "shadowEvalParentLayer")
	}
	if parentCtx.memo == nil {
		parentCtx.memo = make(map[string]*Experiment, 2)
	}
	parentCtx.memo[layerKey] = experiment
	return experiment, true, nil
}

// newLiveParentShadowOptions builds the context used to evaluate a parent layer.
//
// The shadow evaluation only answers "did the user hit the parent", so it keeps the inputs that
// decision needs and drops the query-scoped filters. Those filters matter: a caller may request
// only the child layer, and applying its filters to the parent layer would report a miss and
// silently disable the parent-child relation.
//
// DMPTagResult and OverrideList are shared on purpose. The former is prepared up front and keyed
// by tag, so sharing avoids repeating DMP lookups within one request; the latter must apply to
// the parent layer too, since a user forced into a specific parent group should be treated as
// having hit it.
func newLiveParentShadowOptions(options *Options) *Options {
	shadow := *options
	shadow.SceneIDs = nil
	shadow.LayerKeys = nil
	shadow.ExperimentKeys = nil
	return &shadow
}

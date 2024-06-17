// Package experiment abtest
package experiment

import (
	"context"

	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
)

func sceneIDListFilter(ctx context.Context, layer *protoccacheserver.Layer, options *Options) bool {
	if len(options.SceneIDs) == 0 {
		return true
	}
	for i := range layer.Metadata.SceneIdList {
		if options.SceneIDs[layer.Metadata.SceneIdList[i]] {
			return true
		}
	}
	return false
}

func layerKeyFilter(ctx context.Context, layer *protoccacheserver.Layer, options *Options) bool {
	if len(options.LayerKeys) == 0 {
		return true
	}
	return options.LayerKeys[layer.Metadata.Key]
}

func experimentKeyFilter(ctx context.Context, layer *protoccacheserver.Layer, options *Options) bool {
	if len(options.ExperimentKeys) == 0 {
		return true
	}
	for _, experiment := range layer.ExperimentIndex {
		if options.ExperimentKeys[experiment.Key] {
			return true
		}
	}
	return false
}

package experiment

import (
	"context"

	"github.com/abetterchoice/go-sdk/internal/cache"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/pkg/errors"
)

// GetDefaultExperiments The default experiment on the acquisition layer does not involve any diversion process,
// and is mainly aimed at obtaining a bottom-up strategy.
func (e *executor) GetDefaultExperiments(ctx context.Context, projectID string,
	options *Options) (map[string]*Experiment,
	error) {
	application := cache.GetApplication(projectID)
	if application == nil {
		return nil, errors.Errorf("projectID [%s] not found", projectID)
	}
	options.Application = application
	return e.getDomainDefaultExperiments(ctx, application.TabConfig.ExperimentData.GlobalDomain, options)
}

func (e *executor) getDomainDefaultExperiments(ctx context.Context, domain *protoc_cache_server.Domain,
	options *Options) (
	map[string]*Experiment, error) {
	if domain == nil {
		return nil, nil
	}
	var result = make(map[string]*Experiment)
	if err := e.getHoldoutDomainDefaultExperiments(ctx, domain.HoldoutDomainList, options, result); err != nil {
		return nil, errors.Wrap(err, "getHoldoutDomainDefaultExperiments")
	}
	if err := e.getMultiLayerDomainDefaultExperiments(ctx, domain.MultiLayerDomainList, options, result); err != nil {
		return nil, errors.Wrap(err, "getMultiLayerDomainDefaultExperiments")
	}
	if err := e.getDomainListDefaultExperiments(ctx, domain.DomainList, options, result); err != nil {
		return nil, errors.Wrap(err, "getDomainDefaultExperiments")
	}
	return result, nil
}

func (e *executor) getDomainListDefaultExperiments(ctx context.Context, domains []*protoc_cache_server.Domain,
	options *Options,
	result map[string]*Experiment) error {
	for _, domain := range domains {
		item, err := e.getDomainDefaultExperiments(ctx, domain, options)
		if err != nil {
			return errors.Wrap(err, "getDomainDefaultExperiments")
		}
		for key, value := range item {
			result[key] = value
		}
	}
	return nil
}

func (e *executor) getMultiLayerDomainDefaultExperiments(ctx context.Context,
	multiLayerDomains []*protoc_cache_server.MultiLayerDomain, options *Options,
	result map[string]*Experiment) error {
	for _, multiLayerDomain := range multiLayerDomains {
		for _, layer := range multiLayerDomain.LayerList {
			if layer == nil || layer.Metadata == nil || layer.Metadata.DefaultGroup == nil {
				continue
			}
			if !e.isLayerFilterPass(ctx, layer, options) {
				continue
			}
			result[layer.Metadata.Key] = &Experiment{Group: layer.Metadata.DefaultGroup}
		}
	}
	return nil
}

func (e *executor) getHoldoutDomainDefaultExperiments(ctx context.Context,
	holdoutDomains []*protoc_cache_server.HoldoutDomain, options *Options, result map[string]*Experiment) error {
	for _, holdoutDomain := range holdoutDomains {
		for _, layer := range holdoutDomain.LayerList {
			if layer == nil || layer.Metadata == nil || layer.Metadata.DefaultGroup == nil {
				continue
			}
			if !e.isLayerFilterPass(ctx, layer, options) {
				continue
			}
			result[layer.Metadata.Key] = &Experiment{Group: layer.Metadata.DefaultGroup}
		}
	}
	return nil
}

// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval,
// user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"fmt"
)

// Context // This interface offers the primary APIs for retrieving the results of experiment splitting and
// feature flag evaluations.
//
//go:generate mockgen -source=client.go -package=abc -destination client_mock.go
type Context interface {
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
	GetExperiment(ctx context.Context, projectID string, layerKey string, opts ...ExperimentOption) (*ExperimentResult,
		error)

	// GetExperiments is a batch version of GetExperiment(). It returns the experiment assignments
	// for all experiment layers under the specified project ID. This API is typically used when users
	// need to fetch all experiment results and pass them to downstream services.
	//
	// However, using this API is not recommended as it may lead to logging a large number of unexpected exposures,
	// which could dilute your results. To mitigate this issue, exposures are not logged automatically by default.
	// Instead, you may need to use the exposure logging API to manually log the exposures.
	GetExperiments(ctx context.Context, projectID string, opts ...ExperimentOption) (*ExperimentList, error)

	// GetFeatureFlag evaluates the feature flag with the specified key within the given project.
	// The unit ID is extracted from the provided context.
	// Users can use the strongly-typed parameter retrieval APIs of the FeatureFlag object
	// to fetch the parameter value for this unit ID.
	//
	// Parameters:
	// ctx: golang context.
	// projectID: The ID of the project under which the feature flag is evaluated.
	// key: The key that identifies the feature flag.
	// opts: Additional configuration options.
	//
	// Returns:
	// A pointer to the evaluated FeatureFlag object and an error, if any occurred during the evaluation.
	GetFeatureFlag(ctx context.Context, projectID string, key string, opts ...ConfigOption) (*FeatureFlag, error)

	// GetValueByVariantKey retrieves the parameter value using the globally unique parameter key.
	// It first attempts to find the value among all parameters under all experiments within the given project ID.
	// If the parameter value is not found, it continues to search for this parameter among all feature flags
	// under the project ID until the value is found.
	//
	// Parameters:
	// ctx: The context for the operation.
	// projectID: The ID of the project where the search is performed.
	// key: The globally unique key that identifies the parameter.
	// opts: Additional experiment options.
	//
	// Returns:
	// A pointer to the ValueResult object containing the found value and an error, if any occurred during the search.
	GetValueByVariantKey(ctx context.Context, projectID string, key string, opts ...ExperimentOption) (
		*ValueResult, error)

	// GetRemoteConfig TODO
	// Deprecated: GetRemoteConfig is deprecated, please use GetFeatureFlag instead.
	GetRemoteConfig(ctx context.Context, projectID string, key string, opts ...ConfigOption) (*ConfigResult, error)
}
type userContext struct {
	// Save the error that may occur in NewUserContext to facilitate users to call the API in a chain and
	// handle it uniformly
	err error

	// User attribute tags such as gender, nationality, and age.
	tags map[string][]string

	// The identifier used for traffic splitting and exposure logging can be a user ID,
	// session ID, game room ID, or any other valid unique experimentation identifier. For the same unit ID,
	// it will consistently be assigned to the same experiment across different calls.
	// The whitelist configuration of the experimentation system should use either the unit ID or newUnitID.
	unitID string

	// This is the ID used for traffic splitting. Typically,
	// there's no need to set it as its default value is the same as the unitID. However,
	// in some extreme scenarios, if the ID used in exposure logging differs from the traffic splitting ID,
	// then we can set this ID.
	decisionID string

	// In most cases, there's no need to set this ID.
	// It's primarily used during the migration process, for instance,
	// when part of the experiment layers want to use a different identifier for
	// traffic splitting and exposure logging. This can be used in conjunction with the platform configuration system.
	newUnitID string

	// This ID is only required when you need to separate the traffic splitting and
	// exposure logging ID during the migration stage, when the newUnitID is included.
	newDecisionID string

	// Extended details of the logged exposure.
	// This information will be logged as an additional field in the exposure table,
	// in a format similar to k1=v1; k1=v2.
	expandedData map[string]string
}

// Attribution Pass in each option as needed, including but not limited to setting label information, etc.
type Attribution func(c *userContext)

// NewUserContext In this function, each option in 'opts' is executed sequentially,
// resulting in a new user context. This new context is then used for subsequent experimentation and
// feature flag retrieval. If an internal error occurs, 'err' will be stored in the user context.
// The next time the user context is used, it will short-circuit to simplify chain calls for developers.
// example:
//
//	NewUserContext("123456xA").GetExperiment(context.TODO(), "layerKey_AABB")
//
// The web platform supports whitelisting. If the unitID is in the whitelist, it will take effect.
func NewUserContext(unitID string, opts ...Attribution) Context {
	userCtx := &userContext{
		unitID: unitID,
		tags:   map[string][]string{}, // 避免 tags 为 nil
	}
	for _, opt := range opts {
		opt(userCtx)
	}
	return settingNewUnitIDAndNewDecisionID(userCtx)
}

func settingNewUnitIDAndNewDecisionID(userCtx *userContext) *userContext {
	if len(userCtx.unitID) == 0 { // The exposure logging ID cannot be empty
		userCtx.err = fmt.Errorf("unitID is required")
		return userCtx
	}
	if len(userCtx.decisionID) == 0 { // if the traffic splitting ID is empty, the exposure logging ID will be used
		userCtx.decisionID = userCtx.unitID
	}
	emptyNewUnitID := false
	if len(userCtx.newUnitID) == 0 {
		/*
			Generally, no settings are required,
			and it is mainly used for the grayscale switching process of the account system.
			For example, some experimental layers want to use newUnitID for reporting and offloading,
			while others still use unitID for reporting and offloading.
			This can be used in conjunction with the web backend management system.
		*/
		// If empty, unitID is used
		emptyNewUnitID = true
		userCtx.newUnitID = userCtx.unitID
	}
	if len(userCtx.newDecisionID) == 0 {
		if emptyNewUnitID { // If newUnitID is not configured, the decisionID of unitID is used for offloading.
			userCtx.newDecisionID = userCtx.decisionID
		} else {
			userCtx.newDecisionID = userCtx.newUnitID
		}
	}
	return userCtx
}

// WithTags Set the tag information. If it has not been set, directly replace the tags of userContext,
// otherwise copy it to the tags of userContext in turn.
func WithTags(tags map[string][]string) Attribution {
	return func(c *userContext) {
		if len(tags) == 0 {
			return
		}
		if len(c.tags) == 0 { // direct replacement
			c.tags = tags
			return
		}
		for key, value := range tags { // cover
			c.tags[key] = value
		}
		return
	}
}

// WithTagKV Set the label information kv, if the key already exists, append it
func WithTagKV(key, value string) Attribution {
	return func(c *userContext) {
		c.tags[key] = append(c.tags[key], value)
	}
}

// WithDecisionID The default is consistent with unitID. The ID used for offloading generally does not need to be set.
// If the reporting and offloading IDs need to be separated, this can be achieved by setting decisionID
func WithDecisionID(decisionID string) Attribution {
	return func(c *userContext) {
		if len(decisionID) == 0 { // Empty decisionID is illegal
			c.err = fmt.Errorf("decisionID is required")
			return
		}
		c.decisionID = decisionID
	}
}

// WithNewUnitID Generally, no setting is required.
// It is mainly used for the grayscale process of switching the account system.
// For example, some experimental layers hope to use newUnitID for reporting and offloading.
// Others still use unitID for reporting and offloading,
// which can be used in conjunction with the web backend management system.
// If there is no setting, when encountering an experimental layer that needs to use newUnitID,
// the decision shunt corresponding to unitID will be used by default.
func WithNewUnitID(newUnitID string) Attribution {
	return func(c *userContext) {
		if len(newUnitID) == 0 { // empty newUnitID is illegal
			c.err = fmt.Errorf("newUnitID is required")
			return
		}
		c.newUnitID = newUnitID
	}
}

// WithNewDecisionID Generally, there is no need to set it. The default is consistent with newUnitID.
// If the IDs of reporting and offloading need to be separated, this can be achieved by setting newDecisionID.
// The default is consistent with newUnitID
func WithNewDecisionID(newDecisionID string) Attribution {
	return func(c *userContext) {
		if len(newDecisionID) == 0 { // empty newDecisionID is illegal
			c.err = fmt.Errorf("newDecisionID is required")
			return
		}
		c.newDecisionID = newDecisionID
	}
}

// WithExpandedData Extended information, when exposure is reported,
// this part of the information will be reported to the extended field of the exposure table,
// and stored in the form k1=v1;k1=v2
func WithExpandedData(expandedData map[string]string) Attribution {
	return func(c *userContext) {
		if len(c.expandedData) == 0 {
			c.expandedData = expandedData
			return
		}
		for key, value := range expandedData { // cover
			c.expandedData[key] = value
		}
	}
}

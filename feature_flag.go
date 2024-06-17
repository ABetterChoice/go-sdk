// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval, user feature flag management, exposure data reporting, and logger registration.
package abc

import "context"

// GetFeatureFlag Specify projectID and configuration key to obtain the specific configuration hit by the user,
// similar to abtest experimental distribution, but without the complex experimental layer domain structure
// determine the specific configuration value hit based on the userContext information.
// different unitIDs may hit different configurations,
// but the same unitID will stably hit the same configuration value.
// for more examples see example/feature_flag_test.go
func (c *userContext) GetFeatureFlag(ctx context.Context, projectID string, key string,
	opts ...ConfigOption) (*FeatureFlag, error) {
	config, err := c.GetRemoteConfig(ctx, projectID, key, opts...)
	if err != nil {
		return nil, err
	}
	return &FeatureFlag{ConfigResult: config}, err
}

// FeatureFlag Config
type FeatureFlag struct {
	*ConfigResult
}

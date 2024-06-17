// Package env TODO
package env

import "fmt"

// const Related error codes
const (
	projectIDNotFound = 1001
	invalidGroupID    = 1002
)

var (
	// ErrParamKeyNotFound Experiment parameter key not found
	ErrParamKeyNotFound = fmt.Errorf("param key not found")
)

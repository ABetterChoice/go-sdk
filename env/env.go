// Package env Related enumeration value definitions, etc.
package env

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Type Environment Type
type Type = string

// const Environment type enumeration, only provides formal environment test environment
const (
	TypePrd Type = "prd"
	// TypeTest TODO
	// Deprecated: By default, all use the same formal environment
	TypeTest Type = "test"
)

const (
	// DefaultGlobalGroupID global default experiment group ID. If no experiment is hit,
	// the experiment group ID is returned.
	// It is worth noting that if the global default ID is configured on the web platform,
	// the specified global default ID is returned.
	// Deprecated: It is recommended to use IsDefault in Group to determine whether the default experiment
	// will be run [if the control group or experimental group is not hit, the default experiment may be hit]
	// If no experimental group is hit
	// case1: If the layer has a default experiment set, the layer default experiment is returned
	// case2: If the layer has no default experiment set, and there is no experiment running on the layer,
	// nil, nil is returned
	// case3: If the layer has no default experiment set, and there is an experiment running on the layer,
	// the system default experiment is returned, and err is nil
	DefaultGlobalGroupID = -1
	// DefaultGlobalGroupKey The global default experiment group key. If no experiment is matched,
	// the experiment group ID is returned.
	// Deprecated: It is recommended to use IsDefault in Group to determine whether the default experiment will be run.
	// [If the control group or experimental group is not hit, the default experiment may be hit.]
	DefaultGlobalGroupKey = "defaultSystemGroupKey"
)

// InvokePath Call Path
func InvokePath(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

// TypeBySecretKey parses TAB_DATA_ENV from TAB JWT secret key payload.
// It returns false when payload has no TAB_DATA_ENV.
func TypeBySecretKey(secretKey string) (Type, bool) {
	pairs := strings.Split(secretKey, ".")
	if len(pairs) != 3 {
		return "", false
	}
	payload, err := base64.RawURLEncoding.DecodeString(pairs[1])
	if err != nil {
		return "", false
	}
	type tokenPayload struct {
		TABDataEnv string `json:"TAB_DATA_ENV"`
	}
	data := &tokenPayload{}
	if err = json.Unmarshal(payload, data); err != nil {
		return "", false
	}
	value := strings.ToLower(strings.TrimSpace(data.TABDataEnv))
	if value == "" {
		return "", false
	}
	return Type(value), true
}

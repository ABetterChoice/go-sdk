// Package env TODO
package env

import (
	"encoding/json"

	"github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
)

// SamplingInterval Select sampling interval based on error
func SamplingInterval(config *protoc_cache_server.MetricsConfig, err error) uint32 {
	if err != nil {
		return config.ErrSamplingInterval
	}
	return config.SamplingInterval
}

// ErrMsg error message
func ErrMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// JSONString Object serialization into json string
func JSONString(source interface{}) string {
	if source == nil {
		return ""
	}
	data, _ := json.Marshal(source)
	return string(data)
}

// EventStatus Event Status
func EventStatus(err error) protoc_event_server.MonitorEventStatus {
	if err != nil {
		return protoc_event_server.MonitorEvent_STATUS_UNEXPECTED
	}
	return protoc_event_server.MonitorEvent_STATUS_SUCCESS
}

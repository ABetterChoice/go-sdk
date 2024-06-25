// Package metrics 通用监控上报
package metrics

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/abetterchoice/go-sdk/plugin/log"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
	"github.com/pkg/errors"
)

// Client Monitoring abstract class
type Client interface {
	// Name plugin name
	Name() string
	// Init Initialize the plugin. Multiple initializations require idempotency.
	Init(ctx context.Context, config *protoc_cache_server.MetricsInitConfig) error
	// LogExposure Recording exposure
	LogExposure(ctx context.Context, metadata *Metadata, exposureGroup *protoc_event_server.ExposureGroup) error
	// LogEvent Logging Events
	LogEvent(ctx context.Context, metadata *Metadata, eventGroup *protoc_event_server.EventGroup) error
	// LogMonitorEvent Record monitoring events
	LogMonitorEvent(ctx context.Context, metadata *Metadata,
		monitorEventGroup *protoc_event_server.MonitorEventGroup) error
	// SendData Send data, general interface, reserved
	SendData(ctx context.Context, metadata *Metadata, data [][]string) error
}

// metrics client factory
var (
	clientFactory = map[string]Client{} // Plug-in, supports multiple monitoring reports
	rwMutex       sync.RWMutex
)

const (
	// InitConfigKvToken TODO
	InitConfigKvToken = "system_token"
)

// RegisterClient Registration indicator reporting plug-in implementation
func RegisterClient(client Client) {
	if client == nil {
		return
	}
	if client.Name() == "" { // invalid plugin name
		return
	}
	rwMutex.Lock()
	defer rwMutex.Unlock()
	clientFactory[client.Name()] = client
}

// GetClient Get the monitoring reporting plugin client
func GetClient(name string) (Client, bool) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	c, ok := clientFactory[name]
	return c, ok
}

// WalkFunc Traverse clientFactory, if h returns an error, exit WalkFunc
func WalkFunc(h func(name string, client Client) error) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	for pluginName, c := range clientFactory {
		err := h(pluginName, c)
		if err != nil {
			return err
		}
	}
	return nil
}

var sendDataHook func(metadata *Metadata, data [][]string) error

// RegisterSendDataHook registers the reporting hook. If you need to report multiple channels,
// you can register this hook. The hook will be called synchronously when reporting.
// When the hook function fails to execute, the report will be terminated. If you do not want to terminate,
// please return err = nil.
// In fact, users can use this hook to customize monitoring reporting.
func RegisterSendDataHook(handler func(metadata *Metadata, data [][]string) error) {
	if handler != nil {
		return
	}
	sendDataHook = handler
}

var logExposureHook func(metadata *Metadata, group *protoc_event_server.ExposureGroup) error

// RegisterLogExposureHook registers the exposure data recording hook.
// If you need to report in multiple ways, you can register this hook.
// The hook will be called synchronously when reporting.
// When the hook function fails to execute, the report will be terminated.
// If you do not want to terminate, please return err = nil.
// In fact, users can use this hook to customize monitoring reporting
func RegisterLogExposureHook(handler func(metadata *Metadata, group *protoc_event_server.ExposureGroup) error) {
	if handler != nil {
		return
	}
	logExposureHook = handler
}

// SendData sends data, reports in multiple channels. If the passed clientNames have been registered,
// it will be reported Here, the plug-in interface is directly called for reporting. Generally,
// monitoring reporting components have asynchronous reporting functions,
// so unified asynchronous reporting is not performed here
// Report the specified monitoring reporting plug-in metadata.MetricsPluginName according to the specific event
func SendData(ctx context.Context, metadata *Metadata, data [][]string) (err error) {
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
	if len(data) == 0 {
		return nil
	}
	if !SamplingResult(metadata.SamplingInterval) {
		return nil
	}
	if sendDataHook != nil {
		err := sendDataHook(metadata, data)
		if err != nil {
			return errors.Wrap(err, "sendDataHook")
		}
	}
	c, ok := GetClient(metadata.MetricsPluginName)
	if !ok {
		return nil
	}
	return c.SendData(ctx, metadata, data)
}

// LogExposure sends data and reports in multiple ways. If the clientNames passed in have been registered,
// it will be reported. Here, the plug-in interface is directly called for reporting. Generally,
// monitoring reporting components have asynchronous reporting functions,
// so unified asynchronous reporting is not performed here
// Report the specified monitoring reporting plug-in metadata.MetricsPluginName according to the specific event
func LogExposure(ctx context.Context, metadata *Metadata, group *protoc_event_server.ExposureGroup) (err error) {
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
	if group == nil || len(group.Exposures) == 0 {
		return nil
	}
	if !SamplingResult(metadata.SamplingInterval) {
		return nil
	}
	if logExposureHook != nil {
		err := logExposureHook(metadata, group)
		if err != nil {
			return errors.Wrap(err, "logExposureHook")
		}
	}
	c, ok := GetClient(metadata.MetricsPluginName)
	if !ok {
		return nil
	}
	return c.LogExposure(ctx, metadata, group)
}

// LogMonitorEvent Report the specified monitoring reporting plug-in metadata.MetricsPluginName
// according to the specific event
func LogMonitorEvent(ctx context.Context, metadata *Metadata, group *protoc_event_server.MonitorEventGroup) (
	err error) {
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
	if group == nil || len(group.Events) == 0 {
		return nil
	}
	if !SamplingResult(metadata.SamplingInterval) {
		return nil
	}
	c, ok := GetClient(metadata.MetricsPluginName)
	if !ok {
		return nil
	}
	return c.LogMonitorEvent(ctx, metadata, group)
}

// SamplingResult Sampling results
func SamplingResult(interval uint32) bool {
	if interval == 0 {
		return false
	}
	if interval > 1 {
		randValue := rand.Int63n(int64(interval))
		if randValue != int64(interval)-1 {
			return false
		}
	}
	return true
}

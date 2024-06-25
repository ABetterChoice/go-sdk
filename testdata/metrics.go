// Package testdata 测试数据，请勿使用
package testdata

import (
	"context"

	"github.com/abetterchoice/go-sdk/plugin/metrics"
	"github.com/abetterchoice/protoc_cache_server"
	"github.com/abetterchoice/protoc_event_server"
)

type empty struct{}

// LogExposure TODO
func (a empty) LogExposure(ctx context.Context, metadata *metrics.Metadata,
	exposureGroup *protoc_event_server.ExposureGroup) error {
	return nil
}

// LogEvent TODO
func (a empty) LogEvent(ctx context.Context, metadata *metrics.Metadata,
	eventGroup *protoc_event_server.EventGroup) error {
	return nil
}

// LogMonitorEvent TODO
func (a empty) LogMonitorEvent(ctx context.Context, metadata *metrics.Metadata,
	eventGroup *protoc_event_server.MonitorEventGroup) error {
	return nil
}

// Init TODO
func (a empty) Init(ctx context.Context, config *protoc_cache_server.MetricsInitConfig) error {
	return nil
}

// SendData TODO
func (a empty) SendData(ctx context.Context, metadata *metrics.Metadata, data [][]string) error {
	return nil
}

// Name TODO
func (a empty) Name() string {
	return "empty"
}

var (
	// EmptyMetricsClient TODO
	// Deprecated: for test
	EmptyMetricsClient = &empty{}
)

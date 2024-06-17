// Package metrics TODO
package metrics

// Metadata The metadata reported, in addition to the specific reported data, additional metadata required
type Metadata struct {
	MetricsPluginName string `json:"metricsPluginName"` // Monitoring plugin name
	TableName         string `json:"tableName"`         // Specific table name
	TableID           string `json:"tableId"`           // Specific table ID
	Token             string `json:"token"`             // Token
	SamplingInterval  uint32 `json:"samplingInterval"`  // Sampling interval
}

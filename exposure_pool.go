// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments,
// user configuration data retrieval, user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/abetterchoice/go-sdk/plugin/log"
	"github.com/abetterchoice/protoc_event_server"
)

type experimentExposure struct {
	projectID string
	list      *ExperimentList
	et        protoc_event_server.ExposureType
}

type experimentEvent struct {
	projectID string
	list      *ExperimentList
	latency   time.Duration
	optionStr string
	err       error
}

type remoteConfigExposure struct {
	projectID    string
	configResult *ConfigResult
	et           protoc_event_server.ExposureType
}

type remoteConfigEvent struct {
	projectID    string
	configResult *ConfigResult
	latency      time.Duration
	optionStr    string
	err          error
}

var (
	// ExperimentExposureChanSize TODO
	ExperimentExposureChanSize = 1 << 19
	// ExperimentEventChanSize TODO
	ExperimentEventChanSize = 1 << 19
	// RemoteConfigExposureChanSize TODO
	RemoteConfigExposureChanSize = 1 << 19
	// RemoteConfigEventChanSize TODO
	RemoteConfigEventChanSize = 1 << 19
)

var (
	experimentExposureChan   = make(chan *experimentExposure, ExperimentExposureChanSize)
	experimentEventChan      = make(chan *experimentEvent, ExperimentEventChanSize)
	remoteConfigExposureChan = make(chan *remoteConfigExposure, RemoteConfigExposureChanSize)
	remoteConfigEventChan    = make(chan *remoteConfigEvent, RemoteConfigEventChanSize)
)

var (
	defaultMaxParallelism = 4
)

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	if maxProcs > defaultMaxParallelism {
		return maxProcs
	}
	return defaultMaxParallelism
}

// initExposureConsumer Initialize exposure reporting consumer
func initExposureConsumer() {
	for i := 0; i < maxParallelism(); i++ {
		go watchData()
	}
}

// asyncExposureExperiments asynchronous push
// Record exposure data. If passive exposure is not enabled, you can use the Exposure API for manual exposure
// Manual exposure can avoid the overexposure problem that may be caused by passive exposure. Users can use manual exposure to report the exposure of the experiment they hit
func asyncExposureExperiments(projectID string, list *ExperimentList,
	exposureType protoc_event_server.ExposureType) error {
	select {
	case experimentExposureChan <- &experimentExposure{
		projectID: projectID,
		list:      list,
		et:        exposureType,
	}:
		return nil
	default:
		return fmt.Errorf("experimentExposureChan is full")
	}
}

// asyncExposureExperimentEvent async exposure
func asyncExposureExperimentEvent(projectID string, list *ExperimentList,
	latency time.Duration, optionStr string, err error) error {
	select {
	case experimentEventChan <- &experimentEvent{
		projectID: projectID,
		list:      list,
		latency:   latency,
		optionStr: optionStr,
		err:       err,
	}:
		return nil
	default:
		return fmt.Errorf("experimentEventChan is full")
	}
}

// asyncExposureRemoteConfig async exposure
func asyncExposureRemoteConfig(projectID string, configResult *ConfigResult,
	exposureType protoc_event_server.ExposureType) error {
	select {
	case remoteConfigExposureChan <- &remoteConfigExposure{
		projectID:    projectID,
		configResult: configResult,
		et:           exposureType,
	}:
		return nil
	default:
		return fmt.Errorf("remoteConfigExposureChan is full")
	}
}

// asyncExposureRemoteConfigEvent async exposure
func asyncExposureRemoteConfigEvent(projectID string, configResult *ConfigResult,
	latency time.Duration, optionStr string, err error) error {
	select {
	case remoteConfigEventChan <- &remoteConfigEvent{
		projectID:    projectID,
		configResult: configResult,
		latency:      latency,
		optionStr:    optionStr,
		err:          err,
	}:
		return nil
	default:
		return fmt.Errorf("remoteConfigEventChan is full")
	}
}

func watchData() {
	for {
		logExposure()
	}
}

func logExposure() {
	defer func() {
		recoverErr := recover() // Prevent third-party monitoring reporting plugins from panicking
		if recoverErr != nil {
			body := make([]byte, 1<<10)
			runtime.Stack(body, false)
			log.Errorf("recoverErr:%v\n%s", recoverErr, body)
			return
		}
	}()
	select {
	case eExposure := <-experimentExposureChan:
		if eExposure == nil || eExposure.list == nil || len(eExposure.list.Data) == 0 {
			return
		}
		err := exposureExperiments(context.TODO(), eExposure.projectID, eExposure.list, eExposure.et)
		if err != nil {
			// log.Errorf("exposureExperiments fail:%v", err)
		}
	case eEvent := <-experimentEventChan:
		if eEvent == nil || eEvent.list == nil || len(eEvent.list.Data) == 0 {
			return
		}
		err := exposureExperimentEvent(context.TODO(), eEvent.projectID, eEvent.list, eEvent.latency, eEvent.optionStr,
			eEvent.err)
		if err != nil {
			// log.Errorf("exposureExperimentEvent fail:%v", err)
		}
	case cExposure := <-remoteConfigExposureChan:
		if cExposure == nil || cExposure.configResult == nil {
			return
		}
		err := exposureRemoteConfig(context.TODO(), cExposure.projectID, cExposure.configResult, cExposure.et)
		if err != nil {
			log.Errorf("exposureRemoteConfig fail:%v", err)
		}
	case cEvent := <-remoteConfigEventChan:
		if cEvent == nil || cEvent.configResult == nil {
			return
		}
		err := exposureRemoteConfigEvent(context.TODO(), cEvent.projectID, cEvent.configResult, cEvent.latency,
			cEvent.optionStr, cEvent.err)
		if err != nil {
			log.Errorf("exposureRemoteConfig fail:%v", err)
		}
	}
}

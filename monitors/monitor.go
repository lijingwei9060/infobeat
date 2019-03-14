// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package monitors

import (
	"fmt"
	"sync"

	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"

	"github.com/elastic/beats/heartbeat/scheduler"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/lijingwei9060/infobeat/monitors/jobs"
	"github.com/lijingwei9060/infobeat/monitors/wrappers"
)

// Monitor represents a configured recurring monitoring configuredJob loaded from a config file. Starting it
// will cause it to run with the given scheduler until Stop() is called.
type Monitor struct {
	id             string
	name           string
	typ            string
	pluginName     string
	config         *common.Config
	registrar      *pluginsReg
	uniqueName     string
	scheduler      *scheduler.Scheduler
	configuredJobs []*configuredJob
	enabled        bool

	// internalsMtx is used to synchronize access to critical
	// internal datastructures
	internalsMtx sync.Mutex

	pipelineConnector beat.PipelineConnector

	// stats is the countersRecorder used to record lifecycle events
	// for global metrics + telemetry
	stats           registryRecorder
	factoryMetadata *common.MapStrPointer
}

// String prints a description of the monitor in a threadsafe way. It is important that this use threadsafe
// values because it may be invoked from another thread in cfgfile/runner.
func (m *Monitor) String() string {
	return fmt.Sprintf("Monitor<pluginName: %s, enabled: %t>", m.name, m.enabled)
}

func checkMonitorConfig(config *common.Config, registrar *pluginsReg) error {
	m, err := newMonitor(config, registrar, nil, nil, nil)
	m.Stop() // Stop the monitor to free up the ID from uniqueness checks
	return err
}

// uniqueMonitorIDs is used to keep track of explicitly configured monitor IDs and ensure no duplication within a
// given heartbeat instance.
var uniqueMonitorIDs sync.Map

// ErrDuplicateMonitorID is returned when a monitor attempts to start using an ID already in use by another monitor.
type ErrDuplicateMonitorID struct{ ID string }

func (e ErrDuplicateMonitorID) Error() string {
	return fmt.Sprintf("monitor ID %s is configured for multiple monitors! IDs must be unique values.", e.ID)
}

func newMonitor(
	config *common.Config,
	registrar *pluginsReg,
	pipelineConnector beat.PipelineConnector,
	scheduler *scheduler.Scheduler,
	factoryMetadata *common.MapStrPointer,
) (*Monitor, error) {
	// Extract just the Id, Type, and Enabled fields from the config
	// We'll parse things more precisely later once we know what exact type of
	// monitor we have
	mpi, err := pluginInfo(config)
	if err != nil {
		return nil, err
	}

	monitorPlugin, found := registrar.get(mpi.Type)
	if !found {
		return nil, fmt.Errorf("monitor type %v does not exist, valid types are %v", mpi.Type, registrar.monitorNames())
	}

	m := &Monitor{
		id:                mpi.ID,
		name:              mpi.Name,
		typ:               mpi.Type,
		pluginName:        monitorPlugin.name,
		scheduler:         scheduler,
		configuredJobs:    []*configuredJob{},
		pipelineConnector: pipelineConnector,
		internalsMtx:      sync.Mutex{},
		config:            config,
		stats:             monitorPlugin.stats,
		factoryMetadata:   factoryMetadata,
	}

	if m.id != "" {
		// Ensure we don't have duplicate IDs
		if _, loaded := uniqueMonitorIDs.LoadOrStore(m.id, m); loaded {
			return nil, ErrDuplicateMonitorID{m.id}
		}
	} else {
		// If there's no explicit ID generate one
		hash, err := m.configHash()
		if err != nil {
			return nil, err
		}
		m.id = fmt.Sprintf("auto-%s-%#X", m.typ, hash)
	}

	rawJobs, err := monitorPlugin.create(config)
	wrappedJobs := wrappers.WrapCommon(rawJobs, m.id, m.name, m.typ)

	if err != nil {
		return nil, fmt.Errorf("job err %v", err)
	}

	m.configuredJobs, err = m.makeTasks(config, wrappedJobs)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Monitor) configHash() (uint64, error) {
	unpacked := map[string]interface{}{}
	err := m.config.Unpack(unpacked)
	if err != nil {
		return 0, err
	}
	hash, err := hashstructure.Hash(unpacked, nil)
	if err != nil {
		return 0, err
	}

	return hash, nil
}

func (m *Monitor) makeTasks(config *common.Config, jobs []jobs.Job) ([]*configuredJob, error) {
	mtConf := jobConfig{}
	if err := config.Unpack(&mtConf); err != nil {
		return nil, errors.Wrap(err, "invalid config, could not unpack monitor config")
	}

	var mTasks []*configuredJob
	for _, job := range jobs {
		t, err := newConfiguredJob(job, mtConf, m)
		if err != nil {
			// Failure to compile monitor processors should not crash hb or prevent progress
			if _, ok := err.(ProcessorsError); ok {
				logp.Critical("Failed to load monitor processors: %v", err)
				continue
			}

			return nil, err
		}

		mTasks = append(mTasks, t)
	}

	return mTasks, nil
}

// Start starts the monitor's execution using its configured scheduler.
func (m *Monitor) Start() {
	m.internalsMtx.Lock()
	defer m.internalsMtx.Unlock()

	for _, t := range m.configuredJobs {
		t.Start()
	}
}

// Stop stops the Monitor's execution in its configured scheduler.
// This is safe to call even if the Monitor was never started.
func (m *Monitor) Stop() {
	m.internalsMtx.Lock()
	defer m.internalsMtx.Unlock()
	defer m.freeID()

	for _, t := range m.configuredJobs {
		t.Stop()
	}

}

func (m *Monitor) freeID() {
	// Free up the monitor ID for reuse
	uniqueMonitorIDs.Delete(m.id)
}

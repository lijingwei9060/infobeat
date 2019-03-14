package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/heartbeat/scheduler"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/lijingwei9060/infobeat/monitors"
	"github.com/pkg/errors"

	"github.com/lijingwei9060/infobeat/config"
)

// Infobeat configuration.
type Infobeat struct {
	done      chan struct{}
	config    config.Config
	scheduler *scheduler.Scheduler
	client    beat.Client
}

// New creates an instance of infobeat.
func New(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {
	parsedConfig := config.DefaultConfig
	if err := rawConfig.Unpack(&parsedConfig); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	limit := parsedConfig.Scheduler.Limit
	locationName := parsedConfig.Scheduler.Location
	if locationName == "" {
		locationName = "Local"
	}
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return nil, err
	}

	scheduler := scheduler.NewWithLocation(limit, location)
	bt := &Infobeat{
		done:      make(chan struct{}),
		config:    parsedConfig,
		scheduler: scheduler,
	}
	return bt, nil
}

// Run starts infobeat.
func (bt *Infobeat) Run(b *beat.Beat) error {
	logp.Info("infobeat is running! Hit CTRL-C to stop it.")

	err := bt.RunStaticMonitors(b)
	if err != nil {
		return err
	}

	if err := bt.scheduler.Start(); err != nil {
		return err
	}
	defer bt.scheduler.Stop()

	<-bt.done

	logp.Info("Shutting down.")
	return nil
}

// RunStaticMonitors runs the `heartbeat.monitors` portion of the yaml config if present.
func (bt *Infobeat) RunStaticMonitors(b *beat.Beat) error {
	factory := monitors.NewFactory(bt.scheduler)

	for _, cfg := range bt.config.Monitors {
		created, err := factory.Create(b.Publisher, cfg, nil)
		if err != nil {
			return errors.Wrap(err, "could not create monitor")
		}
		created.Start()
	}
	return nil
}

// Stop stops infobeat.
func (bt *Infobeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

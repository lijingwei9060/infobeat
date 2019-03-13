package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/lijingwei9060/infobeat/config"
)

// Infobeat configuration.
type Infobeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of infobeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Infobeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts infobeat.
func (bt *Infobeat) Run(b *beat.Beat) error {
	logp.Info("infobeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops infobeat.
func (bt *Infobeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

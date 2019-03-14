// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Monitors  []*common.Config `config:"monitors"`
	Scheduler Scheduler        `config:"scheduler"`
	Period    time.Duration    `config:"period"`
}

// Scheduler defines the syntax of a heartbeat.yml scheduler block.
type Scheduler struct {
	Limit    uint   `config:"limit"  validate:"min=0"`
	Location string `config:"location"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}

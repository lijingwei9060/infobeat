package command

import (
	"path/filepath"
	"testing"

	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Tests with different params from config file
	absPath, err := filepath.Abs("../../tests/files/")

	// assert.NotNil(t, absPath)
	// assert.Nil(t, err)

	config := &Config{}

	// Reads  config file
	err = cfgfile.Read(config, absPath+"/command.yml")
	assert.Nil(t, err)
}

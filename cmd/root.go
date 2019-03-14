package cmd

import (
	"github.com/lijingwei9060/infobeat/beater"
	_ "github.com/lijingwei9060/infobeat/monitors/defaults"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "infobeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)

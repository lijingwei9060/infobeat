package versionfetch

import (
	"github.com/lijingwei9060/infobeat/utils/command"
)

// VersionFetcher 获取系统版本信息
type VersionFetcher interface {
	Get(*command.Commander) ([]string, error)
}

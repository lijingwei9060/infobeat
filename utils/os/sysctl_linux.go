package os

import (
	"bufio"
	"bytes"

	"github.com/elastic/beats/libbeat/common"
	"github.com/lijingwei9060/infobeat/utils/command"
)

// GetSysctl 获取服务器的sysctl信息
func GetSysctl() common.MapStr {
	ret := common.MapStr{}
	out, err := command.Commander.RunCommand("sysctl", "-a")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	str := []string{}

	for scanner.Scan() {
		str = append(str, scanner.Text())
	}
	if len(str) > 0 {
		ret.Put("sysctl", str)
	}
	return ret
}

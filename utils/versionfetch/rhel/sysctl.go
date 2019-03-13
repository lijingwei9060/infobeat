package rhel

import (
	"bufio"
	"bytes"

	"github.com/lijingwei9060/infobeat/utils/command"
)

// FetchSysctl 获取sysctl 内容
func FetchSysctl(commander command.Commander) ([]string, error) {
	out, err := commander.RunCommand("sysctl", "-a")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	str := []string{}

	for scanner.Scan() {
		str = append(str, scanner.Text())
	}
	return str, nil
}

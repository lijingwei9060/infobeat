package rhel

import (
	"bufio"
	"bytes"

	"github.com/lijingwei9060/infobeat/utils/command"
)

// FetchRpmVersion 获取所有rpm包的版本
func FetchRpmVersion(commander command.Commander) ([]string, error) {
	out, err := commander.RunCommand("rpm", "-qa", "--queryformat", "'%{name}:%{version}.%{release}\n'")
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

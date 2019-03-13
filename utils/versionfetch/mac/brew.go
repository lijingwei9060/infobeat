package mac

import (
	"bufio"
	"bytes"

	"github.com/lijingwei9060/infobeat/utils/command"
)

// FetchBrewVersion 获取所有rpm包的版本
func FetchBrewVersion(commander command.Commander) ([]string, error) {
	out, err := commander.RunCommand("brew", "list", "--versions")
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

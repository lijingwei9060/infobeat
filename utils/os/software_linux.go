package os

import (
	"bufio"
	"bytes"

	"github.com/lijingwei9060/infobeat/utils/command"

	"github.com/elastic/beats/libbeat/common"
)

func GetSoftware() common.MapStr {
	ret := common.MapStr{}
	out, err := command.Commander.RunCommand("rpm", "-qa", "--queryformat", "'%{name}:%{version}.%{release}\n'")
	if err != nil {
		return ret
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	str := []string{}

	for scanner.Scan() {
		str = append(str, scanner.Text())
	}

	if len(str) > 0 {
		ret.Put("rpm", str)
	}

	dpkg, err := command.Commander.RunCommand("dpkg", "-l", "--showformat='${Package}:${Version}\n'")
	if err != nil {
		return ret
	}

	dpkgscanner = bufio.NewScanner(bytes.NewReader(dpkg))
	dpkgstr := []string{}

	for dpkgscanner.Scan() {
		str = append(dpkgstr, dpkgscanner.Text())
	}

	if len(dpkgstr) > 0 {
		ret.Put("dpkg", dpkgstr)
	}

	return ret
}

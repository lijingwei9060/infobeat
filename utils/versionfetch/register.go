package versionfetch

import "fmt"

var register = make(map[string]*VersionFetcher)

// Register 注册获取版本的方法
func Register(name string, vf *VersionFetcher) error {
	if name == "" {
		return fmt.Errorf("name is nil")
	}

	if vf == nil {
		return fmt.Errorf("vf is nil")
	}

	if _, ok := register[name]; ok {
		return fmt.Errorf("%s duplicated", name)
	}

	register[name] = vf
	return nil
}

// GetVersionFetcher 根据name获取工具
func GetVersionFetcher(name string) (*VersionFetcher, error) {
	if _, exists := register[name]; !exists {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	return register[name], nil
}

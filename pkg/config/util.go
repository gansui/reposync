package config

import (
	"os"
	"path/filepath"
	"strings"
)

func findFirstExistingConfigFile(locations []string) (string, error) {
	for _, file := range fileLocations {
		// resolve env vars
		file = os.ExpandEnv(file)

		// skip if path contains unresolved env vars
		if strings.Contains(file, "$") {
			continue
		}

		// check if file exists
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", os.ErrNotExist
}

func firstValidLocation(locations []string) (string, error) {
	for _, file := range fileLocations {
		// resolve env vars
		file = os.ExpandEnv(file)

		// skip if path contains unresolved env vars
		if strings.Contains(file, "$") {
			continue
		}

		return file, nil
	}

	return "", os.ErrNotExist
}

func createParentDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0644)
	}

	return nil
}

// GetConfigDir returns the directory of the config file
func GetConfigDir(configFile string) string {
	if configFile != "" {
		return filepath.Dir(configFile)
	}

	// 尝试查找配置文件
	file, err := findFirstExistingConfigFile(fileLocations)
	if err == nil {
		return filepath.Dir(file)
	}

	// 默认返回当前目录
	return "."
}

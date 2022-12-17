package config

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"os"
	"path"
	"strings"
)

type Config struct {
	Name        string
	Executables map[string]struct {
		Path string
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func folderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func findConfigFilePath(searchPath string) (string, error) {
	run := true

	configPath := ""

	for run {
		if fileExists(searchPath + "/x.yml") {
			configPath = strings.Clone(searchPath + "/x.yml")
			run = false
		}
		if fileExists(searchPath + "/x.yaml") {
			configPath = strings.Clone(searchPath + "/x.yaml")
			run = false
		}

		if folderExists(searchPath + "/.git") {
			run = false
		}
		if searchPath == "/" {
			run = false
		}
		searchPath = path.Dir(searchPath)
	}

	if configPath == "" {
		return configPath, errors.New("config file not found")
	}

	return configPath, nil
}

func New(cwd string, additionalConfigFiles []string) (*Config, error) {
	configFilePath, err := findConfigFilePath(cwd)
	if err != nil {
		return nil, err
	}
	k := koanf.New(".")
	if err = k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return nil, err
	}
	for _, configFile := range additionalConfigFiles {
		if !fileExists(configFile) {
			continue
		}
		fmt.Println(configFile)
		if err = k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return nil, err
		}
	}
	cfg := &Config{}
	if err = k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

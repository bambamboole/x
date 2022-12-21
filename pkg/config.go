package pkg

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"path"
	"strings"
)

type Config struct {
	Name        string
	Executables map[string]struct {
		Path string
	}
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

func NewConfig(cwd string, additionalConfigFiles []string) (Config, error) {
	cfg := Config{}
	configFilePath, err := findConfigFilePath(cwd)
	if err != nil {
		return cfg, err
	}
	k := koanf.New(".")
	if err = k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return cfg, err
	}
	for _, configFile := range additionalConfigFiles {
		if !fileExists(configFile) {
			continue
		}
		if err = k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return cfg, err
		}
	}
	if err = k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

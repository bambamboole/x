package config

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"path"
	"strings"
	"x/pkg/utils"
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
		if utils.FileExists(searchPath + "/x.yml") {
			configPath = strings.Clone(searchPath + "/x.yml")
			run = false
		}
		if utils.FileExists(searchPath + "/x.yaml") {
			configPath = strings.Clone(searchPath + "/x.yaml")
			run = false
		}

		if utils.FolderExists(searchPath + "/.git") {
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

func New(cwd string, additionalConfigFiles []string) (Config, error) {
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
		if !utils.FileExists(configFile) {
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

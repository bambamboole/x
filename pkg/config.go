package pkg

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Name        string
	Executables map[string]struct {
		Path string
	}
}

func NewConfig(cwd string, additionalConfigFiles []string) (Config, error) {
	cfg := Config{}
	k := koanf.New(".")
	for _, configFile := range findFiles(cwd, []string{"x.yml", "x.yaml", "x.local.yml", "x.local.yaml"}) {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return cfg, err
		}
	}
	for _, configFile := range additionalConfigFiles {
		if !fileExists(configFile) {
			continue
		}
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return cfg, err
		}
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

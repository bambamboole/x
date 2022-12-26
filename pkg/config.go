package pkg

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"path"
)

type Config struct {
	Name        string
	Executables map[string]struct {
		Path string
	}
}

func NewConfig(logger IOLoggerInterface, projectPath string, additionalConfigFiles []string) (Config, error) {
	configFiles := append(
		additionalConfigFiles,
		path.Join(projectPath, "x.yml"),
		path.Join(projectPath, "x.yaml"),
		path.Join(projectPath, "x.local.yml"),
		path.Join(projectPath, "x.local.yaml"),
	)
	cfg := Config{}
	k := koanf.New(".")
	for _, configFile := range configFiles {
		if !fileExists(configFile) {
			continue
		}
		logger.Log("Found config file: "+configFile, DebugOn)
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return cfg, err
		}
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}
	logger.Log("Loaded config:", DebugVeryVerbose)
	logger.Log(cfg, DebugVeryVerbose)

	return cfg, nil
}

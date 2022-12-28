package pkg

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"path"
)

type Config struct {
	engine      *koanf.Koanf
	Name        string
	Executables map[string]struct {
		Path string
	}
}

func (c Config) loadFile(filePath string) error {
	if err := c.engine.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		return err
	}
	return nil
}

func (c Config) populate() error {
	if err := c.engine.Unmarshal("", &c); err != nil {
		return err
	}
	return nil
}

func (c Config) PopulateModuleConfig(moduleName string, cfg interface{}) error {
	if err := c.engine.Unmarshal("modules."+moduleName, cfg); err != nil {
		return err
	}
	return nil
}

func NewConfig(logger IOLoggerInterface, projectPath string, additionalConfigFiles []string) (Config, error) {
	configFiles := append(
		additionalConfigFiles,
		path.Join(projectPath, "x.yml"),
		path.Join(projectPath, "x.yaml"),
		path.Join(projectPath, "x.local.yml"),
		path.Join(projectPath, "x.local.yaml"),
	)
	cfg := Config{engine: koanf.New(".")}
	for _, configFile := range configFiles {
		if !fileExists(configFile) {
			continue
		}
		logger.Log("Found config file: "+configFile, DebugOn)

		if err := cfg.loadFile(configFile); err != nil {
			return cfg, err
		}
	}
	if err := cfg.populate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

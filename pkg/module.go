package pkg

import "github.com/bambamboole/x/pkg/modules"

type Module interface {
	GetConfig() interface{}
	GetScript() (string, error)
}

type ModuleRegistry struct {
	modules map[string]Module
}

func (m *ModuleRegistry) GetModules() map[string]Module {
	return m.modules
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: map[string]Module{"docker": &modules.DockerModule{}},
	}
}

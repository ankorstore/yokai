package fxcore

import (
	"fmt"
	"sort"

	"go.uber.org/fx"
)

// FxModuleInfoRegistry is the registry collecting info about registered modules.
type FxModuleInfoRegistry struct {
	infos map[string]FxModuleInfo
}

// FxModuleInfoRegistryParam allows injection of the required dependencies in [NewFxModuleInfoRegistry].
type FxModuleInfoRegistryParam struct {
	fx.In
	Infos []any `group:"core-module-infos"`
}

// NewFxModuleInfoRegistry returns a new [FxModuleInfoRegistry].
func NewFxModuleInfoRegistry(p FxModuleInfoRegistryParam) *FxModuleInfoRegistry {
	infos := make(map[string]FxModuleInfo)

	for _, info := range p.Infos {
		if castInfo, ok := info.(FxModuleInfo); ok {
			infos[castInfo.Name()] = castInfo
		}
	}

	return &FxModuleInfoRegistry{
		infos: infos,
	}
}

func (r *FxModuleInfoRegistry) Names() []string {
	names := make([]string, len(r.infos))

	i := 0
	for name := range r.infos {
		names[i] = name
		i++
	}

	sort.Strings(names)

	return names
}

// All returns a map of all registered [FxModuleInfo].
func (r *FxModuleInfoRegistry) All() map[string]FxModuleInfo {
	return r.infos
}

// Find returns a [FxModuleInfo] by name.
func (r *FxModuleInfoRegistry) Find(name string) (FxModuleInfo, error) {
	if info, ok := r.infos[name]; ok {
		return info, nil
	}

	return nil, fmt.Errorf("fx module info with name %s was not found", name)
}

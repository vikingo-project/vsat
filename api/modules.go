package api

import "github.com/vikingo-project/vsat/modules"

func (a *APIC) Modules() (*RecordsContainer, error) {
	avaliableModules := modules.GetAvaliableModules()
	var modules []map[string]interface{}
	for _, m := range avaliableModules {
		modules = append(modules, m.GetInfo())
	}
	return &RecordsContainer{Records: modules, Total: int64(len(modules))}, nil
}

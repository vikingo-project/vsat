package modules

import (
	"reflect"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/utils"
)

type Module interface {
	GetName() string
	GetInfo() map[string]interface{}
	GetSettings() interface{}
	GetDefaultSettings() interface{}

	Init(string, int, interface{}, *events.EventsAPI) // IP, Port, Settings, pusher
	Up() error
	Down() error
}

// Settings map[string]interface{}
var AvaliableModules = make(map[string]Module)
var Instances = make(map[string]Module)

func Register(m Module) {
	moduleInfo := m.GetInfo()
	utils.PrintDebug("Register module %v", moduleInfo)
	AvaliableModules[m.GetName()] = m
}

func GetAvaliableModules() []Module {
	var modules []Module
	for _, m := range AvaliableModules {
		modules = append(modules, m)
	}
	return modules
}

func clone(old Module) Module {
	n := reflect.New(reflect.TypeOf(old).Elem())
	val := reflect.ValueOf(old).Elem()
	nVal := n.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		if nvField.CanSet() {
			nvField.Set(val.Field(i))
		}
	}
	return n.Interface().(Module)
}

func GetModuleByName(name string) Module {
	for _, m := range AvaliableModules {
		if m.GetName() == name {
			return clone(m)
		}
	}
	return nil
}

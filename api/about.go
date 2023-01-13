package api

import "github.com/vikingo-project/vsat/shared"

func (a *APIC) About() map[string]interface{} {
	return map[string]interface{}{"Version": shared.Version, "Build": shared.BuildHash}
}

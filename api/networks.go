package api

import "github.com/vikingo-project/vsat/utils"

func (a *APIC) Networks() (*RecordsContainer, error) {
	networks, _ := utils.GetNetworks()
	return &RecordsContainer{Records: networks, Total: int64(len(networks))}, nil
}

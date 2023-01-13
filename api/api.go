package api

import (
	"context"
	"log"

	"github.com/vikingo-project/vsat/manager"
)

type APIC struct {
}

func (a *APIC) GetManager() *manager.Manager {
	mgr := CTX.Value("mgr")
	log.Printf("mgr=%+v", mgr)
	return mgr.(*manager.Manager)
}

type Any interface{}

type RecordsContainer struct {
	Records Any
	Total   int64
}

var (
	Instance APIC            // API instance
	CTX      context.Context // global context
)

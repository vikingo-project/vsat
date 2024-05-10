package utils

import (
	"github.com/vikingo-project/glp"
	lua "github.com/yuin/gopher-lua"
	lfs "layeh.com/gopher-lfs"
)

func NewLuaState() *lua.LState {
	L := lua.NewState(lua.Options{
		MinimizeStackMemory: true,
		CallStackSize:       5000,
	})
	L.OpenLibs()
	glp.Preload(L)
	lfs.Preload(L)
	return L
}

package vshttp

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/vikingo-project/vsat/utils"
	lua "github.com/yuin/gopher-lua"
)

type luaHTTPRes struct {
	http.ResponseWriter
	req      *http.Request
	complete bool
	done     chan bool
}

func serveLua(path, action string, data interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings luaSettings
		err := mapstructure.Decode(data, &settings)
		if err != nil {
			log.Println("failed to decode lua settings")
			return
		}

		state := utils.NewLuaState()
		defer state.Close()

		luaReq := state.NewTable()
		bodyReader := state.NewFunction(func(L *lua.LState) int {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(lua.LString(string(data)))
			return 1
		})
		luaReq.RawSetString(`body`, bodyReader)
		luaReq.RawSetString(`host`, lua.LString(r.Host))
		luaReq.RawSetString(`method`, lua.LString(r.Method))
		luaReq.RawSetString(`referer`, lua.LString(r.Referer()))
		luaReq.RawSetString(`proto`, lua.LString(r.Proto))
		luaReq.RawSetString(`user_agent`, lua.LString(r.UserAgent()))
		if r.URL != nil && len(r.URL.Query()) > 0 {
			query := state.NewTable()
			for k, v := range r.URL.Query() {
				if len(v) > 0 {
					query.RawSetString(k, lua.LString(v[0]))
				}
			}
			luaReq.RawSetString(`query`, query)
		}
		if len(r.Header) > 0 {
			headers := state.NewTable()
			for k, v := range r.Header {
				if len(v) > 0 {
					headers.RawSetString(k, lua.LString(v[0]))
				}
			}
			luaReq.RawSetString(`headers`, headers)
		}
		luaReq.RawSetString(`path`, lua.LString(r.URL.Path))
		luaReq.RawSetString(`raw_path`, lua.LString(r.URL.RawPath))
		luaReq.RawSetString(`raw_query`, lua.LString(r.URL.RawQuery))
		luaReq.RawSetString(`request_uri`, lua.LString(r.RequestURI))
		luaReq.RawSetString(`remote_addr`, lua.LString(r.RemoteAddr))

		if err := state.DoString(settings.Code); err != nil {
			log.Printf("failed to exec lua code %s", err.Error())
			return
		}

		luaRes := state.NewTable()
		luaRes.RawSetString(`status_code`, lua.LNumber(200))
		luaRes.RawSetString(`headers`, state.NewTable())
		luaRes.RawSetString(`body`, lua.LString(""))

		if err := state.CallByParam(lua.P{
			Fn:      state.GetGlobal("handler"),
			NRet:    1,
			Protect: true,
		}, luaReq, luaRes); err != nil {
			log.Printf("failed to call handler; error: %+v", err.Error())
			return
		}
		res := state.CheckTable(1)
		resStatusCode := res.RawGet(lua.LString("status_code"))
		code, _ := strconv.Atoi(resStatusCode.String())
		w.WriteHeader(code)

		resHeaders := res.RawGet(lua.LString("headers"))
		if headers, ok := resHeaders.(*lua.LTable); ok {
			headers.ForEach(func(name lua.LValue, value lua.LValue) {
				w.Header().Add(name.String(), value.String())
			})
		}
		resBody := res.RawGet(lua.LString("body"))
		w.Write([]byte(resBody.String()))
	}
}

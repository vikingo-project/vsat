package vshttp

import (
	"net"
	"net/http"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
)

func proxyHeader(header http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range header {
		k = strings.Replace(k, "-", "_", -1)
		m[k] = v[0]
	}
	return m

}
func Render(tpl string, r *http.Request) ([]byte, error) {
	template, err := pongo2.FromString(tpl)
	if err != nil {
		return []byte{}, err
	}
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	out, err := template.Execute(pongo2.Context{
		"client_ip": clientIP,
		"req": pongo2.Context{
			"path":        r.URL.Path,
			"uri":         r.RequestURI,
			"headers":     proxyHeader(r.Header), // todo
			"remote_addr": r.RemoteAddr,
		},
		"vars": mux.Vars(r),
	})
	if err != nil {
		return []byte{}, err
	}
	return []byte(out), nil
}

package pathproxy

import (
	"github.com/wunderlist/moxy"
	"testing"
)

func Test_pathProxy_StartServer(t *testing.T) {
		cnf := NewEnvProxyConfig("", "/ggl/->google.com", []moxy.FilterFunc{})
		p := PathProxy(cnf)
		p.StartServer()
}

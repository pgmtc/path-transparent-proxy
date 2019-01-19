package pathproxy

import (
	"github.com/pgmtc/orchard-gateway-go/pkg/github.com/pgmtc/httptesting"
	"github.com/wunderlist/moxy"
	"os"
	"strings"
	"testing"
	"time"
)

func setUp(t *testing.T) {
	go httptesting.StartHttpServer(":9091", "Server1 response")
	go httptesting.StartHttpServer(":9092", "Server2 response")
	go httptesting.StartHttpServer(":9093", "Server3 response")

	time.Sleep(1 * time.Second)
	httptesting.AssertResponse(t, "server1 online", "Server1 response:/", "http://localhost:9091")
	httptesting.AssertResponse(t, "server2 online", "Server2 response:/", "http://localhost:9092")
	httptesting.AssertResponse(t, "server3 online", "Server3 response:/", "http://localhost:9093")
}

func Test_pathProxy_StartServer(t *testing.T) {
	setUp(t)

	os.Setenv("PORT", "9090")
	routeArray := []string{
		"/server1->localhost:9091",
		"/server2->localhost:9092",
		"/server3/->localhost:9093",
	}
	routes := strings.Join(routeArray, ",")
	cnf := NewEnvProxyConfig("", routes, []moxy.FilterFunc{})
	p := PathProxy(cnf)
	go p.StartServer()

	time.Sleep(1 * time.Second)

	httptesting.AssertResponse(t, "server1 via proxy", "Server1 response:/server1", "http://localhost:9090/server1")
	httptesting.AssertResponse(t, "server2 via proxy", "Server2 response:/server2/", "http://localhost:9090/server2/")
	httptesting.AssertResponse(t, "server2 via proxy, query params", "Server2 response:/server2/?param=1234", "http://localhost:9090/server2/?param=1234")
	httptesting.AssertResponse(t, "server3 via proxy path stripped", "Server3 response:/", "http://localhost:9090/server3")
	httptesting.AssertResponse(t, "server3 via proxy path stripped ", "Server3 response:/", "http://localhost:9090/server3/")
	httptesting.AssertResponse(t, "server3 via proxy path stripped ", "Server3 response:/?query=abcd", "http://localhost:9090/server3/?query=abcd")

	// Test health
	httptesting.AssertResponse(t, "health check", "{\"status\":\"UP\"}", "http://localhost:9090/health")
}

func Test_pathProxy_StartInvalid(t *testing.T) {
	// Test startup fail
	os.Setenv("PORT", "-1")
	cnf := NewEnvProxyConfig("", "", []moxy.FilterFunc{})
	p := PathProxy(cnf)
	err := p.StartServer()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}


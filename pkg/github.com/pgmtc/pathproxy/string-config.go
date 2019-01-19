package pathproxy

import (
	"fmt"
	"github.com/wunderlist/moxy"
	"strings"
)

type EnvProxyConfig struct {
	routes []routeConfig
	filters []moxy.FilterFunc
	contextPath string
	stripPath bool
}

func (e *EnvProxyConfig) StripPath() bool {
	return e.stripPath
}

func (e *EnvProxyConfig) GetContextPath() string {
	return e.contextPath
}

func (e *EnvProxyConfig) GetRoutes() []routeConfig {
	return e.routes
}

func (e *EnvProxyConfig) GetFilters() []moxy.FilterFunc {
	return e.filters
}

func NewEnvProxyConfig(contextPath string, routesString string, filters []moxy.FilterFunc) *EnvProxyConfig{
	cnf := &EnvProxyConfig{
		contextPath:contextPath,
		routes: parseRoutes(routesString),
		filters: filters,
	}
	return cnf

}

func parseRoutes(routesString string) (routes []routeConfig) {
	routesArr := strings.Split(routesString, ",")
	for _, route := range routesArr {
		route = strings.Trim(route, " ")
		if route == "" {
			continue
		}
		routeParts := strings.Split(route, "->")
		if len(routeParts) != 2 {
			panic(fmt.Sprintf("Unable to parse route %s", route))
		}

		routes = append(routes, routeConfig{
			path: routeParts[0],
			host: routeParts[1],
		})
	}
	return
}



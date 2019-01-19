package pathproxy

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/wunderlist/moxy"
	"net/http"
	"os"
	"strings"
)

const DEFAULT_PORT = "8080"

type routeConfig struct {
	path string
	host string
	contextPath string
}

type ProxyConfig interface {
	GetRoutes() []routeConfig
	GetFilters() []moxy.FilterFunc
	GetContextPath() string
}

type pathProxy struct {
	server *http.Server
	app *negroni.Negroni
	config ProxyConfig
}

func getAddr() string {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = DEFAULT_PORT
	}
	return ":" + envPort
}

func (p *pathProxy) StartServer() (returnError error) {
	router := p.buildRouter()
	p.app = negroni.Classic()
	p.app.UseHandler(router)

	server := &http.Server{Addr: getAddr()}
	server.Handler = p.app
	server.Addr = getAddr()
	err := server.ListenAndServe()
	if err != nil {
		returnError = err
	}
	return
}

func PathProxy(config ProxyConfig) *pathProxy {
	pproxy := &pathProxy{
		config:config,
	}
	return pproxy
}

func (p *pathProxy) buildRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("{\"status\":\"UP\"}"))
	})
	for _, route := range p.config.GetRoutes() {
		p.addRoute(router, route.host, route.path)
	}
	return router
}

func (p *pathProxy) addRoute(router *mux.Router, host string, path string) *mux.Route {
	fullPath := p.config.GetContextPath() + path
	trimmedFullPath := strings.TrimRight(fullPath, "/")
	fmt.Printf("Registering %s -> %s\n", fullPath, host)

	filters := p.config.GetFilters()
	proxy := moxy.NewReverseProxy([]string{host}, filters)

	proxy.Director = func(request *http.Request) {
		targetPath := translateTargetPath(request.URL.Path, p.config.GetContextPath(), path)
		// Modify the request
		request.URL.Path = targetPath
		request.URL.Scheme = "http"
		request.URL.Host = host
		fmt.Println(request.URL)
	}

	return router.PathPrefix(trimmedFullPath).HandlerFunc(proxy.ServeHTTP)
}

func translateTargetPath(requestUri string, contextPath string, routePath string) (targetPath string, ) {
	targetPath = strings.Replace(requestUri, contextPath, "", 1)
	if strings.HasSuffix(routePath, "/") {
		targetPath = strings.Replace(targetPath, strings.TrimSuffix(routePath, "/"), "", 1)
	}
	return
}


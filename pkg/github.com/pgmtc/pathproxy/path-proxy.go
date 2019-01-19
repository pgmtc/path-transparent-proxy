package pathproxy

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/wunderlist/moxy"
	"log"
	"net/http"
	"os"
	"strings"
)

type routeConfig struct {
	path string
	host string
	contextPath string
}

type ProxyConfig interface {
	GetRoutes() []routeConfig
	GetFilters() []moxy.FilterFunc
	GetContextPath() string
	StripPath() bool
}

type pathProxy struct {
	app *negroni.Negroni
	config ProxyConfig
}

func (p *pathProxy) StartServer() {
	router := p.buildRouter()
	p.app = negroni.Classic()
	p.app.UseHandler(router)
	p.app.Run()
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
		writer.Write([]byte("Working fine"))
	})

	if len(p.config.GetRoutes()) == 0 {
		log.Fatal("No routes have been provided. Exiting")
		os.Exit(1)
	}
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
	fmt.Println(len(routePath))
	fmt.Println(strings.Index(routePath, "/"))
	if strings.HasSuffix(routePath, "/") {
		targetPath = strings.Replace(targetPath, routePath, "", 1)
	}
	return
}


package main

import (
    "github.com/pgmtc/orchard-gateway-go/pkg/github.com/pgmtc/pathproxy"
    "github.com/wunderlist/moxy"
    "log"
    "net/http"
    "net/url"
    "os"
)

var CONTEXT_PATH = "/orchard-gateway-msvc"
var ROUTES = ""

func init() {
    CONTEXT_PATH = os.Getenv("CONTEXT_PATH")
    ROUTES = os.Getenv("ROUTES")

    if ROUTES == "" {
       log.Fatal("Missing route information. Provide it in ROUTES environmental variable. Example:\n" +
           "export PATH=\"/path-1->server1:8080,/path-2->server2:8080\"")
    }
}

func main() {
    config := pathproxy.NewEnvProxyConfig(CONTEXT_PATH, ROUTES, []moxy.FilterFunc{AddSecurityHeaders})
    pathproxy.PathProxy(config).StartServer()
}


func AddSecurityHeaders(request *http.Request, response *http.Response) {
    locationHeader := response.Header.Get("Location")
    if locationHeader != "" && request.RequestURI == "/orchard-gateway-msvc/orchard-auth-msvc/saml/SSO" {
        redirectUri, _ := url.Parse(locationHeader)
        redirectUri.Path = CONTEXT_PATH + redirectUri.Path
        response.Header.Set("Location", redirectUri.String())
        log.Printf("Rewriting Location request:\n from: %s \n to: %s\n", locationHeader, response.Header.Get("Location"))
    }
    response.Header.Del("X-Powered-By")
    response.Header.Set("X-Super-Secure", "Yes!!")
}
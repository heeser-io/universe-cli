package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/heeser-io/universe-cli/helper"
	v2 "github.com/heeser-io/universe/api/v2"
	"github.com/heeser-io/universe/services/v2/gateway"
	"github.com/olekukonko/tablewriter"
)

var (
	p *Proxy
)

type ProxyPath struct {
	p       *httputil.ReverseProxy
	gateway *v2.Gateway
	route   gateway.Route
}
type Proxy struct {
	functionProxies sync.Map
	proxies         sync.Map
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (proxy *Proxy) AddFunction(fName, host, port string) error {
	u, err := url.Parse(fmt.Sprintf("http://%v:%s/", host, port))
	if err != nil {
		log.Printf("Error parsing URL")
	}

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// req.Host = host
			// req.URL.Scheme = "http"
			// req.URL.Host = fmt.Sprintf("%s.localhost:%s", fName, port)
			// req.URL.Path = "/"
			req.Host = host
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.URL.Path = "/"
		},
	}

	proxyVal, ok := proxy.functionProxies.Load(fName)

	savedProxy := &httputil.ReverseProxy{}
	if ok {
		savedProxy = proxyVal.(*httputil.ReverseProxy)
	} else {
		savedProxy = p
	}

	proxy.functionProxies.Store(fName, savedProxy)

	return nil
}

func (proxy *Proxy) Add(gateway *v2.Gateway, route gateway.Route, host string) error {
	u, err := url.Parse(fmt.Sprintf("http://%v/", host))
	if err != nil {
		log.Printf("Error parsing URL")
	}

	targetQuery := u.RawQuery
	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = host
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.URL.Path = strings.ReplaceAll(req.URL.Path, fmt.Sprintf("/%s", gateway.Name), "")
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		},
	}

	proxyVal, ok := proxy.proxies.Load(gateway.Name)

	savedProxy := []ProxyPath{}
	if ok {
		savedProxy = proxyVal.([]ProxyPath)
	}

	savedProxy = append(savedProxy, ProxyPath{
		p:       p,
		gateway: gateway,
		route:   route,
	})

	proxy.proxies.Store(gateway.Name, savedProxy)

	return nil
}

func New() *Proxy {
	if p == nil {
		p = &Proxy{
			proxies: sync.Map{},
		}
	}

	return p
}

func FuncProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func ProxyRequestHandler(pp *ProxyPath) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pp.p.ServeHTTP(w, r)
	}
}

func (p *Proxy) Listen() error {
	router := mux.NewRouter()

	proxyPort := "12121"

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TYPE", "NAME", "HOST"})

	fmt.Printf("\n PROXY UP AND RUNNING\n")

	p.functionProxies.Range(func(functionName, val any) bool {
		proxy := val.(*httputil.ReverseProxy)
		router.HandleFunc("/", FuncProxyRequestHandler(proxy))
		// router.HandleFunc(fmt.Sprintf("/%s%s", functionName, `/{all:[a-zA-Z0-9=\-\/]+}`), ProxyRequestHandler(proxy))
		// router.HandleFunc(fmt.Sprintf("/%s", functionName), ProxyRequestHandler(proxy))

		table.Append([]string{
			"function",
			functionName.(string),
			fmt.Sprintf("http://%s.localhost:%s", functionName, proxyPort),
		})
		return true
	})

	p.proxies.Range(func(serviceName, val any) bool {
		ppArr := val.([]ProxyPath)
		// router.HandleFunc(fmt.Sprintf("/%s%s", serviceName, `/{all:[a-zA-Z0-9=\-\/]+}`), ProxyRequestHandler(proxy))

		paths := map[string][]ProxyPath{}
		for _, pp := range ppArr {
			if paths[pp.route.Path] == nil {
				paths[pp.route.Path] = []ProxyPath{}
			}
			paths[pp.route.Path] = append(paths[pp.route.Path], pp)
		}

		table.Append([]string{
			"gateway",
			serviceName.(string),
			fmt.Sprintf("http://localhost:%s/%s", proxyPort, serviceName),
		})

		for routePath, ppArr := range paths {
			for _, pp := range ppArr {
				router.HandleFunc(fmt.Sprintf("/%s", path.Join(serviceName.(string), routePath)), ProxyRequestHandler(&pp)).Methods(pp.route.Method)
				table.Append([]string{
					"route",
					pp.route.Method,
					fmt.Sprintf("http://localhost:%s/%s", proxyPort, path.Join(serviceName.(string), routePath)),
				})
			}
		}

		return true
	})

	table.Render()

	helper.KillPort(proxyPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", proxyPort), router); err != nil {
		panic(err)
	}
	return nil
}

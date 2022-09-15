package gpproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/koding/websocketproxy"
)

type NodeService interface {
	ProxyHTTP(http.ResponseWriter, *http.Request)
	ProxyWS(http.ResponseWriter, *http.Request)
}

type nodeService struct {
	httpURL   *url.URL
	wsURL     *url.URL
	httpProxy *httputil.ReverseProxy
	wsProxy   *websocketproxy.WebsocketProxy
}

func NewNodeService(httpStr string, wsStr string) (NodeService, error) {
	httpURL, err := url.Parse(httpStr)
	if err != nil {
		return nil, fmt.Errorf("node: %w", err)
	}

	wsURL, err := url.Parse(wsStr)
	if err != nil {
		return nil, fmt.Errorf("node: %w", err)
	}

	node := &nodeService{
		httpURL: httpURL,
		wsURL:   wsURL,
	}

	if err := node.init(); err != nil {
		return nil, fmt.Errorf("node: %w", err)
	}

	return node, nil
}

func (node *nodeService) init() error {
	node.wsProxy = &websocketproxy.WebsocketProxy{
		Backend: func(r *http.Request) *url.URL {
			u := *node.wsURL
			u.Fragment = node.wsURL.Fragment
			u.Path = node.wsURL.Path
			u.RawQuery = node.wsURL.RawQuery
			return &u
		},
		Director: func(req *http.Request, out http.Header) {
			out.Set("Host", node.wsURL.Host)
		},
	}

	node.httpProxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.URL = node.httpURL
			req.Host = node.httpURL.Host
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		},
	}

	return nil
}

func (node *nodeService) ProxyHTTP(w http.ResponseWriter, r *http.Request) {
	node.httpProxy.ServeHTTP(w, r)
}

func (node *nodeService) ProxyWS(w http.ResponseWriter, r *http.Request) {
	node.wsProxy.ServeHTTP(w, r)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

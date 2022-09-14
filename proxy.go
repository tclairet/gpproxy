package gpproxy

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type API struct {
	NodeService NodeService
}

func NewAPI(node NodeService) *API {
	api := new(API)
	api.NodeService = node

	return api
}

func (api *API) Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", api.nodeHTTPService)
	router.Handle("/ws", &WSHandler{api.NodeService})
	router.Handle("/eth/gasprice", http.HandlerFunc(api.ethGasPrice))

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	return router
}

func (api *API) ethGasPrice(w http.ResponseWriter, r *http.Request) {
	api.NodeService.ProxyHTTP(w, r)
}

func (api *API) nodeHTTPService(w http.ResponseWriter, r *http.Request) {
	api.NodeService.ProxyHTTP(w, r)
}

type WSHandler struct {
	NodeService
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.NodeService.ProxyWS(w, r)
}

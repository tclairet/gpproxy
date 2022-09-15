package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gpproxy"
	"net/http"
	"strings"
)

func main() {
	logger := logrus.New().WithFields(logrus.Fields{"service": "proxy"})

	if err := LoadConfig(); err != nil {
		logger.Fatal(err)
	}
	node, err := gpproxy.NewNodeService(config.NodeURL, config.NodeWSURL)
	if err != nil {
		logger.Fatal(err)
	}
	api := gpproxy.NewAPI(node)

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: api.Routes()}
	logger.Fatal(server.ListenAndServe())
}

func handleInfuraNodeWS(url string) string {
	url = strings.Replace(url, "http", "ws", 1)
	if strings.Contains(url, "infura.io/v3") {
		url = strings.Replace(url, "infura.io/v3", "infura.io/ws/v3", 1)
	}
	return url
}

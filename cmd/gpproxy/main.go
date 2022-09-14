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
	wsURL := strings.Replace(config.NodeURL, "http", "ws", 1)
	node, err := gpproxy.NewNodeService(config.NodeURL, wsURL)
	if err != nil {
		logger.Fatal(err)
	}
	api := gpproxy.NewAPI(node)

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: api.Routes()}
	logger.Fatal(server.ListenAndServe())
}

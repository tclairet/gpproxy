package gpproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var nodeHTTP, nodeWS string

func TestMain(m *testing.M) {
	nodeHTTP, _ = os.LookupEnv("NODE_URL")
	nodeWS = nodeHTTP
	if nodeHTTP == "" {
		fakeNode := NewFakeEthNode()
		defer fakeNode.Close()
		nodeHTTP = fakeNode.server.URL
		nodeWS = fakeNode.server.URL + "/ws"
	}
	nodeWS = strings.Replace(nodeWS, "http", "ws", 1)
	if strings.Contains(nodeWS, "infura.io/v3") {
		nodeWS = strings.Replace(nodeWS, "infura.io/v3", "infura.io/ws/v3", 1)
	}

	m.Run()
}

func TestHttpAPI(t *testing.T) {
	node, err := NewNodeService(nodeHTTP, "")
	if err != nil {
		t.Fatal(err)
	}

	api := NewAPI(node)
	server := httptest.NewServer(api.Routes())
	defer server.Close()

	tests := []struct {
		name    string
		request *http.Request
	}{
		{
			name:    "proxy http",
			request: makeHTTPRequest(t, http.MethodPost, server.URL+"/", (&Request{Method: "eth_gasPrice"}).Bytes()),
		},
		{
			name:    "eth/gasprice endpoint with rpc request",
			request: makeHTTPRequest(t, http.MethodPost, server.URL+"/eth/gasprice", (&Request{Method: "eth_gasPrice"}).Bytes()),
		},
		{
			name:    "eth/gasprice endpoint without rpc request and POST",
			request: makeHTTPRequest(t, http.MethodPost, server.URL+"/eth/gasprice", nil),
		},
		{
			name:    "eth/gasprice endpoint without rpc request and GET",
			request: makeHTTPRequest(t, http.MethodGet, server.URL+"/eth/gasprice", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.request)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := resp.StatusCode, http.StatusOK; got != want {
				b, _ := io.ReadAll(resp.Body)
				t.Fatalf("got %v, want %v. Body: %s", got, want, b)
			}

			actual := struct {
				Result string `json:"result"`
			}{}
			if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
				t.Fatal(err)
			}
			if actual.Result == "" {
				t.Fatal("result is empty")
			}
		})
	}
}

func TestWSAPI(t *testing.T) {
	node, err := NewNodeService("", nodeWS)
	if err != nil {
		t.Fatal(err)
	}

	api := NewAPI(node)
	server := httptest.NewServer(api.Routes())
	defer server.Close()

	c, _, err := websocket.DefaultDialer.Dial(strings.Replace(server.URL, "http", "ws", 1)+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.WriteJSON(&Request{Method: "eth_gasPrice"}); err != nil {
		t.Fatal(err)
	}
	actual := struct {
		Result string `json:"result"`
	}{}
	if err := c.ReadJSON(&actual); err != nil {
		t.Fatal(err)
	}

	if actual.Result == "" {
		t.Fatal("result is empty")
	}
}

func TestEndpoints(t *testing.T) {
	api := NewAPI(nil)
	server := httptest.NewServer(api.Routes())
	defer server.Close()

	tests := []struct {
		endpoint string
	}{
		{
			endpoint: server.URL + "/metrics",
		},
		{
			endpoint: server.URL + "/healthz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			resp, err := http.Get(tt.endpoint)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := resp.StatusCode, http.StatusOK; got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func makeHTTPRequest(t *testing.T, method string, url string, b []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

type FakeEthNode struct {
	requestReceived *http.Request
	server          *httptest.Server
}

func NewFakeEthNode() *FakeEthNode {
	fakeEthNode := &FakeEthNode{}

	router := mux.NewRouter()
	router.HandleFunc("/", fakeEthNode.handlerHTTP).Methods(http.MethodPost)
	router.HandleFunc("/ws", fakeEthNode.handlerWS)

	fakeEthNode.server = httptest.NewServer(router)

	return fakeEthNode
}

func (fakeEthNode *FakeEthNode) Close() {
	fakeEthNode.server.Close()
}

func (fakeEthNode *FakeEthNode) handlerHTTP(w http.ResponseWriter, r *http.Request) {
	fakeEthNode.requestReceived = r
	b, _ := ioutil.ReadAll(r.Body)

	var rpcRequest Request
	json.Unmarshal(b, &rpcRequest)

	if rpcRequest.Method == "eth_gasPrice" {
		resp := MakeResponse(rpcRequest, "0xbeef")
		RespondWithJSON(w, http.StatusOK, &resp)
		return
	}
}

func (fakeEthNode *FakeEthNode) handlerWS(w http.ResponseWriter, r *http.Request) {
	fakeEthNode.requestReceived = r
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	defer c.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "upgrade: %v", err)
		return
	}
	for {
		var rpcRequest Request

		if err := c.ReadJSON(&rpcRequest); err != nil {
			return
		}

		if rpcRequest.Method == "eth_gasPrice" {
			_ = c.WriteJSON(&Response{
				Result: "0xbeef",
			})
		}
	}
}

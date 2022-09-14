package gpproxy

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHttpAPI(t *testing.T) {
	url, ok := os.LookupEnv("NODE_URL")
	if !ok {
		t.Skip("NODE_URL not set")
	}
	node, err := NewNodeService(url, strings.Replace(url, "http", "ws", 1))
	if err != nil {
		t.Fatal(err)
	}

	api := NewAPI(node)
	server := httptest.NewServer(api.Routes())
	defer server.Close()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "proxy http",
			url:  server.URL + "/",
		},
		{
			name: "eth/gasprice endpoint",
			url:  server.URL + "/eth/gasprice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := POST(t, tt.url, (&Request{Method: "eth_gasPrice"}).Bytes())
			if got, want := resp.StatusCode, http.StatusOK; got != want {
				t.Fatalf("got %v, want %v", got, want)
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
	url, ok := os.LookupEnv("NODE_URL")
	if !ok {
		t.Skip("NODE_URL not set")
	}
	node, err := NewNodeService(url, strings.Replace(url, "http", "ws", 1))
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

func POST(t *testing.T, url string, b []byte) *http.Response {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

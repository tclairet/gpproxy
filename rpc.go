package gpproxy

import (
	"encoding/json"
)

// Request defines ethereum rpc request
type Request struct {
	ID      uint          `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Bytes convert the request in bytes
func (r *Request) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}
